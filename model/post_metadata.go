package model

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Status string

const (
	Draft  Status = "DRAFT"
	Posted Status = "POSTED"
)

type PostMetadata struct {
	ID          string
	Title       string
	BodyUrl     string
	PreviewText string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Tags        []Tag
}

func ToDynamoDbAttributes(post *PostMetadata) map[string]types.AttributeValue {
	if post == nil {
		return nil
	}

	return map[string]types.AttributeValue{
		"ID":          &types.AttributeValueMemberS{Value: post.ID},
		"Title":       &types.AttributeValueMemberS{Value: post.Title},
		"BodyUrl":     &types.AttributeValueMemberS{Value: post.BodyUrl},
		"PreviewText": &types.AttributeValueMemberS{Value: post.PreviewText},
		"Status":      &types.AttributeValueMemberS{Value: string(post.Status)},
		"CreatedAt":   &types.AttributeValueMemberS{Value: post.CreatedAt.Format(time.RFC3339)},
		"UpdatedAt":   &types.AttributeValueMemberS{Value: post.UpdatedAt.Format(time.RFC3339)},
	}
}

func FromDynamoDBAttributeValues(ddbValues []map[string]types.AttributeValue) []*PostMetadata {
	var result []*PostMetadata
	for _, ddbValue := range ddbValues {
		result = append(result, FromDynamoDBAttributeValue(ddbValue))
	}
	return result
}

func FromDynamoDBAttributeValue(ddbValue map[string]types.AttributeValue) *PostMetadata {
	return &PostMetadata{
		ID:          getStringAttribute(ddbValue["ID"]),
		Title:       getStringAttribute(ddbValue["Title"]),
		BodyUrl:     getStringAttribute(ddbValue["BodyUrl"]),
		PreviewText: getStringAttribute(ddbValue["PreviewText"]),
		Status:      Status(getStringAttribute(ddbValue["Status"])),
		CreatedAt:   parseTime(getStringAttribute(ddbValue["CreatedAt"])),
		UpdatedAt:   parseTime(getStringAttribute(ddbValue["UpdatedAt"])),
	}
}

func getStringAttribute(attr types.AttributeValue) string {
	if attrS, ok := attr.(*types.AttributeValueMemberS); ok {
		return attrS.Value
	}
	return ""
}

func parseTime(timeStr string) time.Time {
	parsedTime, _ := time.Parse(time.RFC3339, timeStr)
	return parsedTime
}
