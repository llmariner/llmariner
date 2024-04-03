# LLM Operatror

LLM Operator builds a LLM stack that provides the following functionality:

- LLM fine-tuning job management
- LLM inference (compatible with OpenAI API)
- (LoRA) / fine-tuning Model repository

Additionally it provides the following components as optional:
- Vector DB (e.g., https://milvus.io/)
- Dataset Storage
- GPU Operator
- Monitoring
- MLFlow

Here are some of the challenges:

- Be able to satisfy both the SLO of fine tuning jobs and inference on a limited number of GPUs (e.g., Run a large fine-tuning jobs at midnight when no one is using inference)
- Support heterogeneous GPUs (from A100 to B100)
- Support heterogeneous models (from small models to large models)

# High-level Architecture

![Architecture Diagram](docs/images/architecture_diagram.png)
