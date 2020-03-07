var term;

$(document).ready(function () {
    if (location.pathname === "/logs/") {
        getLogs()
    }
});


function getLogs() {
    if (typeof term === "undefined") {
        term = new ExecTerminal('terminal');
        term.UpdateTerminalFit();
    }
    term.terminal.writeln("hello world!")

    let ws = webSocketConnection("ws://localhost:8888/logs/append");
    ws.onopen = function()
    {
        ws.send(JSON.stringify({}));
    };
    ws.onmessage = function(msg) {
        term.terminal.writeln(msg.data)
    };

}

