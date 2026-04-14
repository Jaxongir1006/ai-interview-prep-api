package embassy

import (
	"context"
	"slices"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/question"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/review"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/session"
	interviewportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/interview"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/sorter"
)

func (e *embassy) ListDashboardSessions(
	ctx context.Context,
	req *interviewportal.ListDashboardSessionsRequest,
) (*interviewportal.ListDashboardSessionsResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}

	status := firstStatus(req.Statuses)
	sessions, err := e.domainContainer.SessionRepo().List(ctx, session.Filter{
		UserID:        &req.UserID,
		Status:        status,
		StartedAtFrom: req.StartedAtFrom,
		StartedAtTo:   req.StartedAtTo,
		Limit:         &limit,
		SortOpts: sorter.SortOpts{
			{F: "started_at", D: sorter.Desc},
			{F: "id", D: sorter.Desc},
		},
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	if len(req.Statuses) > 1 {
		sessions = filterByStatuses(sessions, req.Statuses)
	}
	if req.Cursor != nil {
		sessions = filterAfterCursor(sessions, *req.Cursor)
	}

	items := make([]interviewportal.DashboardSession, 0, len(sessions))
	for i := range sessions {
		item, itemErr := e.toDashboardSession(ctx, &sessions[i])
		if itemErr != nil {
			return nil, errx.Wrap(itemErr)
		}
		if req.TopicID != nil && !hasTopic(item.Topics, *req.TopicID) {
			continue
		}
		items = append(items, *item)
	}

	var nextCursor *string
	if len(items) == limit {
		next := items[len(items)-1].ID
		nextCursor = &next
	}

	return &interviewportal.ListDashboardSessionsResponse{
		Items:      items,
		NextCursor: nextCursor,
	}, nil
}

func (e *embassy) toDashboardSession(
	ctx context.Context,
	s *session.Session,
) (*interviewportal.DashboardSession, error) {
	questions, err := e.domainContainer.QuestionRepo().List(ctx, question.Filter{
		SessionID: &s.ID,
		SortOpts: sorter.SortOpts{
			{F: "position", D: sorter.Asc},
			{F: "id", D: sorter.Asc},
		},
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	topics := uniqueTopics(questions)
	questionIDs := make([]int64, 0, len(questions))
	for i := range questions {
		questionIDs = append(questionIDs, questions[i].ID)
	}

	reviewedCount := 0
	if len(questionIDs) > 0 {
		reviews, reviewErr := e.domainContainer.ReviewRepo().List(ctx, review.Filter{
			SessionQuestionIDs: questionIDs,
		})
		if reviewErr != nil {
			return nil, errx.Wrap(reviewErr)
		}
		reviewedCount = len(reviews)
	}

	return &interviewportal.DashboardSession{
		ID:                    s.ID,
		Title:                 s.Title,
		Status:                s.Status,
		Score:                 s.TotalScore,
		StartedAt:             s.StartedAt,
		CompletedAt:           s.CompletedAt,
		DurationSeconds:       s.TotalDurationSeconds,
		QuestionCount:         s.QuestionCount,
		AnsweredCount:         s.AnsweredCount,
		Topics:                topics,
		ReviewedQuestionCount: reviewedCount,
	}, nil
}

func firstStatus(statuses []string) *string {
	if len(statuses) == 1 {
		return &statuses[0]
	}
	return nil
}

func filterByStatuses(items []session.Session, statuses []string) []session.Session {
	out := make([]session.Session, 0, len(items))
	for i := range items {
		if slices.Contains(statuses, items[i].Status) {
			out = append(out, items[i])
		}
	}
	return out
}

func filterAfterCursor(items []session.Session, cursor string) []session.Session {
	for i := range items {
		if items[i].ID == cursor {
			return items[i+1:]
		}
	}
	return items
}

func uniqueTopics(questions []question.Question) []interviewportal.Topic {
	seen := map[string]struct{}{}
	topics := make([]interviewportal.Topic, 0, len(questions))
	for i := range questions {
		if _, ok := seen[questions[i].TopicKey]; ok {
			continue
		}
		seen[questions[i].TopicKey] = struct{}{}
		topics = append(topics, interviewportal.Topic{
			ID:   questions[i].TopicKey,
			Name: questions[i].TopicName,
		})
	}
	return topics
}

func hasTopic(topics []interviewportal.Topic, topicID string) bool {
	for i := range topics {
		if topics[i].ID == topicID {
			return true
		}
	}
	return false
}
