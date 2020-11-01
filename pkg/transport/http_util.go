package transport

import (
	"fmt"

	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
)

// getHTTPAddress returns a sanitized HTTP Address
func getHTTPAddress(cfg configuration.Configuration) string {
	return fmt.Sprintf("%s:%d", cfg.HTTP.Address, cfg.HTTP.Port)
}
