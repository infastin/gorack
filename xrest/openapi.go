package xrest

import (
	"errors"
	"fmt"
	"path"

	"gopkg.in/yaml.v3"
)

func SpecWithBasePath(spec []byte, basePath string) ([]byte, error) {
	var document map[string]any
	if err := yaml.Unmarshal(spec, &document); err != nil {
		return nil, err
	}

	aVersion, ok := document["openapi"]
	if !ok {
		return nil, errors.New(`invalid OpenAPI spec: "openapi" not found`)
	}

	version, ok := aVersion.(string)
	if !ok {
		return nil, errors.New(`invalid OpenAPI spec: "version" is not a string`)
	}

	switch version {
	case "2.0.0":
		document["basePath"] = path.Clean(basePath)
	case "3.0.0", "3.0.1", "3.0.2", "3.0.3", "3.1.0":
		var servers []map[string]any
		if aServers, ok := document["servers"]; ok {
			servers, ok = aServers.([]map[string]any)
			if !ok {
				return nil, errors.New(`invalid OpenAPI spec: "servers" is not an array of objects`)
			}
		}
		if len(servers) == 0 {
			servers = append(servers, map[string]any{
				"url": path.Clean(basePath),
			})
		} else {
			for _, srv := range servers {
				srv["url"] = path.Clean(basePath)
			}
		}
		document["servers"] = servers
	default:
		return nil, fmt.Errorf("invalid OpenAPI spec: invalid version %q", version)
	}

	return yaml.Marshal(document)
}
