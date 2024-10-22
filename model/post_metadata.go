package model

import "time"

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
