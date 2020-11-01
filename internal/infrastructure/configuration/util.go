package configuration

import "strings"

// setResourceStage changes default resource stage (dev) to current development stage
//	prod, sandbox, test, etc...
func setResourceStage(resource, stage string) string {
	if stage != DevelopmentStage {
		return strings.Replace(resource, DevelopmentStage, stage, -1)
	}

	return resource
}
