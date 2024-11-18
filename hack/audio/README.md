# Chat completion request with audio

## Sample configuration to support audio chat completion with vLLM.

```console
inference-manager-engine:
  runtime:
    runtimeImages:
      # The modified container image of vLLM.
      vllm: llmariner/vllm-openai:0.6.2
  model:
    default:
      runtimeName: vllm
      resources:
        limits:
          cpu: 0
          memory: 0
          # Do not allocate GPU to inference-manager-engine since g5.4xlarge has only one GPU,
          # and it is needed for the fine-tuning job
          nvidia.com/gpu: 1
        requests:
          cpu: 0
          memory: 0
    overrides:
      fixie-ai/ultravox-v0_3:
        preloaded: true
        contextLength: 4096


model-manager-loader:
  models:
  - model: fixie-ai/ultravox-v0_3
    baseModel: meta-llama/Meta-Llama-3.1-8B-instruct
  huggingFaceSecret:
    name: hf-token-secret
    apiKeyKey: apiKeyKey
  downloader:
    kind: huggingFace
    s3:
      endpointUrl: http://minio.minio:9000
      region: dummy
      bucket: llmariner
      pathPrefix: models
    huggingFace:
      cacheDir: /tmp/hf

```
