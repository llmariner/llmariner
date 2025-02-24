# Scripts for deploying the latest version of LLMariner components.

The following is a procedure for deploying the latest container images and Helm files from
the individual repos.

```bash
./set_up_repos.sh
cp ./Chart.yaml ../../deployments/llmariner/Chart.yaml
helmfile apply --skip-diff-on-install --state-values-set llmariner.deployLatest=true
```
