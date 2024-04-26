# LLM Operatror

LLM Operator converts your GPU clusters to a platform for generative AI workloads.

# Key Values

- *Provide LLM as a service.* LLM Operator builds a software stack that provides LLM as a service, including inference, fine-tuning, model management, and training data management.
- *Utilize GPU optimally.* LLM Operator provides auto-scaling of inference-workloads, efficient scheduling of fine-tuning batch jobs, GPU sharing, etc.

# Use Cases

- Develop LLM applications with the API that is compatible with [OpenAI-compatible API](https://platform.openai.com/docs/api-reference).
- Fine-tune models while keeping data safely and securely in your on-premise datacenter.
- Run fine-tuning jobs efficiently with guaranteed SLO and without interference with inference requests.

# High-level Architecture

![Architecture Diagram](docs/images/architecture_diagram.png)

# An Initial Demo Scenario

1. A user uploads a dataset to File Manager.
2. The user creates a fine-tuning job in Job Manager. Job Manager generates a LoRA adapter with the uploaded dataset and stores the LoRA adapter in Model Registry.
3. Inference Manager is notified and imports a new model.
4. The user runs a chatbot using the fine-tuned model.

Please see [the demo video](https://drive.google.com/file/d/1IIDytriu4Cl1O9Wo7fXzHkS1kbqJxfXO/view?usp=sharing).


# Technical Challenges

- Be able to satisfy both the SLO of fine tuning jobs and inference on a limited number of GPUs (e.g., Run a large fine-tuning jobs at midnight when no one is using inference)
- Support heterogeneous GPUs (from A100 to B100)
- Support heterogeneous models (from small models to large models)
