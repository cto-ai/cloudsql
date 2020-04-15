package gcp

import (
	"anteater/cloudsql/logger"
	"anteater/cloudsql/osutil"
	"fmt"
	"io/ioutil"
)

// AuthClient - client for GCP
type AuthClient struct {
	Logger *logger.Logger
}

const googleApplicationCredentialsFilePath = "/ops/gcp.json"

// Authenticate gcp auth function CreateKeyFile must be called before this
func (g *AuthClient) Authenticate(projectID string) error {
	command := fmt.Sprintf("gcloud auth activate-service-account --quiet --key-file %s --project %s", googleApplicationCredentialsFilePath, projectID)

	err := osutil.ExecCmd(command)
	if err != nil {
		return err
	}
	return nil
}

// CreateKeyFile - writes GCP creds file
func (g *AuthClient) CreateKeyFile(credentials string) error {
	data := []byte(credentials)
	err := ioutil.WriteFile(googleApplicationCredentialsFilePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
