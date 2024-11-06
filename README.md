<p align="center">
  <img title="LLMariner" alt="LLMariner" width="55%" src="img/logo.png">
</p>
<p align="center">
  <a href="https://llmariner.ai/"><b>Documentation</b></a> |
  <a href="https://llmariner.slack.com/join/shared_invite/zt-2rbwooslc-LIrUCmK9kklfKsMEirUZbg#/shared-invite/email"><b>Community Slack</b></a>
</p>

---

LLMariner (= LLM + Mariner) is an extensible open source platform to simplify the management of generative AI workloads. Built on Kubernetes, it enables you to efficiently handle both training and inference data within your own clusters. With [OpenAI-compatible APIs](https://platform.openai.com/docs/api-reference), LLMariner leverages an ecosystem of tools, facilitating seamless integration for a wide range of AI-driven applications.

<p align="center">
  <img src="https://llmariner.ai/images/concepts.png" width=80% title="LLMariner concepts" alt="LLMariner concepts">
</p>

## Architecture

LLMariner consists of a control-plane and one or more worker-planes. Both components can operate within a single cluster, but if you want to utilize GPU resources across multiple clusters, they can also be installed into separate clusters:

<dl>
  <dt>Control-Plane components:</dt>
  <dd>Expose the OpenAI-compatible APIs and manage the overall state of LLMariner and receive a request from the client.</dd>
  <dt>Worker-Plane components:</dt>
  <dd>Run every worker cluster, process tasks using compute resources such as GPUs in response to requests from the control-plane.</dd>
</dl>

<p align="center">
  <img src="https://llmariner.ai/images/highlevel_architecture.png" width=75% title="LLMariner High-level Architecture" alt="LLMariner High-level Architecture">
</p>

Refer to the [High-Level Architecture](https://llmariner.ai/docs/overview/how-works/) document and [Technical Details](https://llmariner.ai/docs/dev/architecture/) document for more information.

## Installation

Check out our [installation guide](https://llmariner.ai/docs/setup/install/), which covers several deployment methods, including setting up a test environment using the kind cluster and building a production-ready environment, among others.

## Integration

LLMariner provides OpenAI-compatible APIs, making it easy to integrate with powerful tools such as assistant web UIs, code generation tools, and more. Here are some integration samples:

- **Open WebUI**: A self-hosted web UI that works with OpenAI-compatible APIs. See [Open WebUI](https://llmariner.ai/docs/integration/openwebui/) integration guide for details.
- **Continue**: An open-source AI code assistant inside of VS Code and JetBrains. See [Continue](https://llmariner.ai/docs/integration/continue/) integration guide for details.
- **Weights & Biases (W&B)**: AI developer platform that can enable you to easily see the progress of your fine-tuning jobs, such as training epoch, loss, etc. See [W&B](https://llmariner.ai/docs/integration/wandb/) integration guide for details.

## Directory Structure

- `cli`: CLI
- `deployments`: Helm chart
- `integration-examples`: Examples of integration of other services with LLMariner
- `provision`: provisioning scripts
- `tutorials`: Tutorials

## Talks

- [Transform Your Kubernetes Cluster Into a GenAI Platform: Get Ready-to-Use LLM APIs Today! - Cloud Native & Kubernetes AI Day 2024 North America](https://sched.co/1izue)

## Contributing

See [Contributing Guide](CONTRIBUTING.md).

## License

See [LICENSE](LICENSE).
