#!/bin/bash

IP=http://$(ifconfig docker0 | grep 'inet end' | awk '{print $3}')

if [[ ! -z $1 ]]; then
	IP=$1
fi

SHIPYARD_URL=$IP:8000

if [[ $UID != 0 ]]; then
  echo "Please run this script with sudo:"
  echo "sudo $0 $*"
  exit 1
fi

API_VERSION=${API_VERSION:-v1.12}

echo "Registering agent"
$(./shipyard-agent -url $SHIPYARD_URL -register &> output.tmp)
KEY=$(cat output.tmp | tail -1 | awk '{ print $5 }')

echo "Running agent"
./shipyard-agent -url $SHIPYARD_URL -key $KEY -api-version $API_VERSION
