package shoutrrr

import (
	"os"
	"testing"

	"github.com/keel-hq/keel/constants"
	"github.com/keel-hq/keel/extension/notification"
)

func TestConfigureDisabledWhenNoEnv(t *testing.T) {
	os.Unsetenv(constants.EnvShoutrrrUrl)
	s := &sender{}
	enabled, err := s.Configure(&notification.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enabled {
		t.Fatal("expected sender to be disabled when env is not set")
	}
}

func TestConfigureEmptyAfterSplit(t *testing.T) {
	os.Setenv(constants.EnvShoutrrrUrl, " , , ")
	defer os.Unsetenv(constants.EnvShoutrrrUrl)

	s := &sender{}
	enabled, err := s.Configure(&notification.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enabled {
		t.Fatal("expected sender to be disabled when all URLs are empty")
	}
}

func TestConfigureWithInvalidURL(t *testing.T) {
	os.Setenv(constants.EnvShoutrrrUrl, "not-a-valid-url")
	defer os.Unsetenv(constants.EnvShoutrrrUrl)

	s := &sender{}
	enabled, err := s.Configure(&notification.Config{})
	if enabled {
		t.Fatal("expected sender to be disabled for invalid URL")
	}
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
