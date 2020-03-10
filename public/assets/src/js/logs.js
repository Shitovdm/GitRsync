var term;

$(document).ready(function () {
    if (location.pathname === "/logs/") {
        if (typeof term === "undefined") {
            term = new ExecTerminal('terminal');
            term.UpdateTerminalFit();
        }
        SubscribeOnRuntimeLogs()
    }
});

function SubscribeOnRuntimeLogs() {
    let ws = webSocketConnection("ws://localhost:8888/logs/subscribe/");
    ws.onopen = function()
    {
        ws.send(JSON.stringify({"action":"init"}));
        //ws.send(JSON.stringify({"action":"subscribe"}));
    };
    ws.onmessage = function(msg) {
        term.terminal.writeln(msg.data)
    };
}