//go:build system

package dashboard_test

import (
	"net/http"
	"testing"
	"time"

	stateanalytics "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/analytics"
	stateauth "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	statecandidate "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/candidate"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	stateinterview "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/interview"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"
)

func TestDashboardOverview_Success(t *testing.T) {
	database.Empty(t)

	token := seedDashboardState(t)

	resp := trigger.UserAction(t).GET("/api/v1/dashboard/overview").
		WithQuery("range", "7d").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	resp.Status(http.StatusOK)
	obj := resp.JSON().Object()
	obj.Value("user").Object().Value("full_name").String().IsEqual("John Developer")
	obj.Value("user").Object().Value("email").String().IsEqual("john@example.com")
	obj.Value("stats").Object().Value("total_interviews").Object().Value("value").Number().IsEqual(1)
	obj.Value("stats").Object().Value("average_score").Object().Value("value").Number().IsEqual(87)
	obj.Value("performance").Object().Value("points").Array().Length().IsEqual(1)
	obj.Value("topics").Object().Value("items").Array().Length().IsEqual(2)
	obj.Value("topics").Object().Value("weak").Array().Length().IsEqual(1)
	obj.Value("recent_activity").Object().Value("items").Array().Length().IsEqual(1)
	obj.Value("recommendations").Object().Value("recommended_topics").Array().Length().IsEqual(1)
}

func TestDashboardSmallerEndpoints_Success(t *testing.T) {
	database.Empty(t)

	token := seedDashboardState(t)

	stats := trigger.UserAction(t).GET("/api/v1/dashboard/stats").
		WithQuery("range", "30d").
		WithHeader("Authorization", "Bearer "+token).
		Expect()
	stats.Status(http.StatusOK)
	stats.JSON().Object().Value("range").String().IsEqual("30d")
	stats.JSON().Object().Value("stats").Object().Value("total_interviews").Object().Value("value").Number().IsEqual(1)

	trend := trigger.UserAction(t).GET("/api/v1/dashboard/performance-trend").
		WithQuery("range", "30d").
		WithQuery("topic_id", "algorithms").
		WithHeader("Authorization", "Bearer "+token).
		Expect()
	trend.Status(http.StatusOK)
	trend.JSON().Object().Value("topic").Object().Value("id").String().IsEqual("algorithms")
	trend.JSON().Object().Value("points").Array().Length().IsEqual(1)

	topics := trigger.UserAction(t).GET("/api/v1/dashboard/topics").
		WithQuery("range", "30d").
		WithHeader("Authorization", "Bearer "+token).
		Expect()
	topics.Status(http.StatusOK)
	topics.JSON().Object().Value("items").Array().Length().IsEqual(2)

	activity := trigger.UserAction(t).GET("/api/v1/dashboard/recent-activity").
		WithQuery("limit", "10").
		WithHeader("Authorization", "Bearer "+token).
		Expect()
	activity.Status(http.StatusOK)
	activity.JSON().Object().Value("items").Array().Length().IsEqual(1)
	activity.JSON().Object().Value("items").Array().Element(0).Object().
		Value("status").String().IsEqual("completed")

	recommendations := trigger.UserAction(t).GET("/api/v1/dashboard/recommendations").
		WithHeader("Authorization", "Bearer "+token).
		Expect()
	recommendations.Status(http.StatusOK)
	recommendations.JSON().Object().Value("recommended_topics").Array().Length().IsEqual(1)
	recommendations.JSON().Object().Value("next_interview").Object().Value("difficulty").String().IsEqual("medium")
}

func TestDashboardOverview_EmptyState(t *testing.T) {
	database.Empty(t)

	user := stateauth.GivenUsers(t, map[string]any{
		"username":    nil,
		"email":       "empty@example.com",
		"is_verified": true,
	})[0]
	session := stateauth.GivenSessions(t, map[string]any{"user_id": user.ID})[0]
	statecandidate.GivenProfiles(t, map[string]any{
		"user_id":          user.ID,
		"full_name":        "Empty User",
		"target_role":      "python",
		"experience_level": "junior",
	})

	resp := trigger.UserAction(t).GET("/api/v1/dashboard/overview").
		WithHeader("Authorization", "Bearer "+session.AccessToken).
		Expect()

	resp.Status(http.StatusOK)
	obj := resp.JSON().Object()
	obj.Value("stats").Object().Value("total_interviews").Object().Value("value").Number().IsEqual(0)
	obj.Value("stats").Object().Value("average_score").Object().Value("value").Null()
	obj.Value("performance").Object().Value("points").Array().Length().IsEqual(0)
	obj.Value("topics").Object().Value("items").Array().Length().IsEqual(0)
	obj.Value("recent_activity").Object().Value("items").Array().Length().IsEqual(0)
}

func TestDashboardValidationAndAuthFailures(t *testing.T) {
	database.Empty(t)

	token := seedDashboardState(t)

	invalidRange := trigger.UserAction(t).GET("/api/v1/dashboard/stats").
		WithQuery("range", "bad").
		WithHeader("Authorization", "Bearer "+token).
		Expect()
	invalidRange.Status(http.StatusBadRequest)
	invalidRange.JSON().Object().Value("error").Object().
		Value("code").String().IsEqual("VALIDATION_FAILED")

	invalidLimit := trigger.UserAction(t).GET("/api/v1/dashboard/recent-activity").
		WithQuery("limit", "100").
		WithHeader("Authorization", "Bearer "+token).
		Expect()
	invalidLimit.Status(http.StatusBadRequest)
	invalidLimit.JSON().Object().Value("error").Object().
		Value("code").String().IsEqual("VALIDATION_FAILED")

	unauthorized := trigger.UserAction(t).GET("/api/v1/dashboard/overview").Expect()
	unauthorized.Status(http.StatusUnauthorized)
	unauthorized.JSON().Object().Value("error").Object().
		Value("code").String().IsEqual("UNAUTHORIZED")
}

func seedDashboardState(t *testing.T) string {
	t.Helper()

	user := stateauth.GivenUsers(t, map[string]any{
		"username":    nil,
		"email":       "john@example.com",
		"is_verified": true,
	})[0]
	authSession := stateauth.GivenSessions(t, map[string]any{"user_id": user.ID})[0]

	profile := statecandidate.GivenProfiles(t, map[string]any{
		"user_id":          user.ID,
		"full_name":        "John Developer",
		"target_role":      "python",
		"experience_level": "mid",
	})[0]
	statecandidate.GivenTopicPreferences(t,
		map[string]any{"candidate_profile_id": profile.ID, "topic_key": "algorithms", "priority": 0},
		map[string]any{"candidate_profile_id": profile.ID, "topic_key": "security", "priority": 1},
	)

	now := time.Now().UTC()
	stateanalytics.GivenProgressSummaries(t, map[string]any{
		"user_id":                  user.ID,
		"current_streak":           7,
		"longest_streak":           7,
		"total_interviews_taken":   1,
		"total_time_spent_seconds": 3600,
		"average_score":            87,
		"last_interview_at":        now,
	})
	stateanalytics.GivenTopicStats(t,
		map[string]any{
			"user_id":                  user.ID,
			"topic_key":                "security",
			"attempts":                 12,
			"total_time_spent_seconds": 3600,
			"average_score":            65,
			"best_score":               80,
		},
		map[string]any{
			"user_id":                  user.ID,
			"topic_key":                "api_design",
			"attempts":                 20,
			"total_time_spent_seconds": 7200,
			"average_score":            90,
			"best_score":               96,
		},
	)

	interviewSession := stateinterview.GivenSessions(t, map[string]any{
		"user_id":                user.ID,
		"status":                 "completed",
		"total_score":            87.0,
		"started_at":             now.Add(-2 * time.Hour),
		"completed_at":           now.Add(-time.Hour),
		"total_duration_seconds": 3600,
	})[0]
	stateinterview.GivenQuestions(t,
		map[string]any{
			"session_id": interviewSession.ID,
			"topic_key":  "algorithms",
			"topic_name": "Algorithms",
			"position":   1,
		},
		map[string]any{
			"session_id": interviewSession.ID,
			"topic_key":  "security",
			"topic_name": "Security",
			"position":   2,
		},
	)

	return authSession.AccessToken
}
