#!/bin/bash

set -e
docker build -t ayufan/docker-composer:master .
docker push ayufan/docker-composer:master
