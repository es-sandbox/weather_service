#!/bin/bash

curl -XPOST localhost:9000/data -H "application/json" -d '{"Temp": 42}'
