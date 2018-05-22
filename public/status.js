var xhr = new XMLHttpRequest();

iter();

setInterval(iter, 2000);

function iter() {
    xhr.open('GET', 'http://192.168.0.105:9000/data', true);

    // xhr.setRequestHeader("Access-Control-Allow-Origin", "*");

    xhr.send(); // (1)

    xhr.onreadystatechange = function() { // (3)
        if (xhr.readyState != 4) return;

        if (xhr.status != 200) {
            console.log(xhr.status + ': ' + xhr.statusText);
        } else {
            // console.log(xhr.responseText);

            var weatherInfo = JSON.parse(xhr.responseText);
            lastRecord = weatherInfo.pop()

            console.log("TimeStamp " + lastRecord.TimeStamp);
            console.log("Battery " + lastRecord.Battery);
            console.log("Charging " + lastRecord.Charging);

            console.log(convertTimeStampToData(lastRecord.TimeStamp));

            document.getElementById("timestamp").innerHTML = convertTimeStampToData(lastRecord.TimeStamp)
            document.getElementById("battery").innerHTML = lastRecord.Battery
            document.getElementById("charging").innerHTML = convertIntToBool(lastRecord.Charging)
        }
    }
}

function convertIntToBool(n) {
    if (n == 0) {
        return false
    }
    return true
}

function convertTimeStampToData(unix_timestamp) {
    // Create a new JavaScript Date object based on the timestamp
    // multiplied by 1000 so that the argument is in milliseconds, not seconds.
    var date = new Date(unix_timestamp / 1000000);
    return date.toString();
}