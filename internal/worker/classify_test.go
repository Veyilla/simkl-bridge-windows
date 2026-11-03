package worker

import (
	"errors"
	"simklbridge/internal/simkl"
	"simklbridge/internal/supervisor"
	"testing"
)

func TestInitialUnauthorizedAccountRequiresSignInWithoutCrashRetry(t *testing.T) {
	err := actionRequired(&simkl.APIError{Status: 401, Code: "user_token_required"})
	if !errors.Is(err, supervisor.ErrAuthentication) {
		t.Fatalf("Anmeldefehler als Absturz behandelt: %v", err)
	}
}
