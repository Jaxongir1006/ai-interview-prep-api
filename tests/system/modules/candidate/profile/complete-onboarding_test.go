//go:build system

package profile_test

import (
	"net/http"
	"sort"
	"testing"

	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	statecandidate "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/candidate"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	stateinterview "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/interview"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestCompleteOnboarding_Success(t *testing.T) {
	database.Empty(t)
	stateinterview.GivenDefaultCatalog(t)

	userEntity := auth.GivenUsers(t, map[string]any{
		"email":       "candidate@example.com",
		"is_verified": true,
	})[0]
	sessionEntity := auth.GivenSessions(t, map[string]any{
		"user_id": userEntity.ID,
	})[0]
	profile := statecandidate.GivenProfiles(t, map[string]any{
		"user_id":   userEntity.ID,
		"full_name": "John Doe",
	})[0]
	statecandidate.GivenTopicPreferences(t, map[string]any{
		"candidate_profile_id": profile.ID,
		"topic_key":            "old-topic",
		"priority":             0,
	})

	resp := trigger.UserAction(t).POST("/api/v1/me/complete-onboarding").
		WithHeader("Authorization", "Bearer "+sessionEntity.AccessToken).
		WithJSON(map[string]any{
			"target_role":      "python",
			"experience_level": "junior",
			"preferred_topics": []string{
				"algorithms",
				"system-design",
				"database-design",
			},
		}).
		Expect()

	resp.Status(http.StatusOK)
	obj := resp.JSON().Object().Value("profile").Object()
	obj.Value("id").Number().IsEqual(profile.ID)
	obj.Value("user_id").String().IsEqual(userEntity.ID)
	obj.Value("full_name").String().IsEqual("John Doe")
	obj.Value("target_role").String().IsEqual("python")
	obj.Value("experience_level").String().IsEqual("junior")
	obj.Value("preferred_topics").Array().Length().IsEqual(3)
	obj.Value("onboarding_completed").Boolean().IsTrue()
	obj.Value("onboarding_completed_at").String().NotEmpty()

	updated := statecandidate.GetProfileByUserID(t, userEntity.ID)
	assert.Equal(t, "python", *updated.TargetRole)
	assert.Equal(t, "junior", *updated.ExperienceLevel)
	assert.True(t, updated.OnboardingCompleted)
	assert.NotNil(t, updated.OnboardingCompletedAt)

	preferences := statecandidate.ListTopicPreferencesByProfileID(t, profile.ID)
	sort.Slice(preferences, func(i, j int) bool {
		return preferences[i].Priority < preferences[j].Priority
	})
	assert.Len(t, preferences, 3)
	assert.Equal(t, "algorithms", preferences[0].TopicKey)
	assert.Equal(t, 0, preferences[0].Priority)
	assert.Equal(t, "system-design", preferences[1].TopicKey)
	assert.Equal(t, 1, preferences[1].Priority)
	assert.Equal(t, "database-design", preferences[2].TopicKey)
	assert.Equal(t, 2, preferences[2].Priority)
}

func TestCompleteOnboarding_ProfileNotFound(t *testing.T) {
	database.Empty(t)
	stateinterview.GivenDefaultCatalog(t)

	userEntity := auth.GivenUsers(t, map[string]any{
		"email":       "candidate@example.com",
		"is_verified": true,
	})[0]
	sessionEntity := auth.GivenSessions(t, map[string]any{
		"user_id": userEntity.ID,
	})[0]

	resp := trigger.UserAction(t).POST("/api/v1/me/complete-onboarding").
		WithHeader("Authorization", "Bearer "+sessionEntity.AccessToken).
		WithJSON(validPayload()).
		Expect()

	resp.Status(http.StatusNotFound)
	resp.JSON().Object().Value("error").Object().
		Value("code").String().IsEqual("CANDIDATE_PROFILE_NOT_FOUND")
}

func TestCompleteOnboarding_EmailNotVerified(t *testing.T) {
	database.Empty(t)

	userEntity := auth.GivenUsers(t, map[string]any{
		"email":       "candidate@example.com",
		"is_verified": false,
	})[0]
	sessionEntity := auth.GivenSessions(t, map[string]any{
		"user_id": userEntity.ID,
	})[0]
	statecandidate.GivenProfiles(t, map[string]any{
		"user_id": userEntity.ID,
	})

	resp := trigger.UserAction(t).POST("/api/v1/me/complete-onboarding").
		WithHeader("Authorization", "Bearer "+sessionEntity.AccessToken).
		WithJSON(validPayload()).
		Expect()

	resp.Status(http.StatusBadRequest)
	resp.JSON().Object().Value("error").Object().
		Value("code").String().IsEqual("EMAIL_NOT_VERIFIED")
}

func TestCompleteOnboarding_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		payload map[string]any
	}{
		{
			name: "unknown target role",
			payload: map[string]any{
				"target_role":      "rust",
				"experience_level": "junior",
				"preferred_topics": []string{"algorithms"},
			},
		},
		{
			name: "unknown topic",
			payload: map[string]any{
				"target_role":      "python",
				"experience_level": "junior",
				"preferred_topics": []string{"unknown-topic"},
			},
		},
		{
			name: "duplicate topics",
			payload: map[string]any{
				"target_role":      "python",
				"experience_level": "junior",
				"preferred_topics": []string{"algorithms", "algorithms"},
			},
		},
		{
			name: "missing topics",
			payload: map[string]any{
				"target_role":      "python",
				"experience_level": "junior",
			},
		},
		{
			name: "unknown experience level",
			payload: map[string]any{
				"target_role":      "python",
				"experience_level": "expert",
				"preferred_topics": []string{"algorithms"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			database.Empty(t)
			stateinterview.GivenDefaultCatalog(t)

			userEntity := auth.GivenUsers(t, map[string]any{
				"email":       "candidate@example.com",
				"is_verified": true,
			})[0]
			sessionEntity := auth.GivenSessions(t, map[string]any{
				"user_id": userEntity.ID,
			})[0]
			statecandidate.GivenProfiles(t, map[string]any{
				"user_id": userEntity.ID,
			})

			resp := trigger.UserAction(t).POST("/api/v1/me/complete-onboarding").
				WithHeader("Authorization", "Bearer "+sessionEntity.AccessToken).
				WithJSON(tc.payload).
				Expect()

			resp.Status(http.StatusBadRequest)
		})
	}
}

func TestCompleteOnboarding_AuthFailures(t *testing.T) {
	database.Empty(t)

	resp := trigger.UserAction(t).POST("/api/v1/me/complete-onboarding").
		WithJSON(validPayload()).
		Expect()

	resp.Status(http.StatusUnauthorized)
	resp.JSON().Object().Value("error").Object().
		Value("code").String().IsEqual("UNAUTHORIZED")
}

func validPayload() map[string]any {
	return map[string]any{
		"target_role":      "python",
		"experience_level": "junior",
		"preferred_topics": []string{"algorithms"},
	}
}
