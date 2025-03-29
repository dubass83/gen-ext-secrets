#!/bin/bash

# This script is used to decode all secrets in the default namespace
# and print them in a human-readable format.
# It uses kubectl to get the secrets and jq to decode them.
# It is assumed that kubectl and jq are installed and configured to access the Kubernetes cluster.
# The script will create a JSON file for each secret in the default namespace.
# The JSON file will contain the decoded data of the secret.

set -xeo pipefail

SECRETS=$(kubectl get secrets -n default | grep -v helm | grep -v kubernetes.io/tls | tail -n +2 | awk '{print $1}' | tr '\n' ' ')

for secret in $SECRETS; do
  echo "Secret: $secret"
  echo "writing to $secret.json"
  kubectl get secret "$secret" -o json -n default | jq .data | jq '. |= with_entries(.value |= @base64d)' > "$secret.json"
done