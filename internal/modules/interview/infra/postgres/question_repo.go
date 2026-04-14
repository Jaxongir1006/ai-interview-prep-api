package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/question"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewQuestionRepo(idb bun.IDB) question.Repo {
	return repogen.NewPgRepoBuilder[question.Question, question.Filter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(question.CodeQuestionNotFound).
		WithConflictCodesMap(map[string]string{
			"uk_interview_questions_session_position": question.CodeQuestionPositionConflict,
		}).
		WithFilterFunc(questionFilterFunc).
		Build()
}

func questionFilterFunc(q *bun.SelectQuery, f question.Filter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.SessionID != nil {
		q = q.Where("session_id = ?", *f.SessionID)
	}
	if f.TopicKey != nil {
		q = q.Where("topic_key = ?", *f.TopicKey)
	}
	if f.Difficulty != nil {
		q = q.Where("difficulty = ?", *f.Difficulty)
	}
	if f.QuestionType != nil {
		q = q.Where("question_type = ?", *f.QuestionType)
	}
	if f.Source != nil {
		q = q.Where("source = ?", *f.Source)
	}
	if f.IDs != nil {
		q = q.Where("id IN (?)", bun.In(f.IDs))
	}
	if f.SessionIDs != nil {
		q = q.Where("session_id IN (?)", bun.In(f.SessionIDs))
	}
	if f.TopicKeys != nil {
		q = q.Where("topic_key IN (?)", bun.In(f.TopicKeys))
	}
	if f.Limit != nil {
		q = q.Limit(*f.Limit)
	}
	if f.Offset != nil {
		q = q.Offset(*f.Offset)
	}
	for _, o := range f.SortOpts {
		q = q.Order(o.ToSQL())
	}
	return q
}
