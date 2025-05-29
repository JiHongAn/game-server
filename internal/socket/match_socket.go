package socket

import (
	"encoding/json"
	"fmt"
	"game-server/internal/dto"
	"game-server/internal/pkg/auth"
	"game-server/internal/pkg/database"
	"game-server/internal/service"
	"log"
	"net"
	"sync"
)

type Client struct {
	ID     string
	UserID string
	Conn   net.Conn
}

type SocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	From string      `json:"from,omitempty"`
}

type MatchServer struct {
	listener     net.Listener
	clients      map[string]*Client // socketId -> Client
	clientsMux   sync.RWMutex
	matchService *service.MatchService
}

const (
	INVITE_EXPIRE_MINUTES = 5
)

func NewMatchServer(matchService *service.MatchService) *MatchServer {
	return &MatchServer{
		clients:      make(map[string]*Client),
		matchService: matchService,
	}
}

func (s *MatchServer) Start(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to start match server: %v", err)
	}
	s.listener = listener

	log.Printf("Match Server listening on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *MatchServer) handleConnection(conn net.Conn) {
	clientAddr := conn.RemoteAddr().String()
	log.Printf("New client connected: %s", clientAddr)

	client := &Client{
		ID:   clientAddr,
		Conn: conn,
	}

	defer s.removeClient(client.ID)
	defer conn.Close()

	buffer := make([]byte, 4096)
	authenticated := false

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Error reading from client %s: %v", clientAddr, err)
			return
		}

		var msg SocketMessage
		if err := json.Unmarshal(buffer[:n], &msg); err != nil {
			log.Printf("Error parsing message from %s: %v", clientAddr, err)
			continue
		}

		// 인증 검사
		if !authenticated {
			if msg.Type != "auth" {
				s.sendErrorToClient(client, "Authentication required")
				continue
			}

			if !s.authenticateClient(client, &msg) {
				continue
			}

			authenticated = true
			s.addClient(client)
			log.Printf("Client %s authenticated as user %s", clientAddr, client.UserID)
			continue
		}

		msg.From = client.ID
		s.handleMessage(client, &msg)
	}
}

func (s *MatchServer) authenticateClient(client *Client, msg *SocketMessage) bool {
	var authData map[string]interface{}
	if msgData, ok := msg.Data.(map[string]interface{}); ok {
		authData = msgData
	}

	token, ok := authData["token"].(string)
	if !ok || token == "" {
		s.sendErrorToClient(client, "Token is required")
		return false
	}

	// JWT 토큰 검증
	jwtToken, err := auth.ValidateAccessToken(token)
	if err != nil {
		s.sendErrorToClient(client, "Invalid token")
		return false
	}

	// 사용자 ID 추출
	userID, err := auth.GetUserIDFromToken(jwtToken)
	if err != nil {
		s.sendErrorToClient(client, "Failed to get user ID from token")
		return false
	}

	client.UserID = userID

	// Redis에 socket_id와 user_id 매핑 저장
	if err := database.HSet("socket:users", client.ID, userID); err != nil {
		log.Printf("Failed to store socket mapping: %v", err)
	}
	if err := database.HSet("user:sockets", userID, client.ID); err != nil {
		log.Printf("Failed to store user mapping: %v", err)
	}

	s.sendToClient(client, SocketMessage{
		Type: "auth_success",
		Data: map[string]string{"userId": userID},
	})

	return true
}

func (s *MatchServer) handleMessage(client *Client, msg *SocketMessage) {
	switch msg.Type {
	case "create_match":
		s.handleCreateMatch(client, msg)
	case "invite_friends":
		s.handleInviteFriends(client, msg)
	case "respond_invite":
		s.handleRespondInvite(client, msg)
	case "start_match":
		s.handleStartMatch(client, msg)
	case "leave_match":
		s.handleLeaveMatch(client, msg)
	default:
		s.sendErrorToClient(client, "Unknown message type")
	}
}

func (s *MatchServer) handleCreateMatch(client *Client, msg *SocketMessage) {
	var req dto.CreateMatchRequest
	if err := s.parseMessageData(msg.Data, &req); err != nil {
		s.sendErrorToClient(client, "Invalid create match data")
		return
	}

	// 서비스로 위임
	matchInfo, err := s.matchService.CreateMatch(client.UserID, req.GameID, req.MaxPlayers)
	if err != nil {
		s.sendErrorToClient(client, err.Error())
		return
	}

	s.sendToClient(client, SocketMessage{
		Type: "match_created",
		Data: dto.CreateMatchResponse{
			MatchID:    matchInfo.MatchID,
			GameID:     matchInfo.GameID,
			HostID:     matchInfo.HostID,
			MaxPlayers: matchInfo.MaxPlayers,
			Message:    "Match created successfully",
		},
	})

	log.Printf("Match %s created by user %s for game %s", matchInfo.MatchID, client.UserID, matchInfo.GameID)
}

func (s *MatchServer) handleInviteFriends(client *Client, msg *SocketMessage) {
	var req dto.InviteFriendRequest
	if err := s.parseMessageData(msg.Data, &req); err != nil {
		s.sendErrorToClient(client, "Invalid invite data")
		return
	}

	// 서비스로 위임
	response, err := s.matchService.InviteFriends(client.UserID, req.MatchID, req.FriendIds)
	if err != nil {
		s.sendErrorToClient(client, err.Error())
		return
	}

	// 친구들에게 초대 알림 전송
	for _, friendID := range response.InvitedIds {
		if friendClient := s.getClientByUserID(friendID); friendClient != nil {
			// 초대 정보 다시 조회해서 전송
			inviteKey := fmt.Sprintf("invite:%s:%s", req.MatchID, friendID)
			inviteData, err := database.Get(inviteKey)
			if err == nil && inviteData != "" {
				var invitation dto.MatchInvitation
				if json.Unmarshal([]byte(inviteData), &invitation) == nil {
					s.sendToClient(friendClient, SocketMessage{
						Type: "match_invitation",
						Data: invitation,
					})
				}
			}
		}
	}

	s.sendToClient(client, SocketMessage{
		Type: "friends_invited",
		Data: response,
	})

	log.Printf("User %s invited %d friends to match %s", client.UserID, len(response.InvitedIds), req.MatchID)
}

func (s *MatchServer) handleRespondInvite(client *Client, msg *SocketMessage) {
	var req dto.InviteResponseRequest
	if err := s.parseMessageData(msg.Data, &req); err != nil {
		s.sendErrorToClient(client, "Invalid response data")
		return
	}

	// 서비스로 위임
	response, err := s.matchService.RespondInvite(client.UserID, req.MatchID, req.Response)
	if err != nil {
		s.sendErrorToClient(client, err.Error())
		return
	}

	// 응답 메시지 타입 결정
	msgType := "invite_declined"
	if req.Response == "accept" {
		msgType = "invite_accepted"
	}

	s.sendToClient(client, SocketMessage{
		Type: msgType,
		Data: response,
	})

	// accept인 경우 매치의 모든 플레이어에게 알림
	if req.Response == "accept" {
		players, err := s.matchService.GetMatchPlayers(req.MatchID)
		if err == nil {
			s.notifyMatchPlayers(req.MatchID, SocketMessage{
				Type: "player_joined",
				Data: map[string]interface{}{
					"matchId": req.MatchID,
					"userId":  client.UserID,
					"players": players,
				},
			}, client.UserID)
		}
	}

	log.Printf("User %s %s invitation for match %s", client.UserID, req.Response, req.MatchID)
}

func (s *MatchServer) handleStartMatch(client *Client, msg *SocketMessage) {
	var req dto.StartMatchRequest
	if err := s.parseMessageData(msg.Data, &req); err != nil {
		s.sendErrorToClient(client, "Invalid start match data")
		return
	}

	// 서비스로 위임
	response, err := s.matchService.StartMatch(client.UserID, req.MatchID)
	if err != nil {
		s.sendErrorToClient(client, err.Error())
		return
	}

	// 모든 플레이어에게 게임 시작 알림
	s.notifyMatchPlayers(req.MatchID, SocketMessage{
		Type: "match_started",
		Data: response,
	}, "")

	log.Printf("Match %s started by user %s", req.MatchID, client.UserID)
}

func (s *MatchServer) handleLeaveMatch(client *Client, msg *SocketMessage) {
	// 현재 매치 확인
	matchID, err := database.HGet("user:matches", client.UserID)
	if err != nil || matchID == "" {
		s.sendErrorToClient(client, "You are not in any match")
		return
	}

	// 서비스로 위임
	if err := s.matchService.LeaveMatch(client.UserID); err != nil {
		s.sendErrorToClient(client, err.Error())
		return
	}

	s.sendToClient(client, SocketMessage{
		Type: "match_left",
		Data: dto.SuccessResponse{Message: "Left the match"},
	})

	// 남은 플레이어들에게 알림
	players, err := s.matchService.GetMatchPlayers(matchID)
	if err == nil && len(players) > 0 {
		s.notifyMatchPlayers(matchID, SocketMessage{
			Type: "player_left",
			Data: map[string]interface{}{
				"matchId": matchID,
				"userId":  client.UserID,
				"players": players,
			},
		}, client.UserID)
	}

	log.Printf("User %s left match %s", client.UserID, matchID)
}

func (s *MatchServer) notifyMatchPlayers(matchID string, msg SocketMessage, excludeUserID string) {
	players, err := s.matchService.GetMatchPlayers(matchID)
	if err != nil {
		return
	}

	for _, player := range players {
		if player.UserID != excludeUserID {
			if client := s.getClientByUserID(player.UserID); client != nil {
				s.sendToClient(client, msg)
			}
		}
	}
}

func (s *MatchServer) getClientByUserID(userID string) *Client {
	socketID, err := database.HGet("user:sockets", userID)
	if err != nil || socketID == "" {
		return nil
	}

	s.clientsMux.RLock()
	defer s.clientsMux.RUnlock()
	return s.clients[socketID]
}

func (s *MatchServer) parseMessageData(data interface{}, target interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, target)
}

func (s *MatchServer) sendToClient(client *Client, msg SocketMessage) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	_, err = client.Conn.Write(msgBytes)
	if err != nil {
		log.Printf("Error sending message to client %s: %v", client.ID, err)
	}
}

func (s *MatchServer) sendErrorToClient(client *Client, message string) {
	s.sendToClient(client, SocketMessage{
		Type: "error",
		Data: dto.ErrorResponse{Message: message},
	})
}

func (s *MatchServer) addClient(client *Client) {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()
	s.clients[client.ID] = client
}

func (s *MatchServer) removeClient(clientID string) {
	s.clientsMux.Lock()
	client, exists := s.clients[clientID]
	if exists {
		delete(s.clients, clientID)
	}
	s.clientsMux.Unlock()

	if exists && client.UserID != "" {
		// Redis에서 매핑 제거
		database.HDel("socket:users", clientID)
		database.HDel("user:sockets", client.UserID)

		// 매치에서 제거
		if err := s.matchService.LeaveMatch(client.UserID); err == nil {
			log.Printf("Removed user %s from match due to disconnect", client.UserID)
		}

		log.Printf("Client %s (user: %s) disconnected", clientID, client.UserID)
	} else {
		log.Printf("Client %s disconnected", clientID)
	}
}
