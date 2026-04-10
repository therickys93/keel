package shoutrrr

import (
	"fmt"
	"os"
	"strings"

	shoutrrrLib "github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	shoutrrrTypes "github.com/containrrr/shoutrrr/pkg/types"
	"github.com/keel-hq/keel/constants"
	"github.com/keel-hq/keel/extension/notification"
	"github.com/keel-hq/keel/types"

	log "github.com/sirupsen/logrus"
)

type sender struct {
	router *router.ServiceRouter
	uiURL  string
}

func init() {
	notification.RegisterSender("shoutrrr", &sender{})
}

func (s *sender) Configure(config *notification.Config) (bool, error) {
	raw := os.Getenv(constants.EnvShoutrrrUrl)
	if raw == "" {
		return false, nil
	}

	var urls []string
	for _, u := range strings.Split(raw, ",") {
		if u = strings.TrimSpace(u); u != "" {
			urls = append(urls, u)
		}
	}
	if len(urls) == 0 {
		return false, nil
	}

	r, err := shoutrrrLib.CreateSender(urls...)
	if err != nil {
		return false, fmt.Errorf("shoutrrr: %w", err)
	}
	s.router = r
	s.uiURL = os.Getenv("KEEL_UI_URL")

	log.WithFields(log.Fields{
		"name":  "shoutrrr",
		"count": len(urls),
	}).Info("extension.notification.shoutrrr: sender configured")
	return true, nil
}

func levelPriority(l types.Level) string {
	switch l {
	case types.LevelFatal:
		return "5"
	case types.LevelError:
		return "4"
	case types.LevelWarn:
		return "3"
	case types.LevelSuccess:
		return "2"
	default:
		return "1"
	}
}

func levelTags(l types.Level) string {
	switch l {
	case types.LevelFatal:
		return "rotating_light"
	case types.LevelError:
		return "x"
	case types.LevelWarn:
		return "warning"
	case types.LevelSuccess:
		return "rocket"
	default:
		return "information_source"
	}
}

func (s *sender) Send(event types.EventNotification) error {
	title := fmt.Sprintf("[%s] %s", strings.ToUpper(event.Level.String()), event.Name)
	body := event.Message
	if event.Identifier != "" {
		body += "\n" + event.Identifier
	}

	params := shoutrrrTypes.Params{
		"title":    title,
		// "priority": levelPriority(event.Level),
		// "tags":     levelTags(event.Level),
		// "icon":     constants.KeelLogoURL,
	}

	if s.uiURL != "" {
		params["click"] = s.uiURL
	}

	errs := s.router.Send(body, &params)

	var failed []string
	for _, err := range errs {
		if err != nil {
			failed = append(failed, err.Error())
		}
	}
	if len(failed) > 0 {
		return fmt.Errorf("shoutrrr: %s", strings.Join(failed, "; "))
	}
	return nil
}
