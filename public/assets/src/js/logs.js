var term;

$(document).ready(function () {
    if (location.pathname === "/logs/") {
        if (typeof term === "undefined") {
            term = new ExecTerminal('terminal');
            term.UpdateTerminalFit();
        }
        term.terminal.writeln("hello world!")
        getLogs()
    }
});


function getLogs() {



    let ws = webSocketConnection("ws://localhost:8888/logs/process/");
    ws.onopen = function()
    {
        ws.send(JSON.stringify({"action":"init"}));
    };
    ws.onmessage = function(msg) {
        console.log(msg)
        term.terminal.writeln(msg.data)
    };

}

