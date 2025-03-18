#!/bin/bash

# Replace all the <...> placeholders with the actual parameter values to the .sh.tftpl file.
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
export REGION=<region>
export RESOURCE_GROUP=<resource_group_name>
export CLUSTER_NAME=<k8s_cluster_name> # e.g. "adr-cluster"
export AIO_CLUSTER_NAME=<aio_cluster_resource_name> # e.g. "adr-aio"
export AIO_CLUSTER_CUSTOM_LOCATION_NAME=adr-cl-1
export SUBSCRIPTION_ID=<subscription_id>
export TENANT_ID=<tenant_id>
export SCHEMA_REGISTRY_NAME=<schema_registry_name> # e.g. "adr-sr"
export SCHEMA_REGISTRY_NAMESPACE=<schema_registry_namespace> # e.g. "adr-sr-ns"
export STORAGE_ACCOUNT_NAME=<storage_account_name> # e.g. "adrstgacct", max 24 characters long.
export USER_ASSIGNED_MI_NAME=<managed_identity_name> # e.g. "adr-mi-1"
export KEYVAULT_NAME=<keyvault_name> # e.g. "adr-kv"

export AZURE_CLIENT_ID=<client_id> # managed identity or service principal client id
export AZURE_CLIENT_SECRET=<client_secret> # managed identity or service principal client secret

curl -sfL https://get.k3s.io | sh -

curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash

sudo chmod 777 /etc/rancher/k3s/k3s.yaml

az login
# az login --identity --client-id $AZURE_CLIENT_ID # Managed Identity login
# az login --client-id $AZURE_CLIENT_ID --password $AZURE_CLIENT_SECRET --service-principal --tenant $TENANT_ID # Service Principal login

az account set --subscription "$SUBSCRIPTION_ID"

az extension add --name connectedk8s

az extension add --name azure-iot-ops

az connectedk8s connect --name $CLUSTER_NAME -l $REGION --resource-group $RESOURCE_GROUP --enable-oidc-issuer --enable-workload-identity

SERVICE_ACCOUNT_ISSUER=$(az connectedk8s show --resource-group $RESOURCE_GROUP --name $CLUSTER_NAME --query oidcIssuerProfile.issuerUrl --output tsv)
sudo tee /etc/rancher/k3s/config.yaml > /dev/null <<EOF
kube-apiserver-arg:
- service-account-issuer=$SERVICE_ACCOUNT_ISSUER
- service-account-max-token-expiration=24h
EOF

systemctl daemon-reload
systemctl restart k3s

export OBJECT_ID=$(az ad sp show --id bc313c14-388c-4e7d-a58e-70017303ee3b --query id -o tsv)
az connectedk8s enable-features -n $CLUSTER_NAME -g $RESOURCE_GROUP --custom-locations-oid $OBJECT_ID --features cluster-connect custom-locations

az storage account create --name $STORAGE_ACCOUNT_NAME --resource-group $RESOURCE_GROUP --enable-hierarchical-namespace --verbose

az iot ops schema registry create --name $SCHEMA_REGISTRY_NAME --resource-group $RESOURCE_GROUP --registry-namespace $SCHEMA_REGISTRY_NAMESPACE --sa-resource-id $(az storage account show --name $STORAGE_ACCOUNT_NAME --resource-group $RESOURCE_GROUP -o tsv --query id) --verbose

az iot ops init  --subscription $SUBSCRIPTION_ID -g $RESOURCE_GROUP --cluster $CLUSTER_NAME --debug --no-progress

az iot ops create  --subscription $SUBSCRIPTION_ID  -g $RESOURCE_GROUP  --cluster $CLUSTER_NAME --custom-location $AIO_CLUSTER_CUSTOM_LOCATION_NAME  -n $AIO_CLUSTER_NAME  --sr-resource-id $(az iot ops schema registry list -g $RESOURCE_GROUP --query "[?name=='$SCHEMA_REGISTRY_NAME'].id" -o tsv)  --add-insecure-listener --enable-rsync --debug --yes --no-progress

az identity create --name $USER_ASSIGNED_MI_NAME --resource-group $RESOURCE_GROUP --location $REGION --subscription $SUBSCRIPTION_ID

az keyvault create --resource-group $RESOURCE_GROUP --location $REGION --name $KEYVAULT_NAME --enable-rbac-authorization

az role assignment create --role "Key Vault Secrets Officer" --assignee $(az ad signed-in-user show --query id -o tsv) --scope /subscriptions/$SUBSCRIPTION_ID/resourcegroups/$RESOURCE_GROUP/providers/Microsoft.KeyVault/vaults/$KEYVAULT_NAME

az iot ops secretsync enable  -g $RESOURCE_GROUP  -n $AIO_CLUSTER_NAME  --kv-resource-id $(az keyvault list -g $RESOURCE_GROUP  --query "[?name=='$KEYVAULT_NAME'].id" -o tsv) --mi-user-assigned $(az identity list -g $RESOURCE_GROUP --query  "[?name=='$USER_ASSIGNED_MI_NAME'].id" -o tsv)

az iot ops identity assign  -g $RESOURCE_GROUP  -n $AIO_CLUSTER_NAME  --mi-user-assigned $(az identity list -g $RESOURCE_GROUP --query  "[?name=='$USER_ASSIGNED_MI_NAME'].id" -o tsv)
