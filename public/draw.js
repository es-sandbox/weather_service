var xhr = new XMLHttpRequest();

xhr.open('GET', 'http://192.168.0.105:9000/data', true);

// xhr.setRequestHeader("Access-Control-Allow-Origin", "*");

xhr.send(); // (1)

xhr.onreadystatechange = function() { // (3)
    if (xhr.readyState != 4) return;

    if (xhr.status != 200) {
        console.log(xhr.status + ': ' + xhr.statusText);
    } else {
        console.log(xhr.responseText);

        var weatherInfo = JSON.parse(xhr.responseText);

        var xTempOUT = [], yTempOUT = [];
        var xHumidity = [], yHumidity = [];
        var xTempIN = [], yTempIN = [];
        var xPressure = [], yPressure = [];
        var xWindSpeed = [], yWindSpeed = [];
        var xWindDirection = [], yWindDirection = [];
        var xRainfall = [], yRainfall = [];

        weatherInfo.forEach(function(item, i, arr) {
            xTempOUT.push(item.ID);
            yTempOUT.push(item.TempOUT);

            xHumidity.push(item.ID);
            yHumidity.push(item.Humidity);

            xTempIN.push(item.ID);
            yTempIN.push(item.TempIN);

            xPressure.push(item.ID);
            yPressure.push(item.Pressure);

            xWindSpeed.push(item.ID);
            yWindSpeed.push(item.WindSpeed);

            xWindDirection.push(item.ID);
            yWindDirection.push(item.WindDirection);

            xRainfall.push(item.ID);
            yRainfall.push(item.Rainfall);
        });

        TESTER = document.getElementById('TempOUT');
        Plotly.plot( TESTER, [ {x: xTempOUT, y: yTempOUT} ], { margin: { t: 0 } } );

        TESTER = document.getElementById('Humidity');
        Plotly.plot( TESTER, [ {x: xHumidity, y: yHumidity} ], { margin: { t: 0 } } );

        TESTER = document.getElementById('TempIN');
        Plotly.plot( TESTER, [ {x: xTempIN, y: yTempIN} ], { margin: { t: 0 } } );

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