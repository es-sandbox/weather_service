#!/bin/bash

curl -XPOST localhost:9000/data -H "application/json" -d \
    '{"Temp": 43, "Humidity": 42, "Pressure": 42, "WindSpeed": 42, "WindDirection": "North", "Rainfall": 100500}'
