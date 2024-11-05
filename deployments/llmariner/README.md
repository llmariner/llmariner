# LLMariner

[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/llmariner)](https://artifacthub.io/packages/search?repo=llmariner)

Please refer to the [full installation guide](https://llmariner.ai/docs/setup/install/). The LLMariner chart has some pre-requirements, such as setting up a relational database. The installation guide covers several deployment methods, including setting up a test environment using the kind cluster and building a production-ready environment.

## Configuration

See [Customizing the Chart Before Installing](https://helm.sh/docs/intro/using_helm/#customizing-the-chart-before-installing). To see all configurable options with detailed comments, visit the chart's [values.yaml](./values.yaml), or run these configuration commands:

```console
helm show values oci://public.ecr.aws/cloudnatix/llmariner-charts/llmariner
```

## Install Chart

```console
helm install <RELEASE_NAME> oci://public.ecr.aws/cloudnatix/llmariner-charts/llmariner
```

See [configuration](#configuration) below.
See [helm install](https://helm.sh/docs/helm/helm_install/) for command documentation.

## Uninstall Chart

```console
helm uninstall <RELEASE_NAME>
```

This removes all the Kubernetes components associated with the chart and deletes the release.
See [helm uninstall](https://helm.sh/docs/helm/helm_uninstall/) for command documentation.

## Upgrading Chart

```console
helm upgrade <RELEASE_NAME> oci://public.ecr.aws/cloudnatix/llmariner-charts/llmariner
```

See [helm upgrade](https://helm.sh/docs/helm/helm_upgrade/) for command documentation.
