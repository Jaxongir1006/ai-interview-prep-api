//go:build system

package catalog_test

import (
	"net/http"
	"testing"

	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	stateinterview "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/interview"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"
)

func TestGetOnboardingOptions_Success(t *testing.T) {
	database.Empty(t)
	stateinterview.GivenDefaultCatalog(t)

	userEntity := auth.GivenUsers(t, map[string]any{
		"email":       "candidate@example.com",
		"is_verified": true,
	})[0]
	sessionEntity := auth.GivenSessions(t, map[string]any{
		"user_id": userEntity.ID,
	})[0]

	resp := trigger.UserAction(t).GET("/api/v1/interview/get-onboarding-options").
		WithHeader("Authorization", "Bearer "+sessionEntity.AccessToken).
		Expect()

	resp.Status(http.StatusOK)
	obj := resp.JSON().Object()
	obj.Value("target_roles").Array().Length().IsEqual(4)
	obj.Value("target_roles").Array().Element(0).Object().Value("key").String().IsEqual("python")
	obj.Value("experience_levels").Array().Length().IsEqual(3)
	obj.Value("experience_levels").Array().Element(0).Object().Value("key").String().IsEqual("junior")
	obj.Value("topics").Array().Length().IsEqual(3)
	obj.Value("topics").Array().Element(0).Object().Value("key").String().IsEqual("algorithms")
	obj.Value("topics").Array().Element(0).Object().Value("target_role_keys").Array().Length().IsEqual(4)
}

func TestGetOnboardingOptions_AuthFailures(t *testing.T) {
	database.Empty(t)

	resp := trigger.UserAction(t).GET("/api/v1/interview/get-onboarding-options").
		Expect()

	resp.Status(http.StatusUnauthorized)
	resp.JSON().Object().Value("error").Object().
		Value("code").String().IsEqual("UNAUTHORIZED")
}
