package configuration

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type httpServer struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Endpoint string `json:"endpoint"`
}

func init() {
	viper.SetDefault("http.address", "")
	viper.SetDefault("http.port", 8080)
	viper.SetDefault("http.endpoint", "/live")
}

func (s *httpServer) Load(version string) {
	s.Address = viper.GetString("http.address")
	s.Port = viper.GetInt("http.port")
	s.Endpoint = s.setEndpoint(version)
}

func (s httpServer) setEndpoint(version string) string {
	//	rules
	//	a.	prefix 0. and contains dev stage -> "/alpha" (Use development stage)
	//	b.	prefix > 0 and length equal or above 1 character, contains dev stage -> "/v1/beta" (Use Major version and
	//	dev stage)
	//	c.	prefix > 0 and length equal or above 1 character -> "/v1" (Use Major version, for prod)
	//	d. none of the past cases, default -> "/live"
	stageSlice := strings.Split(version, "-")
	containsStage := len(stageSlice) >= 2

	versionSlice := strings.Split(version, ".")
	containsMajorVer := len(versionSlice) >= 1 && versionSlice[0] != "0"
	switch {
	case strings.HasPrefix(version, "0.") && containsStage:
		return fmt.Sprintf("/%s", stageSlice[1])
	case containsMajorVer && containsStage:
		return fmt.Sprintf("/v%s/%s", versionSlice[0], stageSlice[1])
	case containsMajorVer:
		return fmt.Sprintf("/v%s", versionSlice[0])
	default:
		return viper.GetString("http.endpoint")
	}
}
