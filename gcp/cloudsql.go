package gcp

import (
	"anteater/cloudsql/logger"
	"anteater/cloudsql/osutil"
	"encoding/json"
	"fmt"
)

// CloudSQLClient - Instance of API client
type CloudSQLClient struct {
	Logger *logger.Logger
}

// Instance - CloudSQL instance
type Instance struct {
	Name      string      `json:"name"`
	Region    string      `json:"region"`
	Zone      string      `json:"gceZone"`
	DBVersion string      `json:"databaseVersion"`
	Settings  SQLSettings `json:"settings"`
	State     string      `json:"state"`
	IPs       []IPAddress `json:"ipAddresses"`
}

// SQLSettings - CloudSQL Settings Object
type SQLSettings struct {
	Tier       string `json:"tier"`
	DiskSizeGB string `json:"dataDiskSizeGb"`
	DiskType   string `json:"dataDiskType"`
}

// IPAddress - Object containing info about IP addresses
type IPAddress struct {
	IP   string `json:"ipAddress"`
	Type string `json:"type"`
}

// getInstances - returns all CloudSQL instances
func (c *CloudSQLClient) getInstances() ([]Instance, error) {
	var sqlInstances []Instance
	cmd := "gcloud sql instances list --format json"
	stdout, err := osutil.ExecCmdWithLogs(cmd)
	if err != nil {
		return sqlInstances, err
	}

	err = json.Unmarshal([]byte(stdout), &sqlInstances)
	if err != nil {
		return sqlInstances, err
	}

	return sqlInstances, nil
}

// createInstance - creates and returns a CloudSQL instance
func (c *CloudSQLClient) createInstance(dbname string, region string, dbVersion string, tier string) (Instance, error) {
	var sqlInstance Instance
	cmd := fmt.Sprintf("gcloud sql instances create %v --region %v --database-version %v --tier %v --format json", dbname, region, dbVersion, tier)
	stdout, err := osutil.ExecCmdWithLogs(cmd)
	if err != nil {
		return sqlInstance, err
	}

	err = json.Unmarshal([]byte(stdout), &sqlInstance)
	if err != nil {
		return sqlInstance, err
	}

	return sqlInstance, nil
}

// cloneInstance - creates and returns a CloudSQL instance
func (c *CloudSQLClient) cloneInstance(dbnameToClone string, dbname string) ([]Instance, error) {
	var sqlInstances []Instance
	cmd := fmt.Sprintf("gcloud sql instances clone %v %v --format json", dbnameToClone, dbname)
	stdout, err := osutil.ExecCmdWithLogs(cmd)
	if err != nil {
		return sqlInstances, err
	}

	err = json.Unmarshal([]byte(stdout), &sqlInstances)
	if err != nil {
		return sqlInstances, err
	}

	return sqlInstances, nil
}

// deleteInstance - deletes a CloudSQL instance
func (c *CloudSQLClient) deleteInstance(dbname string) (string, error) {
	cmd := fmt.Sprintf("gcloud sql instances delete %v --quiet --format json", dbname)
	stdout, err := osutil.ExecCmdWithLogs(cmd)
	if err != nil {
		return stdout, err
	}

	return stdout, nil
}

// List - outputs all CloudSQL instances
func (c *CloudSQLClient) List() (string, error) {
	sqlInstances, err := c.getInstances()
	if err != nil {
		return "", err
	}
	output := formatSQLInstances(sqlInstances)
	return output, nil
}

// GetListOfInstances - gets a list of all CloudSQL instances
func (c *CloudSQLClient) GetListOfInstances() ([]string, error) {
	sqlInstances, err := c.getInstances()
	if err != nil {
		return []string{}, err
	}
	var listOfInstances []string
	for _, obj := range sqlInstances {
		listOfInstances = append(listOfInstances, obj.Name)
	}
	return listOfInstances, nil
}

// Provision - creates a CloudSQL instance
func (c *CloudSQLClient) Provision(dbname string, region string, dbVersion string, tier string) (string, error) {
	c.Logger.LogInfo(fmt.Sprintf("⚙️  Please wait. Creating database `%v (%v | %v)` in `%v` with default settings...", dbname, dbVersion, tier, region))
	sqlInstance, err := c.createInstance(dbname, region, dbVersion, tier)
	if err != nil {
		return "", err
	}
	output := formatSQLInstances([]Instance{sqlInstance})
	c.Logger.LogInfo(fmt.Sprintf("✅  Database `%v (%v | %v)` successfully created in `%v`!", dbname, dbVersion, tier, region))
	return output, nil
}

// Clone - clones a CloudSQL instance
func (c *CloudSQLClient) Clone(dbnameToClone string, dbname string) (string, error) {
	c.Logger.LogInfo(fmt.Sprintf("⚙️  Please wait. Creating database `%v` based on `%v`...", dbname, dbnameToClone))
	sqlInstances, err := c.cloneInstance(dbnameToClone, dbname)
	if err != nil {
		return "", err
	}
	output := formatSQLInstances(sqlInstances)
	c.Logger.LogInfo(fmt.Sprintf("✅  Database `%v` (clone of `%v`) successfully created!", dbname, dbnameToClone))
	return output, nil
}

// Delete - deletes a CloudSQL instance
func (c *CloudSQLClient) Delete(dbname string) (string, error) {
	c.Logger.LogInfo(fmt.Sprintf("⚙️  Please wait. Deleting database `%v`...", dbname))
	output, err := c.deleteInstance(dbname)
	if err != nil {
		return "", err
	}
	c.Logger.LogInfo(fmt.Sprintf("✅  Database `%v` successfully deleted!", dbname))
	return output, nil
}

func formatSQLInstances(sqlInstances []Instance) string {
	var formattedData string
	for _, obj := range sqlInstances {
		formattedData += fmt.Sprintf(
			"Name: `%s`\nRegion: `%s`\nZone: `%s`\nTier: `%s`\nDisk Size: `%s GB`\nDisk Type: `%s`\nDatabase Version:`%s`\nIP Addresses: `%+q`\nState: `%s`\n\n",
			obj.Name,
			obj.Region,
			obj.Zone,
			obj.Settings.Tier,
			obj.Settings.DiskSizeGB,
			obj.Settings.DiskType,
			obj.DBVersion,
			obj.IPs,
			obj.State)
	}
	return formattedData
}
