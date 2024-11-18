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
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
}

type PostMetadataDdbDao struct {
	client    DynamoDBAPI
	tableName string
}

func (dao *PostMetadataDdbDao) GetPostMetadata(id string) (*model.PostMetadata, error) {
	ddbInput := &dynamodb.GetItemInput{
		TableName: aws.String(dao.tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	}
	output, err := dao.client.GetItem(context.TODO(), ddbInput)
	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, nil
	}

	result := model.FromDynamoDBAttributeValue(output.Item)

	return result, nil
}

func (dao *PostMetadataDdbDao) UpdatePostMetadata(postMetadataToUpdate *model.PostMetadata) (*model.PostMetadata, error) {
	if postMetadataToUpdate == nil {
		return nil, nil
	}
	postMetadataToUpdate.UpdatedAt = time.Now().Truncate(time.Second)

	attributeValueMap := model.ToDynamoDbAttributes(postMetadataToUpdate)

	ddb_input := &dynamodb.PutItemInput{
		TableName: aws.String(dao.tableName),
		Item:      attributeValueMap,
	}

	_, err := dao.client.PutItem(ddb_input)

	if err != nil {
		return nil, err
	}

	return postMetadataToUpdate, nil
}

func (dao *PostMetadataDdbDao) ListPostMetadata(limit int, lastEvaluatedKey string) ([]*model.PostMetadata, error) {
	ddbInput := &dynamodb.ScanInput{
		TableName: aws.String(dao.tableName),
		Limit:     aws.Int64(int64(limit)),
	}

	if lastEvaluatedKey != "" {
		ddbInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"PrimaryKeyAttribute": &types.AttributeValueMemberS{Value: lastEvaluatedKey},
		}

	}
	output, err := dao.client.Scan(ddbInput)
	if err != nil {
		return nil, err
	}

	result := model.FromDynamoDBAttributeValues(output.Items)
	return result, nil
}

func (dao *PostMetadataDdbDao) CreatePostMetadata(postMetadataToCreate *model.PostMetadata) error {
	return nil
}
