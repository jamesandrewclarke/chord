#!/bin/bash
# Gets credentials to use as secrets for GitHub Actions
# Allows pushing and pulling from the ACR repository

set -ex

ACR_NAME=chord
SERVICE_PRINCIPAL_NAME=gh-actions
SUBSCRIPTION=66f773f7-ab57-4a46-9c43-c11a2caf7c9d

# Obtain the full registry ID
ACR_REGISTRY_ID=$(az acr show --name $ACR_NAME --query "id" --output tsv --subscription $SUBSCRIPTION)
echo "Registry ID: $ACR_REGISTRY_ID"

PASSWORD=$(az ad sp create-for-rbac --name $SERVICE_PRINCIPAL_NAME --scopes $ACR_REGISTRY_ID --role contributor --query "password" --output tsv)
USER_NAME=$(az ad sp list --display-name $SERVICE_PRINCIPAL_NAME --query "[].appId" --output tsv)

# Output the service principal's credentials; use these in your services and
# applications to authenticate to the container registry.
echo "Service principal ID: $USER_NAME"
echo "Service principal password: $PASSWORD"