# Fault Injection Test

This test injects faults by deleting pods while continuously sending
inference requests.

The `llmariner.enableInferenceManagerGracefulShutdown` state value
needs to be set to `true` for this test.

```bash
helmfile apply --skip-diff-on-install \
  --state-values-set llmariner.enableInferenceManagerGracefulShutdown=true
```
