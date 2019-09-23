package default_env_value

import (
	"os"
	"testing"
)

func TestSetDefaultEnvValue(t *testing.T) {
	platform := SetDefaultEnvValue("ENV", "aws")
	err := os.Setenv("ENV", "local")
	if err != nil {
		t.Fatal(err)
	}
	localEnv := os.Getenv("ENV")

	if localEnv == platform {
		t.Errorf("expected environment variable: ENV, to have the value: %s, got %s", localEnv, platform)
	}
}
