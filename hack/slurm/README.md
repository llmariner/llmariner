# Slurm deployment with Slink

[Slink](https://slinky.ai) is a project from SchedMD, and one of its open source tools
is [slurm-operator](https://github.com/SlinkyProject/slurm-operator).

Please [this presentation](https://slurm.schedmd.com/SLUG24/Slinky-Slurm-Operator.pdf) for its overview and future plan.

To deploy Slurm Operator and create a Slurm cluster,
follow [the quick start guide](https://github.com/SlinkyProject/slurm-operator/blob/main/docs/user/quickstart.md)

The above steps will create the following `Cluster` resource and `NodeSet` resource.

```yaml
kind: Cluster
metadata:
  name: slurm
  namespace: slurm
spec:
  server: http://slurm-restapi.slurm:6820
  token:
    secretRef: slurm-token-slurm
```

```yaml
apiVersion: slinky.slurm.net/v1alpha1
kind: NodeSet
metadata:
  name: slurm-compute-debug
  namespace: slurm
spec:
  clusterName: slurm
  replicas: 1
  serviceName: slurm-compute
  minReadySeconds: 0
  updateStrategy:
    rollingUpdate:
      maxUnavailable: 20%
      partition: 0
      paused: false
    type: RollingUpdate
  persistentVolumeClaimRetentionPolicy:
    whenDeleted: Retain
  selector:
    matchLabels:
      app.kubernetes.io/instance: slurm-compute-debug
      app.kubernetes.io/name: slurmd
  template:
    ...
    spec:
      hostname: slurm-compute-debug
      priorityClassName:
      automountServiceAccountToken: false
      dnsConfig:
        searches:
          - slurm-controller.slurm.svc.cluster.local
          - slurm-compute.slurm.svc.cluster.local
      nodeSelector:
        kubernetes.io/os: linux
```

Then Slurm Operator will then create the following pods in the `slurm` namespace.

```console
$ kubectl get pods -n slurm
NAME                              READY   STATUS    RESTARTS      AGE
slurm-accounting-0                1/1     Running   1 (61m ago)   61m
slurm-compute-debug-7xmzv         1/1     Running   0             59m
slurm-controller-0                2/2     Running   1 (60m ago)   61m
slurm-exporter-7b44b6d856-k8vz4   1/1     Running   0             61m
slurm-mariadb-0                   1/1     Running   0             61m
slurm-restapi-5f75db85d9-l9jpj    1/1     Running   0             61m
```


To test the Slurm CLIs:

```bash
kubectl --namespace=slurm exec -it statefulsets/slurm-controller -- bash --login
sinfo
srun hostname
sbatch --wrap="sleep 600"
squeue
```

To test the REST endpoint:

```console
kubectl port-forward -n slurm service/slurm-restapi 6820 &
TOKEN=$(kubectl get secrets -n slurm slurm-token-slurm -o jsonpath='{.data.auth-token}' | base64 -d)

cat << EOF > submit_job.json
{
    "job": {
        "name": "test",
        "tasks": 1,
        "nodes": "1",
        "environment": [
            "PATH=/bin:/usr/bin/:/usr/local/bin/",
            "LD_LIBRARY_PATH=/lib/:/lib64/:/usr/local/lib"
        ],
		"current_working_directory": "/tmp",
        "standard_input": "/dev/null",
        "standard_output": "/tmp/test.out",
        "standard_error": "/tmp/test_error.out",
    },
    "script": "#!/bin/bash\necho HELLO"
  }
EOF

curl \
  -X POST \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-type: application/json" \
  --data-binary @submit_job.json \
  http://localhost:6820/slurm/v0.0.41/job/submit

curl \
  -H "Authorization: Bearer ${TOKEN}" \
  http://localhost:6820/slurm/v0.0.41/jobs/
```

You can check the output of the job by accessing the `slurm-compute-debug` pod:

```bash
kubectl exec -it -n slurm <slurm-compute-debug pod> cat /tmp/test.out
```

Please see https://slurm.schedmd.com/SLUG23/REST-API-SLUG23.pdf and
https://aws.amazon.com/blogs/hpc/using-the-slurm-rest-api-to-integrate-with-distributed-architectures-on-aws/ for example API calls.

Please see https://slurm.schedmd.com/rest_api.html for the API spec.

Please note that the API returns 500 status code (not 401 or 403) even if an invalid auth token is given.
