from mlflow.deployments import get_deploy_client

client = get_deploy_client("http://localhost:5000")
endpoints = client.list_endpoints()

response = client.predict(
    endpoint=endpoints[0].name,
    inputs={"prompt": "Tell me a joke about rabbits"},
)
print(response)
