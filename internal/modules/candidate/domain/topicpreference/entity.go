package topicpreference

import "github.com/rise-and-shine/pkg/pg"

const (
	CodeTopicPreferenceNotFound = "TOPIC_PREFERENCE_NOT_FOUND"
)

type TopicPreference struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	CandidateProfileID int64  `json:"candidate_profile_id"`
	TopicKey           string `json:"topic_key"`
	Priority           int    `json:"priority"`
}
