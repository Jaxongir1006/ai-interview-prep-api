package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/review"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewReviewRepo(idb bun.IDB) review.Repo {
	return repogen.NewPgRepoBuilder[review.Review, review.Filter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(review.CodeReviewNotFound).
		WithConflictCodesMap(map[string]string{
			"uk_interview_reviews_session_question_id": review.CodeReviewQuestionConflict,
		}).
		WithFilterFunc(reviewFilterFunc).
		Build()
}

func reviewFilterFunc(q *bun.SelectQuery, f review.Filter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.SessionQuestionID != nil {
		q = q.Where("session_question_id = ?", *f.SessionQuestionID)
	}
	if f.AnswerID != nil {
		q = q.Where("answer_id = ?", *f.AnswerID)
	}
	if f.ReviewerType != nil {
		q = q.Where("reviewer_type = ?", *f.ReviewerType)
	}
	if f.IDs != nil {
		q = q.Where("id IN (?)", bun.In(f.IDs))
	}
	if f.SessionQuestionIDs != nil {
		q = q.Where("session_question_id IN (?)", bun.In(f.SessionQuestionIDs))
	}
	if f.AnswerIDs != nil {
		q = q.Where("answer_id IN (?)", bun.In(f.AnswerIDs))
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
