#!/bin/bash

curl -XPOST localhost:9000/data -H "application/json" -d \
    '{"Temp": 20, "Humidity": 42, "Pressure": 100, "WindSpeed": 42, "WindDirection": "North", "Rainfall": 32}'

curl -XPOST localhost:9000/data -H "application/json" -d \
    '{"Temp": 25, "Humidity": 43, "Pressure": 105, "WindSpeed": 41, "WindDirection": "North", "Rainfall": 33}'

curl -XPOST localhost:9000/data -H "application/json" -d \
    '{"Temp": 20, "Humidity": 44, "Pressure": 110, "WindSpeed": 40, "WindDirection": "North", "Rainfall": 34}'

curl -XPOST localhost:9000/data -H "application/json" -d \
    '{"Temp": 30, "Humidity": 45, "Pressure": 115, "WindSpeed": 39, "WindDirection": "North", "Rainfall": 35}'

curl -XPOST localhost:9000/data -H "application/json" -d \
    '{"Temp": 35, "Humidity": 39, "Pressure": 95, "WindSpeed": 38, "WindDirection": "North", "Rainfall": 36}'