#!/usr/bin/env bash

set -euf -o pipefail

group="tissue"
location='westus'
cluster_name="tissue-k8s-cluster"
name_prefix="k8s"
new_name=$(echo $(mktemp -u ${name_prefix}XXXX) | tr '[:upper:]' '[:lower:]')
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [[ -z $(az account list -o tsv 2>/dev/null ) ]]; then
    az login -o table
    echo ""
fi

if [[ -z $(az group show -n ${group}) ]]; then
    echo "Creating Resource Group named ${group}"
    az group create -n ${group} -l ${location} 1>/dev/null
else 
    echo "Using Resource Group name ${group}"
fi

if [[ ! -f ${DIR}/.ssh/id_rsa ]]; then
    mkdir -p ${DIR}/.ssh
    echo "Generating ssh keys to use for setting up the Kubernetes cluster"
    ssh-keygen -f ${DIR}/.ssh/id_rsa -t rsa -N '' 1>/dev/null
else
    echo "Using ${DIR}/.ssh/id_rsa to authenticate with the Kubernetes cluster"
fi

cosmosdb_name=$(az cosmosdb list -g ${group} --query "[?starts_with(name, '${name_prefix}')] | [0].name" -o tsv)
if [[ -z "${cosmosdb_name}" ]]; then
    echo "Create CosmosDB instance named ${new_name} with MongoDB wire format"
    az cosmosdb create -g ${group} -n "${new_name}" --kind MongoDB 1>/dev/null
else
    echo "Using the already provisioned CosmosDB named ${cosmosdb_name}"
fi

if [[ -z $(az acs show -g ${group} -n ${cluster_name} -o tsv) ]]; then
    echo "Creating Azure Kubernetes cluster named ${cluster_name} in group ${group}"
    az acs create -g ${group} -n ${cluster_name} --orchestrator-type Kubernetes \
        --ssh-key-value ${DIR}/.ssh/id_rsa.pub 1>/dev/null
else
    echo "Using Azure Kubernetes cluster named ${cluster_name} in group ${group}"
fi

registry=$(az acr list -g ${group} --query "[?starts_with(name, '${name_prefix}')] | [0].loginServer" -o tsv)
if [[ -z ${registry} ]]; then
    echo "Creating Azure Container Registry named ${new_name} in group ${group}"
    registry=$(az acr create -g ${group} -n ${new_name} -l ${location} \
                --admin-enabled true --sku Basic --query "loginServer" -o tsv)
else
    echo "Using Azure Container Registry ${registry} in group ${group}"
fi
echo ""

registry_name=$(az acr list -g ${group} --query "[?starts_with(name, '${name_prefix}')] | [0].name" -o tsv)
creds=$(az acr credential show -g ${group} -n ${registry_name})
username=$(echo $creds | jq ".username" -r)
password=$(echo $creds | jq ".passwords[0].value" -r)
echo "Logging Docker into ${registry} with user: ${username}"
docker login ${registry} -u ${username} -p ${password}
echo "To push to your docker registry run 'docker push ${registry}/myImage:version'"
echo ""

if [[ ! -d ${HOME}.kube/config ]]; then
    echo "Creating ${HOME}.kube/config w/ credentials for managing ${cluster_name}"
    az acs kubernetes get-credentials -g ${group} -n ${cluster_name} \
        --ssh-key-file ${DIR}/.ssh/id_rsa 1>/dev/null
else
    echo "Using ${HOME}.kube/config w/ credentials for managing ${cluster_name}"
fi

if [[ -z $(which kubectl) ]]; then
    echo "We didn't find kubectl installed, so installing..."
    az acs kubernetes install-cli
fi

echo "Your Kubernetes cluster has been deployed and you are ready to connect."
echo "To connect to the cluster run 'kubectl cluster-info'"
echo ""