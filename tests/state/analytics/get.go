package analytics

import (
	"testing"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain/progress"
	analyticspg "github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/infra/postgres"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
)

func GetProgressSummaryByUserID(t *testing.T, userID string) *progress.Summary {
	t.Helper()

	db := database.GetTestDB(t)
	repo := analyticspg.NewProgressSummaryRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	summary, err := repo.Get(ctx, progress.SummaryFilter{UserID: &userID})
	if err != nil {
		t.Fatalf("GetProgressSummaryByUserID: failed to get summary for user %q: %v", userID, err)
	}

	return summary
}
