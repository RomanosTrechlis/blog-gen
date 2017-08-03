package endpoint

import (
	"fmt"
)

type Endpoint interface {
	Upload(destFolder, endpointUsername, endpointPassword, endpointURL string) error
}

func New(endpointType string) (endpoint Endpoint, err error) {
	switch endpointType {
	case "git":
		endpoint = newGitEndpoint()
	default:
		err = fmt.Errorf("no endpoint information found in the config file")
	}
	return endpoint, err
}
