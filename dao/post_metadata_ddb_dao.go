package dao

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/neuralcoral/BlogService/model"
)

type DynamoDBAPI interface {
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
}

type PostMetadataDdbDao struct {
	client    DynamoDBAPI
	tableName string
}

func (dao *PostMetadataDdbDao) GetPostMetadata(id string) (*model.PostMetadata, error) {
	ddb_input := &dynamodb.GetItemInput{
		TableName: aws.String(dao.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
	}
	output, err := dao.client.GetItem(ddb_input)
	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, nil
	}

	result := &model.PostMetadata{
		ID:          id,
		Title:       *output.Item["Title"].S,
		BodyUrl:     *output.Item["BodyUrl"].S,
		PreviewText: *output.Item["PreviewText"].S,
		Status:      model.Status(*output.Item["Status"].S),
		CreatedAt:   parseTime(output.Item["CreatedAt"].S),
		UpdatedAt:   parseTime(output.Item["UpdatedAt"].S),
	}

	return result, nil
}

func (dao *PostMetadataDdbDao) UpdatePostMetadata(postMetadataToUpdate *model.PostMetadata) (*model.PostMetadata, error) {
	if postMetadataToUpdate == nil {
		return nil, nil
	}
	postMetadataToUpdate.UpdatedAt = time.Now().Truncate(time.Second)

	attributeValueMap := convertPostMetadataToDynamoDBAttributes(postMetadataToUpdate)

	ddb_input := &dynamodb.PutItemInput{
		TableName: aws.String(dao.tableName), Item: attributeValueMap,
	}

	_, err := dao.client.PutItem(ddb_input)

	if err != nil {
		return nil, err
	}

	return postMetadataToUpdate, nil
}

func (dao *PostMetadataDdbDao) ListPostMetadata(limit int, lastEvaluatedKey string) ([]*model.PostMetadata, error) {
	return nil, nil
}

func (dao *PostMetadataDdbDao) CreatePostMetadata(postMetadataToCreate *model.PostMetadata) error {
	return nil
}

func convertPostMetadataToDynamoDBAttributes(post *model.PostMetadata) map[string]*dynamodb.AttributeValue {
	if post == nil {
		return nil
	}

	return map[string]*dynamodb.AttributeValue{
		"ID":          {S: aws.String(post.ID)},
		"Title":       {S: aws.String(post.Title)},
		"BodyUrl":     {S: aws.String(post.BodyUrl)},
		"PreviewText": {S: aws.String(post.PreviewText)},
		"Status":      {S: aws.String(string(post.Status))},
		"CreatedAt":   {S: aws.String(post.CreatedAt.Format(time.RFC3339))},
		"UpdatedAt":   {S: aws.String(post.UpdatedAt.Format(time.RFC3339))},
	}

}

func parseTime(timeStr *string) time.Time {
	parsedTime, _ := time.Parse(time.RFC3339, *timeStr)
	return parsedTime
}
