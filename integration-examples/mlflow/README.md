# MLflow Experiment

## Install

Run `./deploy.sh` to deploy MLflow.

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

## MLflow Deployments Server for LLMs (Experimental)

Take the following steps to run a deployment server in a K8s cluster:

```bash
kubectl create secret generic -n mlflow llmariner-api-key \
  --from-literal=secret=${OPENAI_API_KEY}
kubectl apply -n mlflow -f deployment-server.yaml
```

## Prompt Engineering UI (Experimental)

Follow https://mlflow.org/docs/latest/llms/prompt-engineering/index.html.

You can access the MLflow Tracking Server, click "New run" and choose "using Prompt Engineering".

## Run an MLflow Project on Kubernetes

See https://www.mlflow.org/docs/latest/projects.html#kubernetes-execution

## Authentication & Authorization

By default, MLflow has its own user management. If we change that to OIDC and use the same authorizaiton mechanism as LLMOperator,
we need to either add an authorization plugin to MLflow or put a reverse proxy in front of MLflow. The latter approach might work
better if we need to put a similar authorization to other services such as Grafana.

See https://github.com/data-platform-hq/mlflow-oidc-auth and https://www.mlflow.org/docs/latest/auth/index.html#configuration.
