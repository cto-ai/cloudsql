package gcp

import "strings"

// GetRegions - returns list of GCP Compute regions
func GetRegions() []string {
	return []string{"asia-east1", "asia-east2", "asia-northeast1", "asia-northeast2", "asia-northeast3", "asia-south1", "asia-southeast1", "australia-southeast1", "europe-north1", "europe-west1", "europe-west2", "europe-west3", "europe-west4", "europe-west6", "northamerica-northeast1", "southamerica-east1", "us-central1", "us-east1", "us-east4", "us-west1", "us-west2", "us-west3"}
}

// GetDBVersions - returns list of supported database versions
func GetDBVersions() []string {
	return []string{"POSTGRES_11", "POSTGRES_9_6", "MYSQL_5_7", "MYSQL_5_6", "MYSQL_5_5"}
}

// GetTiers - returns list of supported db tiers
func GetTiers(dbVersion string) []string {
	if strings.Contains(dbVersion, "POSTGRES") {
		return []string{"db-f1-micro", "db-g1-small"}
	}
	return []string{"db-f1-micro", "db-g1-small", "db-n1-standard-1", "db-n1-standard-2", "db-n1-standard-4", "db-n1-standard-8", "db-n1-standard-16", "db-n1-standard-32", "db-n1-standard-64", "db-n1-standard-96", "db-n1-highmem-2", "db-n1-highmem-4", "db-n1-highmem-8", "db-n1-highmem-16", "db-n1-highmem-32", "db-n1-highmem-64", "db-n1-highmem-96"}
}
