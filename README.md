# LLMariner (= LLM + Mariner)

LLMariner transforms your GPU clusters into a powerhouse for generative AI workloads.

![alt text](https://github.com/llmariner/.github/blob/main/images/logo.png?raw=true)

This repository contains a Helm chart, cli, tutorials, and a provisioning script for a playground.

Please visit [our documentation site](https://llmariner.readthedocs.io/) to learn LLMariner.

## Development Notes

### How to Update the Helm chart

Run the following command to bump the versions of sub-charts:

```console
python3 scripts/update_chart.py deployments/llmariner/Chart.yaml
```
