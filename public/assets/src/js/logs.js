var term;

$(document).ready(function () {
    let data = {
        "InsuranceProductsGUID": "8a395dc5-b9e7-4864-8c18-c9244fcb5546",
        "policeGUID": "9d539dac-f283-4ed1-92dd-0d5fb3efeee1",
        "IIN": "4201066M072PB2",
        "UrlResult": "aHR0cHM6Ly9kZW5naS5tdHMuYnkvc2VydmljZXMvaW5zdXJhbmNlL2t1cGFsYS9yZXN1bHQ="
    }

    //redirectWithPost("http://kupala.test.francysk.com/onlineoferta/", data)
    if (location.pathname === "/logs/") {
        if (typeof term === "undefined") {
            term = new Terminal({ cols: 178, rows: 40 });
            term.open(document.getElementById('terminal'));
        }
        SubscribeOnRuntimeLogs();
    }
});

function redirectWithPost(url, data) {
    var form = document.createElement('form');
    form.method = 'POST';
    form.action = url;

    for (var key in data) {
        var input = document.createElement('input');
        input.name = key;
        input.value = data[key];
        input.type = 'hidden';
        form.appendChild(input)
    }
    document.body.appendChild(form);
    form.submit();
}



$('body').on('click', '.btn-clear-logs', function (e) {
    e.preventDefault();
    term.clear();
});

$('body').on('click', '.btn-remove-runtime-logs', function (e) {
    e.preventDefault();
    ToggleAjaxPreloader();
    let ws = webSocketConnection("ws://localhost:8888/logs/remove/runtime/");
    ws.onopen = function () {
        ws.send(JSON.stringify({"action": "runtime"}));
    };
    ws.onmessage = function (msg) {
        ToggleAjaxPreloader()
        let body = JSON.parse(msg.data);
        showNotification(body["status"], body["message"]);
        if (body["status"] === "success") {
            term.clear();
        }
    };
});

$('body').on('click', '.btn-remove-all-logs', function (e) {
    e.preventDefault();
    ToggleAjaxPreloader();
    let ws = webSocketConnection("ws://localhost:8888/logs/remove/all/");
    ws.onopen = function () {
        ws.send(JSON.stringify({"action": "all"}));
    };
    ws.onmessage = function (msg) {
        ToggleAjaxPreloader()
        let body = JSON.parse(msg.data);
        showNotification(body["status"], body["message"]);
        if (body["status"] === "success") {
            term.clear();
        }
    };
});

function SubscribeOnRuntimeLogs() {
    let ws = webSocketConnection("ws://localhost:8888/logs/subscribe/");
    ws.onopen = function () {
        ws.send(JSON.stringify({"action": "init"}));
    };
    ws.onmessage = function (msg) {
        term.writeln(msg.data)
    };
}