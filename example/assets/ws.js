
ws = new WebSocket("ws://192.168.43.39:8009/ws/connect");
ws.onopen = function () {

    var message = {
        "action" : "subscribe",
        "path" : "/rooms/22/message"
    };

    var json = JSON.stringify(message);
    ws.send(json)
};

ws.onmessage = function (data) {
    console.log(data)
};

ws.onerror = function (err) {
    console.log("Error => " + err)
};