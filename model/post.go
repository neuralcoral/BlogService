package model

import "time"

type Status string

const (
	Draft  Status = "DRAFT"
	Posted Status = "POSTED"
)

type PostMetadata struct {
	id          string
	title       string
	bodyUrl     string
	previewText string
	status      Status
	createdAt   time.Time
	updatedAt   time.Time
	tags        []Tag
}

type Post struct {
	postMetadata PostMetadata
	content      string
}
