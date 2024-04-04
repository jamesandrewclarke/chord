#!/bin/bash
# Gets credentials to use as secrets for GitHub Actions
# Allows pushing and pulling from the ACR repository

ACR_NAME=chord
SERVICE_PRINCIPAL_NAME=gh-actions
SUBSCRIPTION=bc56aec8-74ab-408c-bf65-f69efd5d2446

# Obtain the full registry ID
ACR_REGISTRY_ID=$(az acr show --name $ACR_NAME --query "id" --output tsv --subscription $SUBSCRIPTION)
echo "Registry ID: $ACR_REGISTRY_ID"

PASSWORD=$(az ad sp create-for-rbac --name $SERVICE_PRINCIPAL_NAME --scopes $ACR_REGISTRY_ID --role acrpull --query "password" --output tsv --subscription $SUBSCRIPTION)
USER_NAME=$(az ad sp list --display-name $SERVICE_PRINCIPAL_NAME --query "[].appId" --output tsv --subscription $SUBSCRIPTION)

# Output the service principal's credentials; use these in your services and
# applications to authenticate to the container registry.
echo "Service principal ID: $USER_NAME"
echo "Service principal password: $PASSWORD"