var term;

$(document).ready(function () {
    if (location.pathname === "/logs/") {
        if (typeof term === "undefined") {
            term = new Terminal({ cols: 178, rows: 40 });
            term.open(document.getElementById('terminal'));
        }
        SubscribeOnRuntimeLogs();
    }
});

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