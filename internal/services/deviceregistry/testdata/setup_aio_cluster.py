# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This python script is used to arc enable and install Azure IoT Operations extension.
# This should be identical to the script used in the AIO onboarding flow, plus setting custom location enabled.
# See this link for the exact commands: https://review.learn.microsoft.com/en-us/azure/iot-operations/get-started-end-to-end-sample/quickstart-deploy?branch=main

import argparse
import json
import logging as logger
import os
import platform
import shutil
import signal
import stat
import subprocess
import time
from subprocess import PIPE, Popen

def register_az_providers():
    print("Registering required Azure providers...")
    required_providers = [
        "Microsoft.ExtendedLocation",
        "Microsoft.Kubernetes",
        "Microsoft.KubernetesConfiguration",
        "Microsoft.IoTOperations",
        "Microsoft.DeviceRegistry",
        "Microsoft.SecretSyncController"
    ]
    for provider in required_providers:
        print("Registering provider " + provider + "...")
        
        # `az provider register -n <provider> --only-show-errors`
        register_provider_cmd = f"az provider register -n {provider} --only-show-errors"
        response = subprocess.run(register_provider_cmd, stdout=PIPE, stderr=PIPE, shell=True)
        if response.returncode != 0:
            raise Exception("Failed to register provider " + provider + ": " + str(response.stderr))
        print("Successfully registered provider " + provider)
    
    print("Successfully registered required Azure providers")
    return

def install_az_extensions():
    print("Adding required Azure CLI extension...")
    required_extensions = [
        "aksarc",
        "azure-iot-ops",
        "connectedk8s",
        "customlocation",
        "k8s-extension"
    ]
    for extension in required_extensions:
        print("Installing extension " + extension + "...")
        
        # `az extension add --upgrade --name <extension> --only-show-errors`
        install_extension_cmd = f"az extension add --upgrade --name {extension} --only-show-errors"
        response = subprocess.run(install_extension_cmd, stdout=PIPE, stderr=PIPE, shell=True)
        if response.returncode != 0:
            raise Exception("Failed to install extension " + extension + ": " + str(response.stderr))
        print("Successfully installed extension " + extension)
    
    print("Successfully installed required Azure extensions")
    return

def onboard_k8s_cluster(resource_group_name, cluster_name, location):
    print("Connecting k8s cluster " + cluster_name + " to Azure...")
    
    # `az connectedk8s connect --name <cluster_name> --location <location> --resource-group <resource_group_name> --kube-config /etc/rancher/k3s/k3s.yaml --only-show-errors`
    connect_cluster_cmd = f"az connectedk8s connect --name {cluster_name} --location {location} --resource-group {resource_group_name} --kube-config /etc/rancher/k3s/k3s.yaml --only-show-errors"
    response = subprocess.run(connect_cluster_cmd, stdout=PIPE, stderr=PIPE, shell=True)
    if response.returncode != 0:
        raise Exception("Failed to connect k8s cluster " + cluster_name + ": " + str(response.stderr))
    
    print("Successfully connected k8s cluster " + cluster_name + " to Azure. Getting object ID of Microsoft Entra ID Application for Azure Arc service...")

    # `az ad sp show --id bc313c14-388c-4e7d-a58e-70017303ee3b --query id -o tsv --only-show-errors`
    get_object_id_cmd = f"az ad sp show --id bc313c14-388c-4e7d-a58e-70017303ee3b --query id -o tsv --only-show-errors"
    response = subprocess.run(get_object_id_cmd, stdout=PIPE, stderr=PIPE, shell=True, text=True)
    if response.returncode != 0:
        raise Exception("Failed to get object ID of Microsoft Entra ID Application for Azure Arc service: " + str(response.stderr))
    object_id = response.stdout

    print("Successfully got object ID of Microsoft Entra ID Application for Azure Arc service. Enabling features on the cluster...")
    
    # `az connectedk8s enable-features -n <cluster_name> -g <resource_group> --custom-locations-oid <object_id> --features cluster-connect custom-locations --kube-config /etc/rancher/k3s/k3s.yaml --only-show-errors`
    enable_features_cmd = f"az connectedk8s enable-features -n {cluster_name} -g {resource_group_name} --custom-locations-oid {object_id.strip()} --features cluster-connect custom-locations --kube-config /etc/rancher/k3s/k3s.yaml --only-show-errors"
    response = subprocess.run(enable_features_cmd, stdout=PIPE, stderr=PIPE, shell=True)
    if response.returncode != 0:
        raise Exception("Failed to enable features on k8s cluster " + cluster_name + ": " + str(response.stderr))
    
    print("Successfully enabled features on k8s cluster " + cluster_name)
    return

def setup_schema_registry(storage_account, schema_registry, schema_registry_namespace, resource_group_name, location):
    print("Setting up storage account " + storage_account + " for schema registry " + schema_registry + "...")
    
    # `az storage account create --name <storage_account> --location <location> --resource-group <resource_group_name> --enable-hierarchical-namespace --only-show-errors`
    create_sa_cmd = f"az storage account create --name {storage_account} --location {location} --resource-group {resource_group_name} --enable-hierarchical-namespace --only-show-errors"
    response = subprocess.run(create_sa_cmd, stdout=PIPE, stderr=PIPE, shell=True)
    if response.returncode != 0:
        raise Exception("Failed to create storage account " + storage_account + ": " + str(response.stderr))

    print("Successfully created storage account " + storage_account + ". Setting up schema registry " + schema_registry)

    # `az storage account show --name <storage_account> -o tsv --query id --only-show-errors`
    get_sa_id_cmd = f"az storage account show --name {storage_account} -o tsv --query id --only-show-errors"
    response = subprocess.run(get_sa_id_cmd, stdout=PIPE, stderr=PIPE, shell=True, text=True)
    if response.returncode != 0:
        raise Exception("Failed to get storage account id of " + storage_account + ": " + str(response.stderr))
    storage_account_id = response.stdout

    # `az iot ops schema registry create --name <schema_registry> --resource-group <resource_group_name> --registry-namespace <schema_registry_namespace> --sa-resource-id <storage_account_id> --only-show-errors`
    create_sr_cmd = f"az iot ops schema registry create --name {schema_registry} --resource-group {resource_group_name} --registry-namespace {schema_registry_namespace} --sa-resource-id {storage_account_id.strip()} --only-show-errors"
    response = subprocess.run(create_sr_cmd, stdout=PIPE, stderr=PIPE, shell=True)
    if response.returncode != 0:
        raise Exception("Failed to create schema registry " + schema_registry + ": " + str(response.stderr))
    
    print("Successfully created schema registry " + schema_registry)
    return

def install_azure_iot_ops_extension(cluster_name, resource_group_name, schema_registry, custom_location, aio_create_timeout=15*60):
    print("Initializing cluster for installing Azure IoT Operations extension...")

    # `az iot ops init --cluster <cluster_name> --resource-group <resource_group_name> --no-progress --only-show-errors`
    aio_init_cmd = f"az iot ops init --cluster {cluster_name} --resource-group {resource_group_name} --no-progress --only-show-errors"
    response = subprocess.run(aio_init_cmd, stdout=PIPE, stderr=PIPE, shell=True)
    if response.returncode != 0:
        raise Exception("Failed to initialize Azure IoT Operations extension on cluster " + cluster_name + ": " + str(response.stderr))
    
    print("Successfully initialized cluster. Getting the schema registry resource ID...")

    # `az iot ops schema registry show --name <schema_registry> --resource-group <resource_group_name> -o tsv --query id --only-show-errors`
    get_sr_id_cmd = f"az iot ops schema registry show --name {schema_registry} --resource-group {resource_group_name} -o tsv --query id --only-show-errors"
    response = subprocess.run(get_sr_id_cmd, stdout=PIPE, stderr=PIPE, shell=True, text=True)
    if response.returncode != 0:
        raise Exception("Failed to get schema registry id of " + schema_registry + ": " + str(response.stderr))
    schema_registry_id = response.stdout

    print("Successfully got schema registry resource ID. Installing Azure IoT Operations extension...")
    # Install the Azure IoT Operations extension onto cluster.
        # Note: there is a known issue when trying to install the extension on a Kind cluster.
        # The schema registry portion of the extension will timeout after 30+ minutes and fail to install.
        # But this will not block other parts of AIO installation and the acceptance tests for device registry,
        # so set a timeout cancellation for this step to avoid blocking the rest of the acceptance tests.
    try:
        # `az iot ops create --cluster <cluster_name> --resource-group <resource_group_name> --name <cluster_name>-instance  --sr-resource-id <schema_registry_id> --broker-frontend-replicas 1 --broker-frontend-workers 1  --broker-backend-part 1  --broker-backend-workers 1 --broker-backend-rf 2 --broker-mem-profile Low --custom-location <custom_location> --no-progress --yes --only-show-errors`
        # Default timeout is 15 minutes.
        aio_create_cmd = f"az iot ops create --cluster {cluster_name} --resource-group {resource_group_name} --name {cluster_name}-instance  --sr-resource-id {schema_registry_id.strip()} --broker-frontend-replicas 1 --broker-frontend-workers 1  --broker-backend-part 1  --broker-backend-workers 1 --broker-backend-rf 2 --broker-mem-profile Low --custom-location {custom_location} --no-progress --yes --only-show-errors"
        response = subprocess.run(aio_create_cmd, timeout=aio_create_timeout, stdout=PIPE, stderr=PIPE, shell=True)
        print(response)
        if response.returncode != 0:
            raise Exception("Failed to install Azure IoT Operations extension on cluster " + cluster_name + ": " + str(response.stderr))
        print("Successfully installed Azure IoT Operations extension on cluster " + cluster_name)
        return
    except subprocess.TimeoutExpired:
        # Swallow timeout exception and continue with the rest of the acceptance tests
        print("Warning: installing Azure IoT Operations extension on cluster " + cluster_name + " timed out after " + str(aio_create_timeout) + " seconds but should be working for tests. Continuing with the rest of the acceptance tests")
        return
    except Exception as e:
        raise Exception("Failed to install Azure IoT Operations extension on cluster " + cluster_name + ": " + str(e))

def setup_aio_arc_enabled_cluster():
    parser = argparse.ArgumentParser(
        description='Install AIO extension onto cluster')
    parser.add_argument('--subscriptionId', type=str, required=True)
    parser.add_argument('--resourceGroupName', type=str, required=True)
    parser.add_argument('--clusterName', type=str, required=True)
    parser.add_argument('--location', type=str, required=True)
    parser.add_argument('--customLocation', type=str, required=True)
    parser.add_argument('--storageAccount', type=str, required=True)
    parser.add_argument('--schemaRegistry', type=str, required=True)
    parser.add_argument('--schemaRegistryNamespace', type=str, required=True)

    try:
        args = parser.parse_args()
    except Exception as e:
        raise Exception("Failed to parse arguments." + str(e))

    # register_az_providers()
    install_az_extensions()

    onboard_k8s_cluster(args.resourceGroupName, args.clusterName, args.location)

    setup_schema_registry(args.storageAccount, args.schemaRegistry, args.schemaRegistryNamespace, args.resourceGroupName, args.location)

    install_azure_iot_ops_extension(args.clusterName, args.resourceGroupName, args.schemaRegistry, args.customLocation)

if __name__ == "__main__":
    setup_aio_arc_enabled_cluster()
