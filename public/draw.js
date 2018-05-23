var xhr = new XMLHttpRequest();

drawChartEvent();

// setInterval(iter, 2000);

function drawChartEvent(type) {
    var endpoint = 'http://192.168.0.105:9000/data';
    if (type == 'dailyRecords') {
        endpoint = 'http://192.168.0.105:9000/data/last_day';
    }
    if (type == 'hourlyRecords') {
        endpoint = 'http://192.168.0.105:9000/data/last_hour';
    }
    if (type == 'lastMinuteRecords') {
        endpoint = 'http://192.168.0.105:9000/data/last_minute';
    }

    xhr.open('GET', endpoint, true);

    // xhr.setRequestHeader("Access-Control-Allow-Origin", "*");

    xhr.send(); // (1)

    xhr.onreadystatechange = function () { // (3)
        if (xhr.readyState != 4) return;

        if (xhr.status != 200) {
            console.log(xhr.status + ': ' + xhr.statusText);
        } else {
            // console.log(xhr.responseText);

            var weatherInfo = JSON.parse(xhr.responseText);
            drawTemperatureOutsideChart(weatherInfo);
            drawHumidityChart(weatherInfo);
            drawTemperatureInsideChart(weatherInfo);
            drawPressureChart(weatherInfo);
            drawWindSpeedChart(weatherInfo);
            drawWindDirectionChart(weatherInfo);
            drawRainfallChart(weatherInfo);
        }
    };
}

function drawTemperatureOutsideChart(arr) {
    var xTempOutside = [], yTempOutside = [];
    arr.forEach(function (item, i, arr) {
        xTempOutside.push(item.ID);
        yTempOutside.push(item.TempOUT)
    });

    var layout = { xaxis: { title: 'Timeline' }, yaxis: { title: 'Temperature Outside' } };
    var tempOutside = document.getElementById('temperature-outside');
    Plotly.purge(tempOutside);
    Plotly.plot(tempOutside, [ {x: xTempOutside, y: yTempOutside} ], layout);
}

function drawHumidityChart(arr) {
    var xHumidity = [], yHumidity = [];
    arr.forEach(function (item, i, arr) {
        xHumidity.push(item.ID);
        yHumidity.push(item.Humidity);
    });

    var layout = { xaxis: { title: 'Timeline' }, yaxis: { title: 'Humidity' } };
    var humidity = document.getElementById('humidity');
    Plotly.purge(humidity);
    Plotly.plot(humidity, [ {x: xHumidity, y: yHumidity} ], layout);
}

function drawTemperatureInsideChart(arr) {
    var xTempInside = [], yTempInside = [];
    arr.forEach(function (item, i, arr) {
        xTempInside.push(item.ID);
        yTempInside.push(item.TempIN);
    });

    var layout = { xaxis: { title: 'Timeline' }, yaxis: { title: 'Temperature Inside' } };
    var tempInside = document.getElementById('temperature-inside');
    Plotly.purge(tempInside);
    Plotly.plot(tempInside, [ {x: xTempInside, y: yTempInside} ], layout);
}

function drawPressureChart(arr) {
    var xPressure = [], yPressure = [];
    arr.forEach(function (item, i, arr) {
        xPressure.push(item.ID);
        yPressure.push(item.Pressure);
    });

    var layout = { xaxis: { title: 'Timeline' }, yaxis: { title: 'Pressure' } };
    var pressure = document.getElementById('pressure');
    Plotly.purge(pressure);
    Plotly.plot(pressure, [ {x: xPressure, y: yPressure} ], layout);
}

function drawWindSpeedChart(arr) {
    var xWindSpeed = [], yWindSpeed = [];
    arr.forEach(function (item, i, arr) {
        xWindSpeed.push(item.ID);
        yWindSpeed.push(item.WindSpeed);
    });

    var layout = { xaxis: { title: 'Timeline' }, yaxis: { title: 'WindSpeed' } };
    var windSpeed = document.getElementById('wind-speed');
    Plotly.purge(windSpeed);
    Plotly.plot(windSpeed, [ {x: xWindSpeed, y: yWindSpeed} ], layout);
}

function drawWindDirectionChart(arr) {
    var xWindDirection = [], yWindDirection = [];
    arr.forEach(function (item, i, arr) {
        xWindDirection.push(item.ID);
        yWindDirection.push(item.WindDirection);
    });

    var layout = { xaxis: { title: 'Timeline' }, yaxis: { title: 'WindDirection' } };
    var windDirection = document.getElementById('wind-direction');
    Plotly.purge(windDirection);
    Plotly.plot(windDirection, [ {x: xWindDirection, y: yWindDirection} ], layout);
}

function drawRainfallChart(arr) {
    var xRainfall = [], yRainfall = [];
    arr.forEach(function (item, i, arr) {
        xRainfall.push(item.ID);
        yRainfall.push(item.Rainfall);
    });

    var layout = { xaxis: { title: 'Timeline' }, yaxis: { title: 'Rainfall' } };
    var rainfall = document.getElementById('rainfall');
    Plotly.purge(rainfall);
    Plotly.plot(rainfall, [ {x: xRainfall, y: yRainfall} ], layout);
}