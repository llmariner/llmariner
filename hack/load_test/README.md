# Load test

This is a simple load test for chat completion requests.

You can check the requests rate with the following Prometheus query:

```
sum by (status_code) (rate(llmariner_inference_manager_server_request_count[1m]))
```
