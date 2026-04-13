//go:build system

package user_test

import (
	"net/http"
	"testing"
	"time"

	stateanalytics "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/analytics"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	statecandidate "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/candidate"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	statefilevault "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/filevault"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestGetMe_Success(t *testing.T) {
	database.Empty(t)

	userEntity := auth.GivenUsers(t, map[string]any{
		"email":          "candidate@example.com",
		"phone_number":   "+998901234567",
		"is_verified":    true,
		"last_login_at":  timePtr(time.Now().Add(-2 * time.Hour)),
		"last_active_at": timePtr(time.Now().Add(-1 * time.Hour)),
	})[0]
	sessionEntity := auth.GivenSessions(t, map[string]any{
		"user_id": userEntity.ID,
	})[0]
	auth.GivenOAuthAccounts(t, map[string]any{
		"user_id":          userEntity.ID,
		"provider":         "google",
		"provider_user_id": "google-user-1",
		"provider_email":   "candidate@example.com",
	})
	profile := statecandidate.GivenProfiles(t, map[string]any{
		"user_id":                 userEntity.ID,
		"full_name":               "John Candidate",
		"bio":                     "Practicing backend interviews",
		"location":                "Tashkent",
		"target_role":             "Golang Backend Developer",
		"experience_level":        "mid",
		"interview_goal_per_week": 4,
		"onboarding_completed":    true,
		"onboarding_completed_at": time.Now().Add(-30 * time.Minute),
	})[0]
	statecandidate.GivenTopicPreferences(t,
		map[string]any{
			"candidate_profile_id": profile.ID,
			"topic_key":            "golang-concurrency",
			"priority":             0,
		},
		map[string]any{
			"candidate_profile_id": profile.ID,
			"topic_key":            "postgres-indexing",
			"priority":             1,
		},
	)
	stateanalytics.GivenProgressSummaries(t, map[string]any{
		"user_id":                  userEntity.ID,
		"current_streak":           2,
		"longest_streak":           5,
		"total_interviews_taken":   8,
		"total_time_spent_seconds": 5400,
		"average_score":            78.5,
		"last_interview_at":        time.Now().Add(-24 * time.Hour),
	})
	avatar := statefilevault.GivenFiles(t, map[string]any{
		"entity_type":      "candidate_profile",
		"entity_id":        profile.ID,
		"association_type": "avatar",
		"sort_order":       1,
	})[0]

	resp := trigger.UserAction(t).GET("/api/v1/auth/get-me").
		WithHeader("Authorization", "Bearer "+sessionEntity.AccessToken).
		Expect()

	resp.Status(http.StatusOK)
	obj := resp.JSON().Object()
	obj.Value("user").Object().Value("id").String().IsEqual(userEntity.ID)
	obj.Value("user").Object().Value("email").String().IsEqual("candidate@example.com")
	obj.Value("user").Object().Value("oauth_providers").Array().Length().IsEqual(1)
	obj.Value("profile").Object().Value("full_name").String().IsEqual("John Candidate")
	obj.Value("profile").Object().Value("preferred_topics").Array().Length().IsEqual(2)
	obj.Value("profile").Object().Value("onboarding_completed").Boolean().IsTrue()
	obj.Value("profile").Object().Value("onboarding_completed_at").String().NotEmpty()
	obj.Value("progress_summary").Object().Value("total_interviews_taken").Number().IsEqual(8)
	obj.Value("avatar").Object().Value("file_id").String().IsEqual(avatar.ID)
	obj.Value("avatar").Object().Value("download_url").String().
		IsEqual("/api/v1/filevault/download?id=" + avatar.ID)

	linkedOAuth := auth.ListOAuthAccountsByUserID(t, userEntity.ID)
	assert.Len(t, linkedOAuth, 1)
}

func TestGetMe_MeAliasSuccess(t *testing.T) {
	database.Empty(t)

	userEntity := auth.GivenUsers(t, map[string]any{
		"email":       "candidate@example.com",
		"is_verified": true,
	})[0]
	sessionEntity := auth.GivenSessions(t, map[string]any{
		"user_id": userEntity.ID,
	})[0]

	resp := trigger.UserAction(t).GET("/api/v1/me").
		WithHeader("Authorization", "Bearer "+sessionEntity.AccessToken).
		Expect()

	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("user").Object().
		Value("id").String().IsEqual(userEntity.ID)
}

func TestGetMe_AuthFailures(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		wantStatus int
		wantErr    string
	}{
		{
			name:       "missing authorization header",
			token:      "",
			wantStatus: http.StatusUnauthorized,
			wantErr:    "UNAUTHORIZED",
		},
		{
			name:       "invalid token",
			token:      "invalid-token",
			wantStatus: http.StatusUnauthorized,
			wantErr:    "UNAUTHORIZED",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			database.Empty(t)

			req := trigger.UserAction(t).GET("/api/v1/auth/get-me")
			if tc.token != "" {
				req = req.WithHeader("Authorization", "Bearer "+tc.token)
			}
			resp := req.Expect()

			resp.Status(tc.wantStatus)
			resp.JSON().Object().Value("error").Object().Value("code").String().IsEqual(tc.wantErr)
		})
	}
}

func timePtr(v time.Time) *time.Time { return &v }
