package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/answer"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewAnswerRepo(idb bun.IDB) answer.Repo {
	return repogen.NewPgRepoBuilder[answer.Answer, answer.Filter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(answer.CodeAnswerNotFound).
		WithConflictCodesMap(map[string]string{
			"uk_interview_answers_session_question_id": answer.CodeAnswerQuestionConflict,
		}).
		WithFilterFunc(answerFilterFunc).
		Build()
}

func answerFilterFunc(q *bun.SelectQuery, f answer.Filter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.SessionQuestionID != nil {
		q = q.Where("session_question_id = ?", *f.SessionQuestionID)
	}
	if f.IDs != nil {
		q = q.Where("id IN (?)", bun.In(f.IDs))
	}
	if f.SessionQuestionIDs != nil {
		q = q.Where("session_question_id IN (?)", bun.In(f.SessionQuestionIDs))
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
