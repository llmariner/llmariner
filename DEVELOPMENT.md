# Development Guide

## Versioning

The release version follows [Semantic Versioning](https://semver.org/).

## Style Guide

### Commit Message

The commit message follows [Conventional Commits specification](https://conventionalcommits.org/).

Each commit message consists of a header and an optional message. The header has a special format that includes a type, optional scope, optional breaking-change flag, and a subject:

```
<type>[(<scope>)][!]: <subject>

[<message>]
```

> [!NOTE]
> The bot will automatically adds labels to the pull-request based on the commit message header. This rule is defined in the [labeler.yml](.github/labeler.yml).

#### Header Type

* **feat**: A new feature
* **fix**: A bug fix
* **docs**: Documentation only changes
* **ci**: Changes to our CI configuration files and scripts
* **chore**: Changes that do not affect the meaning of the code (formatting, fixing a typo, etc)
* **refactor**: A code change that neither fixes a bug nor adds a feature

### Helm Chart

The deployment manifest is packaged in a [Helm](https://helm.sh/) chart. The chart is uploaded to the OCI repository and registered in [ArtifactHub](https://artifacthub.io/packages/helm/llmariner/llmariner).

To ensure the values in values.yaml are correct and to make it easier for users to understand the configuration settings, we generate a values schema using the [helm-tool](https://github.com/cert-manager/helm-tool).

You can generate a values.schema.json and verify the chart by running:

```console
# just generate schema
make generate-chart-schema

# generate schema and verify chart
make helm-lint
```

#### Comments for chart values

For improved usability, write a comment for each value. Each comment includes a description, optional reference, optional example, and optional `helm-tool` tags. For `helm-tool` tags, refer to [Tags](https://github.com/cert-manager/helm-tool?tab=readme-ov-file#tags).

Sample:

```yaml
# The application name.
name: llmariner

# The port number for serving sample APIs.
# +docs:type=number
port: 1234

# Specify the key=value settings.
# +docs:property
sampleMapValues:
  keyA: value-1
  keyB: value-2

# Optional sample values to describe the comments for chart values.
# For more information, see [Sample Document](http://example.com).
#
# For example:
# tls:
#   foo: bar
#   names:
#   - alice
#
# +docs:property
# sampleOptionalValue:
#   foo: ""
#   names: []
```

> [!TIP]
> - The type is generally inferred based on the default value; however, for numbers, you should explicitly specify `+docs:type=number`.
> - Add the `+docs:property` tag to the commented-out value for detecting the setting value.
> - When setting a non-empty value to a map-type setting, ensure to set the `+docs:property` tag to prevent the value type from being locked.

### Code

LLMariner, like most Go projects, delegates almost all stylistic choices to `gofmt`.
We also use some linters. Please verify that the code meets our guidelines by running:

```console
make lint
```

> [!NOTE]
> Besides the linters mentioned above, `golangci-lint` is also executed for [pull requests](.github/workflows/ci-pre-merge.yml).

## Testing

### Unit Test

We are using the standard Go test command. The following is the entrypoint for running all unit tests:

```console
make test
```

### Integration Test

There are two options for running tests: locally or using [GitHub Actions](https://github.com/llmariner/llmariner/actions/workflows/manual-integration-test.yaml). Please note that the GitHub Actions option is only available to users who have write permission for this repository.

For testing locally, please refer to [Provision LLMariner for development](provision/dev/README.md).
