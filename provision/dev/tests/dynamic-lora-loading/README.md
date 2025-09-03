# Test for Dynamic LoRA loading

The test requires vLLM as a default runtime and use Inference Sim container image.

Here is an example command to provision LLMariner:

```bash
helmfile apply --skip-diff-on-install \
  --state-values-set llmariner.useInferenceSim=true \
  --state-values-set llmariner.enableDynamicLoRALoading=true \
  --state-values-set llmariner.defaultRuntime=vllm \
  --state-values-set llmariner.enableHuggingFaceDownload=true \
  --state-values-set llmariner.enableDriftedPodUpdater=true \
  --state-values-set llmariner.enableInferenceManagerGracefulShutdown=true
```
