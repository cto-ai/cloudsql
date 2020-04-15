package main

import (
	"anteater/cloudsql/gcp"
	"anteater/cloudsql/logger"
	"encoding/json"
	"fmt"
	"os"

	ctoai "github.com/cto-ai/sdk-go"
)

// GCPConfig GOOGLE_APPLICATION_CREDENTIALS secrets strucutre
type GCPConfig struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

const (
	GoogleApplicationCredentials = "GOOGLE_APPLICATION_CREDENTIALS"
	List                         = "List"
	Provision                    = "Provision"
	Delete                       = "Delete"
	Clone                        = "Clone"
	Exit                         = "Exit"
)

func main() {
	var ux = ctoai.NewUx()
	var sdk = ctoai.NewSdk()
	var prompt = ctoai.NewPrompt()

	logger := &logger.Logger{
		Sdk: sdk,
	}

	gcpAuthClient := &gcp.AuthClient{
		Logger: logger,
	}

	gcpCloudSQLClient := &gcp.CloudSQLClient{
		Logger: logger,
	}

	gcpCreds, err := sdk.GetSecret(GoogleApplicationCredentials)
	if err != nil {
		logger.LogError("Failed to retrieve secrets", err)
		ux.Print("😅  Failed to retrieve secret containing GCP credentials")
		os.Exit(1)
	}

	var gcpConfig GCPConfig
	err = json.Unmarshal([]byte(gcpCreds), &gcpConfig)
	if err != nil {
		logger.LogError("Failed to unmarshal GCP credentials", err)
		ux.Print("😅  Invalid GCP credentials")
		os.Exit(1)
	}
	err = gcpAuthClient.CreateKeyFile(gcpCreds)
	if err != nil {
		logger.LogError("Failed to write keyfile", err)
		ux.Print("😅  Failed to set up GCP credentials")
		os.Exit(1)
	}

	err = gcpAuthClient.Authenticate(gcpConfig.ProjectID)
	if err != nil {
		logger.LogError("Failed to Authenticate GCP", err)
		ux.Print("😅  Failed to authenticate with Google Cloud")
		os.Exit(1)
	}
	ux.Print("🔓  User Authenticated\n")

	actionChoices := []string{List, Provision, Clone, Delete, Exit}

	var runAgain bool = true
	for {
		runAgain = promptForAction(ux, prompt, gcpCloudSQLClient, logger, actionChoices)
		if !runAgain {
			break
		}
	}
	os.Exit(0)
}

func promptForAction(ux *ctoai.Ux, prompt *ctoai.Prompt, gcpCloudSQLClient *gcp.CloudSQLClient, logger *logger.Logger, actionChoices []string) bool {
	resp, err := prompt.List("actions", "Please select an action to perform", actionChoices, ctoai.OptListAutocomplete(true))
	if err != nil {
		logger.LogError("Failed to prompt action list", err)
		ux.Print("😅  Failed to retrieve and display list of actions")
		return false
	}

	switch resp {
	case List:
		output, err := gcpCloudSQLClient.List()
		if err != nil {
			logger.LogError("Failed to list CloudSQL instances", err)
			ux.Print(fmt.Sprintf("😅  Failed to list CloudSQL instance:\n%+v", err))
			return false
		}
		ux.Print(output)
	case Provision:
		dbname, err := prompt.Input("dbname", "Please enter a name for your database")
		if err != nil {
			logger.LogError("Failed to prompt for database name", err)
			ux.Print("😅  Failed to retrieve input")
			return false
		}
		region, err := prompt.List("region", "Please select a region", gcp.GetRegions(), ctoai.OptListAutocomplete(true))
		if err != nil {
			logger.LogError("Failed to prompt for region", err)
			ux.Print("😅  Failed to retrieve input")
			return false
		}
		dbVersion, err := prompt.List("region", "Please select a database version", gcp.GetDBVersions(), ctoai.OptListAutocomplete(true))
		if err != nil {
			logger.LogError("Failed to prompt for database version", err)
			ux.Print("😅  Failed to retrieve input")
			return false
		}
		tier, err := prompt.List("region", "Please select a database tier", gcp.GetTiers(dbVersion), ctoai.OptListAutocomplete(true))
		if err != nil {
			logger.LogError("Failed to prompt for database tier", err)
			ux.Print("😅  Failed to retrieve input")
			return false
		}
		output, err := gcpCloudSQLClient.Provision(dbname, region, dbVersion, tier)
		if err != nil {
			logger.LogError("Failed to provision CloudSQL instance", err)
			ux.Print(fmt.Sprintf("😅  Failed to provision CloudSQL instance:\n%+v", err))
			return false
		}
		ux.Print(output)
		return true
	case Delete:
		existingInstances, err := gcpCloudSQLClient.GetListOfInstances()
		if err != nil {
			logger.LogError("Failed to retrieve list of CloudSQL instances", err)
			ux.Print(fmt.Sprintf("😅  Failed to retrieve list of CloudSQL instances:\n%+v", err))
			return false
		}
		if len(existingInstances) == 0 {
			ux.Print("😅  No existing CloudSQL instances in this project")
			return true
		}
		dbname, err := prompt.List("dbname", "Please enter the name of the database you would like to delete", existingInstances, ctoai.OptListAutocomplete(true))
		if err != nil {
			logger.LogError("Failed to prompt for database name", err)
			ux.Print("😅  Failed to retrieve input")
			return false
		}
		confirm, err := prompt.Input("confirm", "Please confirm deletion by typing `yes`")
		if err != nil {
			logger.LogError("Failed to prompt for database deletion confirmation", err)
			ux.Print("😅  Failed to retrieve input")
			return false
		}
		if confirm != "yes" {
			ux.Print("❌  Deletion not confirmed, returning to main menu")
			return true
		}
		output, err := gcpCloudSQLClient.Delete(dbname)
		if err != nil {
			logger.LogError("Failed to delete CloudSQL instance", err)
			ux.Print(fmt.Sprintf("😅  Failed to delete CloudSQL instance:\n%+v", err))
			return false
		}
		ux.Print(output)
		return true
	case Clone:
		existingInstances, err := gcpCloudSQLClient.GetListOfInstances()
		if err != nil {
			logger.LogError("Failed to retrieve list of CloudSQL instances", err)
			ux.Print(fmt.Sprintf("😅  Failed to retrieve list of CloudSQL instances:\n%+v", err))
			return false
		}
		if len(existingInstances) == 0 {
			ux.Print("😅  No existing CloudSQL instances in this project")
			return true
		}
		dbnameToClone, err := prompt.List("dbnameToClone", "Please enter the name of the database you would like to clone", existingInstances, ctoai.OptListAutocomplete(true))
		if err != nil {
			logger.LogError("Failed to prompt for database name", err)
			ux.Print("😅  Failed to retrieve input")
			return false
		}
		dbname, err := prompt.Input("dbname", "Please enter a name for the new database you would like to create")
		if err != nil {
			logger.LogError("Failed to prompt for database name", err)
			ux.Print("😅  Failed to retrieve input")
			return false
		}
		output, err := gcpCloudSQLClient.Clone(dbnameToClone, dbname)
		if err != nil {
			logger.LogError("Failed to clone CloudSQL instance", err)
			ux.Print(fmt.Sprintf("😅  Failed to clone CloudSQL instance:\n%+v", err))
			return false
		}
		ux.Print(output)
		return true
	case Exit:
		return false
	}
	return true
}
