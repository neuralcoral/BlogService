package dao

import "github.com/neuralcoral/BlogService/model"

type PostMetadataDao interface {
	GetPostmetadata() (model.PostMetadata, error)
	UpdatePostmetadata(postMetadataToUpdate model.PostMetadata) (model.PostMetadata, error)
	ListPostmetadata(limit int, lastEvaluatedKey string) ([]model.PostMetadata, string, error)
	CreatePostmetadata(model.PostMetadata)
}
