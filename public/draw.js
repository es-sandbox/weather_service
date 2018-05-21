var xhr = new XMLHttpRequest();

xhr.open('GET', 'http://localhost:8080/api', true);

// xhr.setRequestHeader("Access-Control-Allow-Origin", "*");

xhr.send(); // (1)

xhr.onreadystatechange = function() { // (3)
    if (xhr.readyState != 4) return;

    if (xhr.status != 200) {
        console.log(xhr.status + ': ' + xhr.statusText);
    } else {
        console.log(xhr.responseText);

        var weatherInfo = JSON.parse(xhr.responseText);

        var xTemp = [], yTemp = [];
        var xHumidity = [], yHumidity = [];
        var xPressure = [], yPressure = [];
        var xWindSpeed = [], yWindSpeed = [];
        var xWindDirection = [], yWindDirection = [];
        var xRainfall = [], yRainfall = [];

        weatherInfo.forEach(function(item, i, arr) {
            xTemp.push(item.ID);
            yTemp.push(item.Temp);

            xHumidity.push(item.ID);
            yHumidity.push(item.Humidity);

            xPressure.push(item.ID);
            yPressure.push(item.Pressure);

            xWindSpeed.push(item.ID);
            yWindSpeed.push(item.WindSpeed);

            xWindDirection.push(item.ID);
            yWindDirection.push(item.WindDirection);

            xRainfall.push(item.ID);
            yRainfall.push(item.Rainfall);
        });

        TESTER = document.getElementById('Temp');
        Plotly.plot( TESTER, [ {x: xTemp, y: yTemp} ], { margin: { t: 0 } } );

        TESTER = document.getElementById('Humidity');
        Plotly.plot( TESTER, [ {x: xHumidity, y: yHumidity} ], { margin: { t: 0 } } );

        TESTER = document.getElementById('Pressure');
        Plotly.plot( TESTER, [ {x: xPressure, y: yPressure} ], { margin: { t: 0 } } );

        TESTER = document.getElementById('WindSpeed');
        Plotly.plot( TESTER, [ {x: xWindSpeed, y: yWindSpeed} ], { margin: { t: 0 } } );

        TESTER = document.getElementById('WindDirection');
        Plotly.plot( TESTER, [ {x: xWindDirection, y: yWindDirection} ], { margin: { t: 0 } } );

        TESTER = document.getElementById('Rainfall');
        Plotly.plot( TESTER, [ {x: xRainfall, y: yRainfall} ], { margin: { t: 0 } } );
    }
};