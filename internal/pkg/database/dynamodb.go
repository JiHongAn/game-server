package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var dynamoClient *dynamodb.Client

// InitDynamoDB DynamoDB 클라이언트 초기화
// 게임 이벤트 저장용, DynamoDB Stream으로 S3에 백업
func InitDynamoDB() error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(getEnv("AWS_REGION", "ap-northeast-2")),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)

	// 테이블 존재 확인 및 생성
	err = createTableIfNotExists()
	if err != nil {
		log.Printf("Warning: Failed to create DynamoDB table: %v", err)
	}

	log.Println("DynamoDB client initialized for game events storage")
	return nil
}

func createTableIfNotExists() error {
	tableName := getEnv("DYNAMODB_TABLE_NAME", "game-events")

	// 테이블 존재 확인
	_, err := dynamoClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})

	if err == nil {
		log.Printf("DynamoDB table %s already exists", tableName)
		return nil
	}

	// 테이블 생성 (DynamoDB Stream 활성화)
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("gameId"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("eventId"),
				KeyType:       types.KeyTypeRange,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("gameId"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("eventId"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
		StreamSpecification: &types.StreamSpecification{
			StreamEnabled:  aws.Bool(true),
			StreamViewType: types.StreamViewTypeNewAndOldImages,
		},
	}

	_, err = dynamoClient.CreateTable(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	log.Printf("DynamoDB table %s created with stream enabled for S3 backup", tableName)
	return nil
}

// SaveEventToDynamoDB 게임 이벤트를 DynamoDB에 저장
// DynamoDB Stream을 통해 자동으로 S3에 백업됨
func SaveEventToDynamoDB(event interface{}) error {
	if dynamoClient == nil {
		return fmt.Errorf("DynamoDB client not initialized")
	}

	tableName := getEnv("DYNAMODB_TABLE_NAME", "game-events")

	item, err := attributevalue.MarshalMap(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}

	_, err = dynamoClient.PutItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to put item to DynamoDB: %w", err)
	}

	return nil
}

// QueryEventsByGameID 특정 게임의 모든 이벤트 조회 (분석용)
func QueryEventsByGameID(gameID string) ([]map[string]interface{}, error) {
	if dynamoClient == nil {
		return nil, fmt.Errorf("DynamoDB client not initialized")
	}

	tableName := getEnv("DYNAMODB_TABLE_NAME", "game-events")

	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("gameId = :gameId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":gameId": &types.AttributeValueMemberS{Value: gameID},
		},
	}

	result, err := dynamoClient.Query(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to query DynamoDB: %w", err)
	}

	var events []map[string]interface{}
	for _, item := range result.Items {
		var event map[string]interface{}
		err := attributevalue.UnmarshalMap(item, &event)
		if err != nil {
			log.Printf("Failed to unmarshal item: %v", err)
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
