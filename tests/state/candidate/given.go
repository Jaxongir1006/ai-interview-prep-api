package candidate

import (
	"testing"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/profile"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/topicpreference"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/infra/postgres"
	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/anymap"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	"github.com/spf13/cast"
)

func GivenProfiles(t *testing.T, data ...map[string]any) []profile.CandidateProfile {
	t.Helper()

	if len(data) == 0 {
		data = []map[string]any{{}}
	}

	db := database.GetTestDB(t)
	repo := postgres.NewProfileRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	items := make([]profile.CandidateProfile, 0, len(data))

	for i, d := range data {
		userID := anymap.String(d, "user_id", "")
		if userID == "" {
			t.Fatalf("GivenProfiles[%d]: user_id is required", i)
		}

		item := &profile.CandidateProfile{
			UserID:                userID,
			FullName:              anymap.StringPtr(d, "full_name", nil),
			Bio:                   anymap.StringPtr(d, "bio", nil),
			Location:              anymap.StringPtr(d, "location", nil),
			TargetRole:            anymap.StringPtr(d, "target_role", nil),
			ExperienceLevel:       anymap.StringPtr(d, "experience_level", nil),
			InterviewGoalPerWeek:  cast.ToInt(d["interview_goal_per_week"]),
			OnboardingCompleted:   anymap.Bool(d, "onboarding_completed", false),
			OnboardingCompletedAt: anymap.TimePtr(d, "onboarding_completed_at", nil),
		}
		if _, ok := d["interview_goal_per_week"]; !ok {
			item.InterviewGoalPerWeek = 3
		}

		created, err := repo.Create(ctx, item)
		if err != nil {
			t.Fatalf("GivenProfiles[%d]: failed to create profile: %v", i, err)
		}

		items = append(items, *created)
	}

	return items
}

func GivenTopicPreferences(t *testing.T, data ...map[string]any) []topicpreference.TopicPreference {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenTopicPreferences: at least one topic preference data map is required")
	}

	db := database.GetTestDB(t)
	repo := postgres.NewTopicPreferenceRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	items := make([]topicpreference.TopicPreference, 0, len(data))

	for i, d := range data {
		profileID := cast.ToInt64(d["candidate_profile_id"])
		topicKey := anymap.String(d, "topic_key", "")
		if profileID == 0 || topicKey == "" {
			t.Fatalf("GivenTopicPreferences[%d]: candidate_profile_id and topic_key are required", i)
		}

		item := &topicpreference.TopicPreference{
			CandidateProfileID: profileID,
			TopicKey:           topicKey,
			Priority:           cast.ToInt(d["priority"]),
		}

		created, err := repo.Create(ctx, item)
		if err != nil {
			t.Fatalf("GivenTopicPreferences[%d]: failed to create topic preference: %v", i, err)
		}

		items = append(items, *created)
	}

	return items
}
