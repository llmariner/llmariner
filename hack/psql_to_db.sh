#!/usr/bin/env bash

set -euo pipefail

kubectl exec -it -n postgres postgres-0 -- env PGPASSWORD=ps_password psql -h localhost -U ps_user -p 5432 -d job_manager
