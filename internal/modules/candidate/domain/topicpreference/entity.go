package topicpreference

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

const (
	CodeTopicPreferenceNotFound = "TOPIC_PREFERENCE_NOT_FOUND"
)

type TopicPreference struct {
	bun.BaseModel `bun:"table:candidate_topic_preferences,alias:ctp"`

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	CandidateProfileID int64  `json:"candidate_profile_id"`
	TopicKey           string `json:"topic_key"`
	Priority           int    `json:"priority"`

	CreatedAt time.Time `bun:",nullzero" json:"created_at"`
}

func (m *TopicPreference) BeforeAppendModel(_ context.Context, query bun.Query) error {
	if _, ok := query.(*bun.InsertQuery); ok {
		m.CreatedAt = time.Now()
	}
	return nil
}
