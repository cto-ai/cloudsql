![CTO Banner](assets/banner.png)

# CloudSQL Op

An Op that facilitates the management of [CloudSQL](https://cloud.google.com/sql) instances in GCP (Google Cloud Platform).

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
* Provision new **public** CloudSQL instances
* Delete existing CloudSQL instances
* Clone existing CloudSQL instances

### CloudSQL Op Walkthrough - Watch video

[![CloudSQL Op Walkthrough - Watch video](https://img.youtube.com/vi/iM1LHHRqGQU/hqdefault.jpg)](https://youtu.be/iM1LHHRqGQU)

## Local Development / Running from Source

**1. Clone the repo:**

```bash
git clone <git url>
```

**2. Navigate into the directory and install dependencies:**

```bash
cd cloudsql
```

**3. Run the Op from your current working directory with:**

```bash
ops run . --build
```

## Contributing

See the [Contributing Docs](CONTRIBUTING.md) for more information

## Contributors

<table>
  <tr>
    <td align="center"><a href="https://github.com/minsohng"><img src="https://avatars2.githubusercontent.com/u/19717602?s=100" width="100px;" alt=""/><br /><sub><b>Min Sohng</b></sub></a><br/></td>
    <td align="center"><a href="https://github.com/ruxandrafed"><img src="https://avatars2.githubusercontent.com/u/11021586?s=100" width="100px;" alt=""/><br /><sub><b>Ruxandra Fediuc</b></sub></a><br/></td>
  </tr>
</table>

## License

[MIT](LICENSE)