var xhr = new XMLHttpRequest();

xhr.open('GET', 'http://localhost:8080/data.json', true);

// xhr.setRequestHeader("Access-Control-Allow-Origin", "*");

xhr.send(); // (1)

var x = [], y = [];

xhr.onreadystatechange = function() { // (3)
    if (xhr.readyState != 4) return;

    if (xhr.status != 200) {
        console.log(xhr.status + ': ' + xhr.statusText);
    } else {
        console.log(xhr.responseText);

        var weatherInfo = JSON.parse(xhr.responseText);

        weatherInfo.forEach(function(item, i, arr) {
            x.push(item.ID);
            y.push(item.Temp);
        });

        TESTER = document.getElementById('tester');
        Plotly.plot( TESTER, [ {x, y} ], { margin: { t: 0 } } );
    }
};