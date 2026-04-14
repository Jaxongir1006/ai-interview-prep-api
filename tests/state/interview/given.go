package interview

import (
	"testing"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/answer"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/question"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/review"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/session"
	interviewpg "github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/infra/postgres"
	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/anymap"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func GivenSessions(t *testing.T, data ...map[string]any) []session.Session {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenSessions: at least one session data map is required")
	}

	db := database.GetTestDB(t)
	repo := interviewpg.NewSessionRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	items := make([]session.Session, 0, len(data))
	now := time.Now().UTC()

	for i, d := range data {
		userID := anymap.String(d, "user_id", "")
		if userID == "" {
			t.Fatalf("GivenSessions[%d]: user_id is required", i)
		}

		item := &session.Session{
			ID:                   anymap.String(d, "id", uuid.NewString()),
			UserID:               userID,
			Title:                anymap.String(d, "title", "Python Backend Interview"),
			TargetRole:           anymap.String(d, "target_role", "python"),
			ExperienceLevel:      anymap.String(d, "experience_level", "mid"),
			Difficulty:           anymap.String(d, "difficulty", session.DifficultyMedium),
			Status:               anymap.String(d, "status", session.StatusCompleted),
			QuestionCount:        cast.ToInt(d["question_count"]),
			AnsweredCount:        cast.ToInt(d["answered_count"]),
			TotalScore:           float64Ptr(d, "total_score", lo.ToPtr(87.0)),
			TotalDurationSeconds: cast.ToInt64(d["total_duration_seconds"]),
			StartedAt:            anymap.Time(d, "started_at", now.Add(-time.Hour)),
			CompletedAt:          anymap.TimePtr(d, "completed_at", lo.ToPtr(now)),
			AbandonedAt:          anymap.TimePtr(d, "abandoned_at", nil),
		}
		if _, ok := d["question_count"]; !ok {
			item.QuestionCount = 5
		}
		if _, ok := d["answered_count"]; !ok {
			item.AnsweredCount = 5
		}
		if _, ok := d["total_duration_seconds"]; !ok {
			item.TotalDurationSeconds = 3600
		}

		created, err := repo.Create(ctx, item)
		if err != nil {
			t.Fatalf("GivenSessions[%d]: failed to create session: %v", i, err)
		}

		items = append(items, *created)
	}

	return items
}

func GivenQuestions(t *testing.T, data ...map[string]any) []question.Question {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenQuestions: at least one question data map is required")
	}

	db := database.GetTestDB(t)
	repo := interviewpg.NewQuestionRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	items := make([]question.Question, 0, len(data))

	for i, d := range data {
		sessionID := anymap.String(d, "session_id", "")
		if sessionID == "" {
			t.Fatalf("GivenQuestions[%d]: session_id is required", i)
		}

		item := &question.Question{
			SessionID:        sessionID,
			TopicKey:         anymap.String(d, "topic_key", "algorithms"),
			TopicName:        anymap.String(d, "topic_name", "Algorithms"),
			Difficulty:       anymap.String(d, "difficulty", session.DifficultyMedium),
			QuestionType:     anymap.String(d, "question_type", question.TypeTechnical),
			QuestionText:     anymap.String(d, "question_text", "Explain a hash map."),
			ExpectedAnswer:   anymap.StringPtr(d, "expected_answer", nil),
			EvaluationRubric: nil,
			Source:           anymap.String(d, "source", question.SourceAI),
			SourceQuestionID: anymap.StringPtr(d, "source_question_id", nil),
			AIProvider:       anymap.StringPtr(d, "ai_provider", nil),
			AIModel:          anymap.StringPtr(d, "ai_model", nil),
			PromptVersion:    anymap.StringPtr(d, "prompt_version", nil),
			Position:         cast.ToInt(d["position"]),
		}
		if _, ok := d["position"]; !ok {
			item.Position = i + 1
		}

		created, err := repo.Create(ctx, item)
		if err != nil {
			t.Fatalf("GivenQuestions[%d]: failed to create question: %v", i, err)
		}

		items = append(items, *created)
	}

	return items
}

func GivenAnswers(t *testing.T, data ...map[string]any) []answer.Answer {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenAnswers: at least one answer data map is required")
	}

	db := database.GetTestDB(t)
	repo := interviewpg.NewAnswerRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	items := make([]answer.Answer, 0, len(data))
	now := time.Now().UTC()

	for i, d := range data {
		questionID := cast.ToInt64(d["session_question_id"])
		if questionID == 0 {
			t.Fatalf("GivenAnswers[%d]: session_question_id is required", i)
		}

		item := &answer.Answer{
			SessionQuestionID: questionID,
			AnswerText:        anymap.String(d, "answer_text", "My answer"),
			TimeSpentSeconds:  cast.ToInt64(d["time_spent_seconds"]),
			SubmittedAt:       anymap.Time(d, "submitted_at", now),
		}
		if _, ok := d["time_spent_seconds"]; !ok {
			item.TimeSpentSeconds = 300
		}

		created, err := repo.Create(ctx, item)
		if err != nil {
			t.Fatalf("GivenAnswers[%d]: failed to create answer: %v", i, err)
		}

		items = append(items, *created)
	}

	return items
}

func GivenReviews(t *testing.T, data ...map[string]any) []review.Review {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenReviews: at least one review data map is required")
	}

	db := database.GetTestDB(t)
	repo := interviewpg.NewReviewRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	items := make([]review.Review, 0, len(data))
	now := time.Now().UTC()

	for i, d := range data {
		questionID := cast.ToInt64(d["session_question_id"])
		if questionID == 0 {
			t.Fatalf("GivenReviews[%d]: session_question_id is required", i)
		}

		item := &review.Review{
			SessionQuestionID: questionID,
			AnswerID:          int64Ptr(d, "answer_id", nil),
			ReviewerType:      anymap.String(d, "reviewer_type", review.ReviewerTypeAI),
			Score:             float64Ptr(d, "score", lo.ToPtr(87.0)),
			CorrectnessRate:   float64Ptr(d, "correctness_rate", lo.ToPtr(0.87)),
			Feedback:          anymap.StringPtr(d, "feedback", nil),
			ReviewedAt:        anymap.Time(d, "reviewed_at", now),
		}

		created, err := repo.Create(ctx, item)
		if err != nil {
			t.Fatalf("GivenReviews[%d]: failed to create review: %v", i, err)
		}

		items = append(items, *created)
	}

	return items
}

func int64Ptr(data map[string]any, key string, defaultVal *int64) *int64 {
	v, ok := data[key]
	if !ok {
		return defaultVal
	}
	if v == nil {
		return nil
	}
	out := cast.ToInt64(v)
	return &out
}

func float64Ptr(data map[string]any, key string, defaultVal *float64) *float64 {
	v, ok := data[key]
	if !ok {
		return defaultVal
	}
	if v == nil {
		return nil
	}
	out := cast.ToFloat64(v)
	return &out
}
