package analytics

import (
	"testing"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain/progress"
	analyticspg "github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/infra/postgres"
	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/anymap"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	"github.com/spf13/cast"
)

func GivenProgressSummaries(t *testing.T, data ...map[string]any) []progress.Summary {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenProgressSummaries: at least one progress summary data map is required")
	}

	db := database.GetTestDB(t)
	repo := analyticspg.NewProgressSummaryRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	items := make([]progress.Summary, 0, len(data))

	for i, d := range data {
		userID := anymap.String(d, "user_id", "")
		if userID == "" {
			t.Fatalf("GivenProgressSummaries[%d]: user_id is required", i)
		}

		item := &progress.Summary{
			UserID:                userID,
			CurrentStreak:         cast.ToInt(d["current_streak"]),
			LongestStreak:         cast.ToInt(d["longest_streak"]),
			TotalInterviewsTaken:  cast.ToInt(d["total_interviews_taken"]),
			TotalTimeSpentSeconds: cast.ToInt64(d["total_time_spent_seconds"]),
			AverageScore:          cast.ToFloat64(d["average_score"]),
			LastInterviewAt:       anymap.TimePtr(d, "last_interview_at", (*time.Time)(nil)),
		}

		created, err := repo.Create(ctx, item)
		if err != nil {
			t.Fatalf("GivenProgressSummaries[%d]: failed to create progress summary: %v", i, err)
		}

		items = append(items, *created)
	}

	return items
}

func GivenTopicStats(t *testing.T, data ...map[string]any) []progress.TopicStat {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenTopicStats: at least one topic stat data map is required")
	}

	db := database.GetTestDB(t)
	repo := analyticspg.NewTopicStatRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	items := make([]progress.TopicStat, 0, len(data))

	for i, d := range data {
		userID := anymap.String(d, "user_id", "")
		topicKey := anymap.String(d, "topic_key", "")
		if userID == "" || topicKey == "" {
			t.Fatalf("GivenTopicStats[%d]: user_id and topic_key are required", i)
		}

		item := &progress.TopicStat{
			UserID:                userID,
			TopicKey:              topicKey,
			Attempts:              cast.ToInt(d["attempts"]),
			TotalTimeSpentSeconds: cast.ToInt64(d["total_time_spent_seconds"]),
			AverageScore:          cast.ToFloat64(d["average_score"]),
			BestScore:             cast.ToFloat64(d["best_score"]),
			LastPracticedAt:       anymap.TimePtr(d, "last_practiced_at", (*time.Time)(nil)),
		}

		created, err := repo.Create(ctx, item)
		if err != nil {
			t.Fatalf("GivenTopicStats[%d]: failed to create topic stat: %v", i, err)
		}

		items = append(items, *created)
	}

	return items
}
