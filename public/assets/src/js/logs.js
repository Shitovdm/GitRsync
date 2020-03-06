var term;



function getLogs() {
    if (typeof term === "undefined") {
        term = new ExecTerminal('terminal');
        term.UpdateTerminalFit();
    }
    term.terminal.writeln("hello world!")
}

