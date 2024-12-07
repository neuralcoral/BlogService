package dao

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/neuralcoral/BlogService/model"
	"github.com/stretchr/testify/assert"
)

type MockDynamoDBClient struct {
	GetItemFunc func(context context.Context, input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	PutItemFunc func(context context.Context, input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	ScanFunc    func(context context.Context, input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
}

func (m *MockDynamoDBClient) GetItem(context context.Context, input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return m.GetItemFunc(context, input)
}

func (m *MockDynamoDBClient) PutItem(context context.Context, input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return m.PutItemFunc(context, input)
}

func (m *MockDynamoDBClient) Scan(context context.Context, input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return m.ScanFunc(context, input)
}

func TestGetPostMetadata_Succeeds(t *testing.T) {
	ID := "123"
	Title := "Title Post"
	BodyUrl := "http://example.com/bodyText"
	PreviewText := "This is a preview"
	Status := model.Draft
	CreatedAt := time.Now().Format(time.RFC3339)
	UpdatedAt := time.Now().Format(time.RFC3339)
	getItemFunc := func(context context.Context, input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
		return &dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"ID":          &types.AttributeValueMemberS{Value: ID},
				"Title":       &types.AttributeValueMemberS{Value: Title},
				"BodyUrl":     &types.AttributeValueMemberS{Value: BodyUrl},
				"PreviewText": &types.AttributeValueMemberS{Value: PreviewText},
				"Status":      &types.AttributeValueMemberS{Value: string(Status)},
				"CreatedAt":   &types.AttributeValueMemberS{Value: CreatedAt},
				"UpdatedAt":   &types.AttributeValueMemberS{Value: UpdatedAt},
			},
		}, nil
	}

	sut := setupMockDynamoDBForGet(t, getItemFunc)

	result, err := sut.GetPostMetadata(context.Background(), "123")

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

	getItemFunc := func(context context.Context, input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
		return nil, errors.New("mock error for testing")
	}

	sut := setupMockDynamoDBForGet(t, getItemFunc)

	result, err := sut.GetPostMetadata(context.Background(), "123")

	if result != nil {
		t.Fatalf("unexpected result: %v", result)
	}

	assert.Equal(t, expectedErr, err)
}

func TestGetPostMetadata_EmptyOutput_ReturnsEmptyOutput(t *testing.T) {
	getItemFunc := func(context context.Context, input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
		return nil, nil
	}

	sut := setupMockDynamoDBForGet(t, getItemFunc)

	result, err := sut.GetPostMetadata(context.Background(), "123")

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

	putItemFunc := func(ctx context.Context, input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
		return &dynamodb.PutItemOutput{
			Attributes: map[string]types.AttributeValue{
				"ID":          &types.AttributeValueMemberS{Value: ID},
				"Title":       &types.AttributeValueMemberS{Value: Title},
				"BodyUrl":     &types.AttributeValueMemberS{Value: BodyUrl},
				"PreviewText": &types.AttributeValueMemberS{Value: PreviewText},
				"Status":      &types.AttributeValueMemberS{Value: string(Status)},
				"CreatedAt":   &types.AttributeValueMemberS{Value: CreatedAt},
				"UpdatedAt":   &types.AttributeValueMemberS{Value: UpdatedAt},
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

	result, err := sut.UpdatePostMetadata(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	assert.Equal(t, expectedResult, result)

}

func TestUpdatePostMetadata_DynamoDBFailure_ReturnsErr(t *testing.T) {
	expectedErr := errors.New("mock error for testing")
	putItemFunc := func(context context.Context, input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
		return nil, errors.New("mock error for testing")
	}
	sut := setupMockDynamoDBForPut(t, putItemFunc)

	input := &model.PostMetadata{}

	result, err := sut.UpdatePostMetadata(context.Background(), input)

	if result != nil {
		t.Fatalf("unexpected result: %v", result)
	}

	assert.Equal(t, expectedErr, err)
}

func TestUpdatePostMetadata_EmptyInput_ReturnsEmpty(t *testing.T) {

	putItemFunc := func(context context.Context, input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
		return nil, nil
	}

	sut := setupMockDynamoDBForPut(t, putItemFunc)

	result, err := sut.UpdatePostMetadata(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	assert.Nil(t, result)
}

func TestListPostMetadata_Success(t *testing.T) {
	ID1 := "123"
	ID2 := "456"
	Title1 := "Title1 Post"
	Title2 := "Title2 Post"
	BodyUrl1 := "http://example.com/bodyText/1"
	BodyUrl2 := "http://example.com/bodyText/2"
	PreviewText1 := "This is a preview1"
	PreviewText2 := "This is a preview2"
	Status := model.Draft
	CreatedAt := time.Now().Format(time.RFC3339)
	UpdatedAt := time.Now().Format(time.RFC3339)

	scanFunc := func(ctx context.Context, input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
		items := []map[string]types.AttributeValue{}
		items = append(items, map[string]types.AttributeValue{
			"ID":          &types.AttributeValueMemberS{Value: ID1},
			"Title":       &types.AttributeValueMemberS{Value: Title1},
			"BodyUrl":     &types.AttributeValueMemberS{Value: BodyUrl1},
			"PreviewText": &types.AttributeValueMemberS{Value: PreviewText1},
			"Status":      &types.AttributeValueMemberS{Value: string(Status)},
			"CreatedAt":   &types.AttributeValueMemberS{Value: CreatedAt},
			"UpdatedAt":   &types.AttributeValueMemberS{Value: UpdatedAt},
		})
		items = append(items, map[string]types.AttributeValue{
			"ID":          &types.AttributeValueMemberS{Value: ID2},
			"Title":       &types.AttributeValueMemberS{Value: Title2},
			"BodyUrl":     &types.AttributeValueMemberS{Value: BodyUrl2},
			"PreviewText": &types.AttributeValueMemberS{Value: PreviewText2},
			"Status":      &types.AttributeValueMemberS{Value: string(Status)},
			"CreatedAt":   &types.AttributeValueMemberS{Value: CreatedAt},
			"UpdatedAt":   &types.AttributeValueMemberS{Value: UpdatedAt},
		})

		consumedCapacity := &types.ConsumedCapacity{}
		lastEvaluatedKey := map[string]types.AttributeValue{}
		return &dynamodb.ScanOutput{
			ConsumedCapacity: consumedCapacity,
			Count:            int32(len(items)),
			Items:            items,
			LastEvaluatedKey: lastEvaluatedKey,
			ScannedCount:     5,
		}, nil
	}
	sut := setupMockDynamoDBForScan(t, scanFunc)

	parsedCreatedAt, _ := time.Parse(time.RFC3339, CreatedAt)
	parsedUpdatedAt, _ := time.Parse(time.RFC3339, UpdatedAt)
	var expectedResult []*model.PostMetadata
	expectedResult = append(expectedResult, &model.PostMetadata{
		ID:          ID1,
		Title:       Title1,
		BodyUrl:     BodyUrl1,
		PreviewText: PreviewText1,
		Status:      Status,
		CreatedAt:   parsedCreatedAt,
		UpdatedAt:   parsedUpdatedAt,
	})
	expectedResult = append(expectedResult, &model.PostMetadata{
		ID:          ID2,
		Title:       Title2,
		BodyUrl:     BodyUrl2,
		PreviewText: PreviewText2,
		Status:      Status,
		CreatedAt:   parsedCreatedAt,
		UpdatedAt:   parsedUpdatedAt,
	})

	result, err := sut.ListPostMetadata(context.Background(), 2, "")

	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	assert.Equal(t, expectedResult, result)
}

func TestListPostMetadata_Succeeds_WithLastEvaluatedKey(t *testing.T) {
	ID1 := "123"
	ID2 := "456"
	Title2 := "Title2 Post"
	BodyUrl2 := "http://example.com/bodyText/2"
	PreviewText2 := "This is a preview2"
	Status := model.Draft
	CreatedAt := time.Now().Format(time.RFC3339)
	UpdatedAt := time.Now().Format(time.RFC3339)

	scanFunc2 := func(ctx context.Context, input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
		items := []map[string]types.AttributeValue{}
		items = append(items, map[string]types.AttributeValue{
			"ID":          &types.AttributeValueMemberS{Value: ID2},
			"Title":       &types.AttributeValueMemberS{Value: Title2},
			"BodyUrl":     &types.AttributeValueMemberS{Value: BodyUrl2},
			"PreviewText": &types.AttributeValueMemberS{Value: PreviewText2},
			"Status":      &types.AttributeValueMemberS{Value: string(Status)},
			"CreatedAt":   &types.AttributeValueMemberS{Value: CreatedAt},
			"UpdatedAt":   &types.AttributeValueMemberS{Value: UpdatedAt},
		})
		consumedCapacity := &types.ConsumedCapacity{}
		lastEvaluatedKey := map[string]types.AttributeValue{}
		return &dynamodb.ScanOutput{
			ConsumedCapacity: consumedCapacity,
			Count:            int32(len(items)),
			Items:            items,
			LastEvaluatedKey: lastEvaluatedKey,
			ScannedCount:     5,
		}, nil
	}

	scanFunc := func(ctx context.Context, input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
		if key, ok := input.ExclusiveStartKey["PrimaryKeyAttribute"]; ok {
			if id, ok := key.(*types.AttributeValueMemberS); ok {
				input.ExclusiveStartKey["PrimaryKeyAttribute"] = id
				return scanFunc2(ctx, input)
			}
		}
		return nil, errors.New("unexpected error")
	}

	sut := setupMockDynamoDBForScan(t, scanFunc)

	parsedCreatedAt, _ := time.Parse(time.RFC3339, CreatedAt)
	parsedUpdatedAt, _ := time.Parse(time.RFC3339, UpdatedAt)
	var expectedResult []*model.PostMetadata
	expectedResult = append(expectedResult, &model.PostMetadata{
		ID:          ID2,
		Title:       Title2,
		BodyUrl:     BodyUrl2,
		PreviewText: PreviewText2,
		Status:      Status,
		CreatedAt:   parsedCreatedAt,
		UpdatedAt:   parsedUpdatedAt,
	})

	result, err := sut.ListPostMetadata(context.Background(), 1, ID1)

	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	assert.Equal(t, expectedResult, result)
}

func TestListPostMetadata_DynamoDBFailure_ReturnsErr(t *testing.T) {
	expectedErr := errors.New("mock error for testing")
	scanFunc := func(ctx context.Context, input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
		return nil, errors.New("mock error for testing")
	}
	sut := setupMockDynamoDBForScan(t, scanFunc)

	result, err := sut.ListPostMetadata(context.Background(), 2, "")

	if result != nil {
		t.Fatalf("unexpected result: %v", result)
	}

	assert.Equal(t, expectedErr, err)

}

func setupMockDynamoDBForGet(t testing.TB, getItemFunc func(context.Context, *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)) PostMetadataDdbDao {
	t.Helper()
	mockDynamoDBClient := &MockDynamoDBClient{
		GetItemFunc: getItemFunc,
	}

	sut := PostMetadataDdbDao{
		client: mockDynamoDBClient,
	}

	return sut
}

func setupMockDynamoDBForPut(t testing.TB, putItemFunc func(context.Context, *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)) PostMetadataDdbDao {
	t.Helper()
	mockDynamoDBClient := &MockDynamoDBClient{
		PutItemFunc: putItemFunc,
	}

	sut := PostMetadataDdbDao{
		client: mockDynamoDBClient,
	}

	return sut
}

func setupMockDynamoDBForScan(t testing.TB, scanFunc func(context.Context, *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)) PostMetadataDdbDao {
	t.Helper()
	MockDynamoDBClient := &MockDynamoDBClient{
		ScanFunc: scanFunc,
	}
	sut := PostMetadataDdbDao{
		client: MockDynamoDBClient,
	}
	return sut
}
