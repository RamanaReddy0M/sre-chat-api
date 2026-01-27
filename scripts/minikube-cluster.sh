#!/usr/bin/env bash
# Minikube 3-node production cluster setup for SRE Chat API
# Node A (application): type=application
# Node B (database): type=database
# Node C (dependent_services): type=dependent_services (observability, vault, etc.)

set -e

PROFILE="${MINIKUBE_PROFILE:-sre-chat-api}"
NODES=3

echo "==> Starting Minikube cluster (profile: ${PROFILE}, nodes: ${NODES})..."
minikube start \
  --profile="${PROFILE}" \
  --nodes="${NODES}" \
  --driver="${MINIKUBE_DRIVER:-docker}" \
  --memory="${MINIKUBE_MEMORY:-4096}" \
  --cpus="${MINIKUBE_CPUS:-2}"

echo "==> Waiting for all nodes to be Ready..."
minikube --profile="${PROFILE}" kubectl -- wait --for=condition=Ready nodes --all --timeout=120s

echo "==> Applying node labels..."

# Node C: control-plane node — dependent services (observability, vault, etc.)
minikube --profile="${PROFILE}" kubectl -- label nodes "${PROFILE}" type=dependent_services --overwrite

# Node A: first worker — application
minikube --profile="${PROFILE}" kubectl -- label nodes "${PROFILE}-m02" type=application --overwrite

# Node B: second worker — database
minikube --profile="${PROFILE}" kubectl -- label nodes "${PROFILE}-m03" type=database --overwrite

echo ""
echo "==> Cluster ready. Node labels:"
minikube --profile="${PROFILE}" kubectl -- get nodes --show-labels
echo ""
echo "Node roles:"
echo "  - ${PROFILE}       : type=dependent_services (observability, vault, etc.)"
echo "  - ${PROFILE}-m02   : type=application"
echo "  - ${PROFILE}-m03   : type=database"
