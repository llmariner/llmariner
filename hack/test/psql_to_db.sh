#!/usr/bin/env bash

set -euo pipefail

kubectl exec -it -n postgres deploy/postgres -- psql -h localhost -U ps_user --no-password -p 5432 -d job_manager
