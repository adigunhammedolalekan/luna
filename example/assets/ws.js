
var ws = new WebSocket("ws://localhost:8009/ws/connect");
var roomId = "22"; //Should be a real/dynamic room ID for a real application
ws.onopen = function () {

    var message = {
        "action" : "subscribe",
        "path" : "/rooms/" + roomId + "/message"
    };

    var json = JSON.stringify(message);
    ws.send(json)
};

ws.onmessage = function (data) {
    console.log(data.data);

    message = JSON.parse(data.data);
    document.getElementById("content").innerHTML += message.message + "<br>";
};

ws.onerror = function (err) {
    console.log("Error => " + err)
};

function sendMessage() {

    text = document.getElementById("chat_box").value;
    var message = {
        "action" : "message",
        "path" : "/rooms/" + roomId + "/message",
        "data" : {
            "message" : text,
            "others" : "Other json data here"
        }
    };

    document.getElementById("chat_box").value = "";
    var json = JSON.stringify(message);
    ws.send(json)
}