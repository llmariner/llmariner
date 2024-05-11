# MLflow Experiment

## Install

Run `./deploy.sh` to deploy Mlflow.

Run the following commands to get credentials.

```bash
export MLFLOW_TRACKING_USERNAME=$(kubectl get secret --namespace mlflow mlflow-tracking -o jsonpath="{ .data.admin-user }" | base64 -d)
export MLFLOW_TRACKING_PASSWORD=$(kubectl get secret --namespace mlflow mlflow-tracking -o jsonpath="{.data.admin-password }" | base64 -d)
echo ${MLFLOW_TRACKING_USERNAME}
echo ${MLFLOW_TRACKING_PASSWORD}
```

Run the following command to access the tracking server:

```bash
kubectl port-forward -n mlflow service/mlflow-tracking 9000:80
```

You can access `http://localhost:9000` and login MLflow with the above username and password.

## MLflow LLM Evaluate

MLflow provides ["MLflow LLM Evaluate"](https://mlflow.org/docs/latest/llms/llm-evaluate/index.html).

The following example runs the sample script:

```bash
export OPENAI_API_BASE='OpenAPI endpoint URL (e.g., http://localhost:8080/v1)'
export OPENAI_BASE_URL='OpenAPI endpoint URL (e.g., http://localhost:8080/v1)'
export OPENAI_API_KEY='your-api-key-here'

python eval.py
```

Here is the output:

```
See evaluation table below:
            inputs  ... token_count
0  What is MLflow?  ...          42
1   What is Spark?  ...          26

[2 rows x 4 columns]
```

You can also access http://localhost:9000 to see the results.

> [!NOTE]
> This is currently not fully tested. The scoring information might not be available.

## MLflow Deployments Server (Experimental)

Run:

```bash
cat << EOF | envsubst > config.yaml
endpoints:
- name: completions
  endpoint_type: llm/v1/completions
  model:
    provider: openai
    name: google-gemma-2b-it-q4
    config:
      openai_api_base: $OPENAI_API_BASE
      openai_api_key: $OPENAI_API_KEY
EOF

mlflow deployments start-server --config-path config.yaml
```

Then access `http://localhost:5000` or run `python test_endpoint.py`.
