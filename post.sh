#!/bin/bash

PORT=${1:-9000}

curl -XPOST localhost:$PORT/data -H "application/json" -d \
    '{"TempOUT": 22, "Humidity": 12, "TempIN": 22.50,  "Pressure": 910.22, "WindSpeed": 12.50, "WindDirection": 12, "Rainfall": 10, "Battery": 72, "Thunder": 1, "Light": 52.2, "Charging": 1, "Fire": true, "Snow": true, "SOS": true, "Latitude": 47.6414021, "Longitude": 26.2440167}'
