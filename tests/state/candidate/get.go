package candidate

import (
	"testing"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/profile"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/topicpreference"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/infra/postgres"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
)

func GetProfileByUserID(t *testing.T, userID string) *profile.CandidateProfile {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewProfileRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	p, err := repo.Get(ctx, profile.Filter{UserID: &userID})
	if err != nil {
		t.Fatalf("GetProfileByUserID: failed to get profile for user %q: %v", userID, err)
	}

	return p
}

func ListTopicPreferencesByProfileID(t *testing.T, profileID int64) []topicpreference.TopicPreference {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewTopicPreferenceRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	items, err := repo.List(ctx, topicpreference.Filter{CandidateProfileID: &profileID})
	if err != nil {
		t.Fatalf("ListTopicPreferencesByProfileID: failed to list topic preferences for profile %d: %v", profileID, err)
	}

	return items
}
