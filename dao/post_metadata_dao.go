package dao

import "github.com/neuralcoral/BlogService/model"

type PostMetadataDao interface {
	GetPostMetadata(id string) (*model.PostMetadata, error)
	UpdatePostMetadata(postMetadataToUpdate *model.PostMetadata) (*model.PostMetadata, error)
	ListPostMetadata(limit int, lastEvaluatedKey string) ([]*model.PostMetadata, string, error)
	CreatePostMetadata(postMetadataToCreate *model.PostMetadata) error
}
