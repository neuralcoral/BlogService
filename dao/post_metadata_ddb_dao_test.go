package dao

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/neuralcoral/BlogService/model"
	"github.com/stretchr/testify/assert"
)

type MockDynamoDBClient struct {
	GetItemFunc func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	PutItemFunc func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	ScanFunc    func(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
}

func (m *MockDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return m.GetItemFunc(input)
}

func (m *MockDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return m.PutItemFunc(input)
}

func (m *MockDynamoDBClient) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return m.ScanFunc(input)
}

func TestGetPostMetadata_Succeeds(t *testing.T) {
	ID := "123"
	Title := "Title Post"
	BodyUrl := "http://example.com/bodyText"
	PreviewText := "This is a preview"
	Status := model.Draft
	CreatedAt := time.Now().Format(time.RFC3339)
	UpdatedAt := time.Now().Format(time.RFC3339)
	mockDynamoDBClient := &MockDynamoDBClient{
		GetItemFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{
				Item: map[string]*dynamodb.AttributeValue{
					"ID":          {S: aws.String(ID)},
					"Title":       {S: aws.String(Title)},
					"BodyUrl":     {S: aws.String(BodyUrl)},
					"PreviewText": {S: aws.String(PreviewText)},
					"Status":      {S: aws.String(string(Status))},
					"CreatedAt":   {S: aws.String(CreatedAt)},
					"UpdatedAt":   {S: aws.String(UpdatedAt)},
				},
			}, nil
		},
	}

	sut := PostMetadataDdbDao{
		client: mockDynamoDBClient,
	}

	result, err := sut.GetPostMetadata("123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parsedCreatedAt, _ := time.Parse(time.RFC3339, CreatedAt)
	parsedUpdatedAt, _ := time.Parse(time.RFC3339, UpdatedAt)

	expected := &model.PostMetadata{
		ID:          ID,
		Title:       Title,
		BodyUrl:     BodyUrl,
		PreviewText: PreviewText,
		Status:      Status,
		CreatedAt:   parsedCreatedAt,
		UpdatedAt:   parsedUpdatedAt,
	}

	assert.Equal(t, expected, result)
}
