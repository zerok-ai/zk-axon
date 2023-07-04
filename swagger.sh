#!/bin/bash
THIS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
export ROOT_DIR="$THIS_DIR"

echo $ROOT_DIR
rm -rf $ROOT_DIR/swagger
mkdir $ROOT_DIR/swagger
mkdir $ROOT_DIR/swagger/docs
swag init --pd -dir $ROOT_DIR/internal/scenarioDataPersistence/ -g routes.go --instanceName scenarioDataPersistence -o $ROOT_DIR/swagger/docs
jq -s '.[0]' $ROOT_DIR/swagger/docs/*.json >> $ROOT_DIR/swagger/combined.json
