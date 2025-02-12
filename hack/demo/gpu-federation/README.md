# GPU Federation

This directory contains the instruction and scripts for running a demo for GPU federation.

## A Demo on a Single machine

In this demo, we set up the following three Kind clusters in a single machine (e.g., single EC2 instance).

- `tenant-cluster`
- `gpu-worker-cluster-large`
- `gpu-worker-cluster-small`

LLMariner worker-plane components will connect to the LLMariner control-plane hosted at https://api.llm.staging.cloudnatix.com/v1.


The username of the demo account is `demo+gpu@cloudnatix.com`. The password is stored in the company's 1 Password.

### Step 1. Set up Kind clusters and deploy LLMariner

Run:

```bash
`./set_up_kind_based_env.sh
```

### Step 2. Submit GPU jobs

TODO(kenji): Fill. You can access the frontend or run `llma hidden jobs clusters list` to check the status.


### Step 3. Tear down the demo env

Run:

```bash
`./tear_down_kind_based_env.sh
```
