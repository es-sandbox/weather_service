#!/bin/bash

PORT=${1:-9000}

curl -XPOST localhost:$PORT/data -H "application/json" -d \
    '{"TempOUT": 22, "Humidity": 12, "TempIN": 23.50,  "Pressure": 810.22, "WindSpeed": 12.50, "WindDirection": 12, "Rainfall": 10, "Battery": 72, "Thunder": 1, "Light": 52.2, "Charging": 1}'

curl -XPOST localhost:$PORT/data -H "application/json" -d \
    '{"TempOUT": 23, "Humidity": 11, "TempIN": 24.50,  "Pressure": 710.22, "WindSpeed": 13.50, "WindDirection": 13, "Rainfall": 11, "Battery": 72, "Thunder": 0, "Light": 62.2, "Charging": 1}'

curl -XPOST localhost:$PORT/data -H "application/json" -d \
    '{"TempOUT": 24, "Humidity": 10, "TempIN": 25.50,  "Pressure": 610.22, "WindSpeed": 14.50, "WindDirection": 11, "Rainfall": 12, "Battery": 72, "Thunder": 0, "Light": 72.2, "Charging": 1}'

curl -XPOST localhost:$PORT/data -H "application/json" -d \
    '{"TempOUT": 23, "Humidity": 13, "TempIN": 21.50,  "Pressure": 910.22, "WindSpeed": 15.50, "WindDirection": 10, "Rainfall": 13, "Battery": 72, "Thunder": 1, "Light": 82.2, "Charging": 1}'

curl -XPOST localhost:$PORT/data -H "application/json" -d \
    '{"TempOUT": 22, "Humidity": 14, "TempIN": 19.50,  "Pressure": 930.22, "WindSpeed": 11.50, "WindDirection": 9, "Rainfall": 19, "Battery": 72, "Thunder": 1, "Light": 52.2, "Charging": 1}'