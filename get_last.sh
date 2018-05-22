#!/bin/bash

PORT=${1:-9000}

curl -XGET localhost:$PORT/data/last
echo
