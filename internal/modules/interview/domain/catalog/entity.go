package catalog

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

const (
	CodeTargetRoleNotFound      = "INTERVIEW_TARGET_ROLE_NOT_FOUND"
	CodeExperienceLevelNotFound = "INTERVIEW_EXPERIENCE_LEVEL_NOT_FOUND"
	CodeTopicNotFound           = "INTERVIEW_TOPIC_NOT_FOUND"
)

type TargetRole struct {
	bun.BaseModel `bun:"table:interview_target_roles,alias:itr"`

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	Key         string  `json:"key"`
	Name        string  `json:"name"`
	Description *string `json:"description"`

	DisplayOrder int  `json:"display_order"`
	IsActive     bool `json:"is_active"`

	CreatedAt time.Time `json:"created_at" bun:",nullzero"`
	UpdatedAt time.Time `json:"updated_at" bun:",nullzero"`
}

func (m *TargetRole) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}

type ExperienceLevel struct {
	bun.BaseModel `bun:"table:interview_experience_levels,alias:iel"`

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	Key         string  `json:"key"`
	Name        string  `json:"name"`
	Description *string `json:"description"`

	DisplayOrder int  `json:"display_order"`
	IsActive     bool `json:"is_active"`

	CreatedAt time.Time `json:"created_at" bun:",nullzero"`
	UpdatedAt time.Time `json:"updated_at" bun:",nullzero"`
}

func (m *ExperienceLevel) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}

type Topic struct {
	bun.BaseModel `bun:"table:interview_topics,alias:it"`

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	Key         string  `json:"key"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Category    *string `json:"category"`

	DisplayOrder int  `json:"display_order"`
	IsActive     bool `json:"is_active"`

	CreatedAt time.Time `json:"created_at" bun:",nullzero"`
	UpdatedAt time.Time `json:"updated_at" bun:",nullzero"`
}

func (m *Topic) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}

type TargetRoleTopic struct {
	bun.BaseModel `bun:"table:interview_target_role_topics,alias:itrtt"`

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	TargetRoleID int64 `json:"target_role_id"`
	TopicID      int64 `json:"topic_id"`

	DisplayOrder int  `json:"display_order"`
	IsActive     bool `json:"is_active"`

	CreatedAt time.Time `json:"created_at" bun:",nullzero"`
	UpdatedAt time.Time `json:"updated_at" bun:",nullzero"`
}

func (m *TargetRoleTopic) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
