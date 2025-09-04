#!/usr/bin/env bash

set -euo pipefail

# TODO(kenji): The default API key needs to be configured so that it is excluded from rate-limiting.
export LLMARINER_API_KEY=default-key-secret

go run main.go completions.go http.go
