# LLM Operatror

This repository contains a Helm chart, tutorial, and AWS provisioning script.
Please visit [our documentation site](https://llm-operator.readthedocs.io/) for more details.


## Development Notes

### How to Update the Helm chart

Run the following command to bump the versions of sub-charts:

```console
python3 scripts/update_chart.py deployments/llm-operator/Chart.yaml
```
