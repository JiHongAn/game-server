# Game Server

Go로 작성된 게임 서버 애플리케이션입니다. HTTP API와 WebSocket을 통한 실시간 매칭 기능을 제공합니다.

## 🏗️ 프로젝트 구조

```
game-server/
├── cmd/                    # 애플리케이션 진입점
│   └── server/            # 메인 서버 애플리케이션
├── internal/              # 내부 패키지 (외부에서 import 불가)
│   ├── config/           # 설정 관리
│   ├── domain/           # 도메인 모델 (엔티티, 값 객체)
│   ├── handler/          # HTTP 핸들러
│   ├── service/          # 비즈니스 로직 (ORM 직접 사용)
│   ├── middleware/       # 미들웨어
│   ├── socket/           # 웹소켓 관련
│   └── pkg/              # 내부 공통 유틸리티
│       ├── auth/         # 인증 관련
│       ├── errors/       # 에러 처리
│       └── response/     # 응답 처리
├── pkg/                  # 외부에서 사용 가능한 공개 라이브러리
│   └── dto/              # 데이터 전송 객체
├── configs/              # 설정 파일들
├── build/                # 빌드 결과물
├── Makefile             # 빌드 자동화
├── go.mod               # Go 모듈 정의
└── README.md            # 프로젝트 설명
```

## 🚀 시작하기

### 필요 조건

- Go 1.23 이상
- Make (선택사항)

### 설치 및 실행

1. **저장소 클론**
   ```bash
   git clone <repository-url>
   cd game-server
   ```

2. **의존성 설치**
   ```bash
   make deps
   # 또는
   go mod download
   ```

3. **환경 변수 설정**
   ```bash
   cp configs/env.example .env
   # .env 파일을 편집하여 필요한 값들을 설정
   ```

4. **애플리케이션 실행**
   ```bash
   make run
   # 또는
   go run cmd/server/main.go
   ```

### 개발 모드

Hot reload를 지원하는 개발 모드로 실행:

```bash
make dev
```

## 🔧 사용 가능한 명령어

```bash
make build         # 애플리케이션 빌드
make run           # 애플리케이션 실행
make dev           # 개발 모드 (hot reload)
make test          # 테스트 실행
make test-coverage # 테스트 커버리지 리포트
make lint          # 코드 린팅
make fmt           # 코드 포맷팅
make clean         # 빌드 결과물 정리
make help          # 도움말
```

## 📡 API 엔드포인트

### HTTP API

- `GET /` - Health check
- `GET /games` - 게임 목록 조회 (JWT 인증 필요)

### WebSocket

- `:8081` - 매치 서버 (실시간 매칭)

## 🏛️ 아키텍처

이 프로젝트는 간단한 레이어드 아키텍처를 따릅니다:

- **Domain Layer**: 비즈니스 엔티티와 규칙 (`internal/domain`)
- **Service Layer**: 비즈니스 로직 및 ORM 사용 (`internal/service`)
- **Handler Layer**: HTTP 핸들러 (`internal/handler`)
- **Middleware Layer**: 인증, 로깅 등 (`internal/middleware`)

### ORM 사용

Repository 패턴 대신 Service 계층에서 ORM을 직접 사용하여 데이터베이스와 상호작용합니다.

## 🔐 인증

JWT(JSON Web Token)를 사용한 Bearer 토큰 인증을 지원합니다.

## 🧪 테스트

```bash
# 모든 테스트 실행
make test

# 커버리지 리포트 생성
make test-coverage
```

## 📝 환경 변수

| 변수명 | 설명 | 기본값 |
|--------|------|--------|
| `PORT` | HTTP 서버 포트 | `8080` |
| `MATCH_PORT` | 매치 서버 포트 | `8081` |
| `JWT_SECRET` | JWT 시크릿 키 | `your-secret-key` |
| `JWT_EXPIRES_IN` | JWT 만료 시간 (시간) | `24` |
| `DB_HOST` | 데이터베이스 호스트 | `localhost` |
| `DB_PORT` | 데이터베이스 포트 | `5432` |
| `DB_USERNAME` | 데이터베이스 사용자명 | `postgres` |
| `DB_PASSWORD` | 데이터베이스 비밀번호 | - |
| `DB_DATABASE` | 데이터베이스 이름 | `gameserver` |

## 🤝 기여하기

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 라이선스

이 프로젝트는 MIT 라이선스 하에 배포됩니다.