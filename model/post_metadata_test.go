package model

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestToDynamoDbAttributes_Succeeds(t *testing.T) {
	currentTime := time.Now().Truncate(time.Second) // Truncate to seconds for consistency
	post := &PostMetadata{
		ID:          "ID1",
		Title:       "Title1",
		BodyUrl:     "BodyUrl1",
		PreviewText: "PreviewText1",
		Status:      Posted,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}

	expected := map[string]types.AttributeValue{
		"ID":          &types.AttributeValueMemberS{Value: "ID1"},
		"Title":       &types.AttributeValueMemberS{Value: "Title1"},
		"BodyUrl":     &types.AttributeValueMemberS{Value: "BodyUrl1"},
		"PreviewText": &types.AttributeValueMemberS{Value: "PreviewText1"},
		"Status":      &types.AttributeValueMemberS{Value: string(Posted)},
		"CreatedAt":   &types.AttributeValueMemberS{Value: currentTime.Format(time.RFC3339)},
		"UpdatedAt":   &types.AttributeValueMemberS{Value: currentTime.Format(time.RFC3339)},
	}

	result := ToDynamoDbAttributes(post)

	assert.Equal(t, expected, result, "The DynamoDB attribute values should match the expected map")
}

func TestToDynamoDbAttributes_NilPost_ReturnsNil(t *testing.T) {
	result := ToDynamoDbAttributes(nil)
	assert.Nil(t, result, "Expected nil when input post is nil")
}

func TestFromDynamoDBAttributeValue_Succeeds(t *testing.T) {
	currentTime := time.Now().Truncate(time.Second)
	var values []map[string]types.AttributeValue
	ddbValue := map[string]types.AttributeValue{
		"ID":          &types.AttributeValueMemberS{Value: "ID1"},
		"Title":       &types.AttributeValueMemberS{Value: "Title1"},
		"BodyUrl":     &types.AttributeValueMemberS{Value: "BodyUrl1"},
		"PreviewText": &types.AttributeValueMemberS{Value: "PreviewText1"},
		"Status":      &types.AttributeValueMemberS{Value: string(Posted)},
		"CreatedAt":   &types.AttributeValueMemberS{Value: currentTime.Format(time.RFC3339)},
		"UpdatedAt":   &types.AttributeValueMemberS{Value: currentTime.Format(time.RFC3339)},
	}

	values = append(values, ddbValue)
	values = append(values, ddbValue)

	var expectedResult []*PostMetadata
	expected := &PostMetadata{
		ID:          "ID1",
		Title:       "Title1",
		BodyUrl:     "BodyUrl1",
		PreviewText: "PreviewText1",
		Status:      Posted,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}

	expectedResult = append(expectedResult, expected)
	expectedResult = append(expectedResult, expected)

	result := FromDynamoDBAttributeValues(values)

	for i, item := range result {
		assert.Equal(t, expectedResult[i], item, "The PostMetadata content should match")
	}
}
