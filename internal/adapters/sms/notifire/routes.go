package notifire

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseRoutes parses sms-routes from a yaml string (SMS_ROUTES env var).
// Empty input yields a nil slice and no error, so routes are simply not sent.
func ParseRoutes(raw string) ([]RouteSt, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	var routes []RouteSt

	if err := yaml.Unmarshal([]byte(raw), &routes); err != nil {
		return nil, fmt.Errorf("fail to parse sms-routes yaml: %w", err)
	}

	return routes, nil
}
