package dao

import (
	"context"
	"time"

	"github.com/neuralcoral/BlogService/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBAPI interface {
	GetItem(context context.Context, input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	PutItem(context context.Context, input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	Scan(context context.Context, input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
}

type PostMetadataDdbDao struct {
	client    DynamoDBAPI
	tableName string
}

func (dao *PostMetadataDdbDao) GetPostMetadata(context context.Context, id string) (*model.PostMetadata, error) {
	ddbInput := &dynamodb.GetItemInput{
		TableName: aws.String(dao.tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	}
	output, err := dao.client.GetItem(context, ddbInput)
	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, nil
	}

	result := model.FromDynamoDBAttributeValue(output.Item)

	return result, nil
}

func (dao *PostMetadataDdbDao) UpdatePostMetadata(context context.Context, postMetadataToUpdate *model.PostMetadata) (*model.PostMetadata, error) {
	if postMetadataToUpdate == nil {
		return nil, nil
	}
	postMetadataToUpdate.UpdatedAt = time.Now().Truncate(time.Second)

	attributeValueMap := model.ToDynamoDbAttributes(postMetadataToUpdate)

	ddbInput := &dynamodb.PutItemInput{
		TableName: aws.String(dao.tableName),
		Item:      attributeValueMap,
	}

	_, err := dao.client.PutItem(context, ddbInput)

	if err != nil {
		return nil, err
	}

	return postMetadataToUpdate, nil
}

func (dao *PostMetadataDdbDao) ListPostMetadata(context context.Context, limit int, lastEvaluatedKey string) ([]*model.PostMetadata, error) {
	ddbInput := &dynamodb.ScanInput{
		TableName: aws.String(dao.tableName),
		Limit:     aws.Int32(int32(limit)),
	}

	if lastEvaluatedKey != "" {
		ddbInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"PrimaryKeyAttribute": &types.AttributeValueMemberS{Value: lastEvaluatedKey},
		}

	}
	output, err := dao.client.Scan(context, ddbInput)
	if err != nil {
		return nil, err
	}

	result := model.FromDynamoDBAttributeValues(output.Items)
	return result, nil
}

func (dao *PostMetadataDdbDao) CreatePostMetadata(postMetadataToCreate *model.PostMetadata) error {
	return nil
}
