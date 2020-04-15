# CloudSQL Op

## Requirements

### Ops Platform

Running this op requires you to have access to the [Ops Platform](https://cto.ai/platform). Please review the [documentation](https://cto.ai/docs/overview) for detailed instructions on how to install the Ops CLI and/or Ops Slack application.

### GCP Credentials

A GCP Service Account is required to use this Op. The following predefined roles are required for all of the features in this Op to function as expected:

* Cloud SQL Admin
* Cloud SQL Viewer

Please refer to [this URL](https://cloud.google.com/iam/docs/creating-managing-service-accounts) for instructions on how to create a service account with the above mentioned permissions. Once created, you will need to also generate a private key for the respective service account and download it to your computer (JSON). When ready, run the following command to save the credentials as a secret in your Ops team, replacing <key_file> with the full path to your credentials JSON file:

```sh
ops secrets:set -k GOOGLE_APPLICATION_CREDENTIALS -v "$(cat <key_file> | tr -d '\n')"
```

## Usage

### CLI

```sh
ops run cloudsql
```

### Slack

```
/ops run cloudsql
```

## Features

* List all existing CloudSQL instances
* Provision new public CloudSQL instances
* Delete existing CloudSQL instances
* Clone existing CloudSQL instances
