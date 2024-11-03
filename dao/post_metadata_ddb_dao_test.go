package dao

import (
	"errors"
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
	getItemFunc := func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
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
	}

	sut := setupMockDynamoDBForGet(t, getItemFunc)

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

func TestGetPostMetadata_DynamoDBFailure_ReturnsErr(t *testing.T) {
	expectedErr := errors.New("mock error for testing")

	getItemFunc := func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
		return nil, errors.New("mock error for testing")
	}

	sut := setupMockDynamoDBForGet(t, getItemFunc)

	result, err := sut.GetPostMetadata("123")

	if result != nil {
		t.Fatalf("unexpected result: %v", result)
	}

	assert.Equal(t, expectedErr, err)
}

func TestGetPostMetadata_EmptyOutput_ReturnsEmptyOutput(t *testing.T) {
	getItemFunc := func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
		return nil, nil
	}

	sut := setupMockDynamoDBForGet(t, getItemFunc)

	result, err := sut.GetPostMetadata("123")

	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	assert.Nil(t, result)

}

func TestUpdatePostMetadata_Succeeds(t *testing.T) {
	ID := "123"
	Title := "Title Post"
	BodyUrl := "http://example.com/bodyText"
	PreviewText := "This is a preview"
	Status := model.Draft
	CreatedAt := time.Now().Format(time.RFC3339)
	UpdatedAt := time.Now().Format(time.RFC3339)

	parsedCreatedAt, _ := time.Parse(time.RFC3339, CreatedAt)
	parsedUpdatedAt, _ := time.Parse(time.RFC3339, UpdatedAt)

	expectedResult := &model.PostMetadata{
		ID:          ID,
		Title:       Title,
		BodyUrl:     BodyUrl,
		PreviewText: PreviewText,
		Status:      Status,
		CreatedAt:   parsedCreatedAt,
		UpdatedAt:   parsedUpdatedAt,
	}

	putItemFunc := func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
		return &dynamodb.PutItemOutput{
			Attributes: map[string]*dynamodb.AttributeValue{
				"ID":          {S: aws.String(ID)},
				"Title":       {S: aws.String(Title)},
				"BodyUrl":     {S: aws.String(BodyUrl)},
				"PreviewText": {S: aws.String(PreviewText)},
				"Status":      {S: aws.String(string(Status))},
				"CreatedAt":   {S: aws.String(CreatedAt)},
				"UpdatedAt":   {S: aws.String(UpdatedAt)},
			},
		}, nil
	}

	sut := setupMockDynamoDBForPut(t, putItemFunc)

	input := &model.PostMetadata{
		ID:          ID,
		Title:       Title,
		BodyUrl:     BodyUrl,
		PreviewText: PreviewText,
		Status:      Status,
		CreatedAt:   parsedCreatedAt,
	}

	result, err := sut.UpdatePostMetadata(input)

	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	assert.Equal(t, expectedResult, result)

}

func TestUpdatePostMetadata_DynamoDBFailure_ReturnsErr(t *testing.T) {
	expectedErr := errors.New("mock error for testing")
	putItemFunc := func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
		return nil, errors.New("mock error for testing")
	}
	sut := setupMockDynamoDBForPut(t, putItemFunc)

	input := &model.PostMetadata{}

	result, err := sut.UpdatePostMetadata(input)

	if result != nil {
		t.Fatalf("unexpected result: %v", result)
	}

	assert.Equal(t, expectedErr, err)
}

func TestUpdatePostMetadata_EmptyInput_ReturnsEmpty(t *testing.T) {

	putItemFunc := func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
		return nil, nil
	}

	sut := setupMockDynamoDBForPut(t, putItemFunc)

	result, err := sut.UpdatePostMetadata(nil)

	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	assert.Nil(t, result)
}

func setupMockDynamoDBForGet(t testing.TB, getItemFunc func(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)) PostMetadataDdbDao {
	t.Helper()
	mockDynamoDBClient := &MockDynamoDBClient{
		GetItemFunc: getItemFunc,
	}

	sut := PostMetadataDdbDao{
		client: mockDynamoDBClient,
	}

	return sut
}

func setupMockDynamoDBForPut(t testing.TB, putItemFunc func(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)) PostMetadataDdbDao {
	t.Helper()
	mockDynamoDBClient := &MockDynamoDBClient{
		PutItemFunc: putItemFunc,
	}

	sut := PostMetadataDdbDao{
		client: mockDynamoDBClient,
	}

	return sut
}
