var connections = [];
var stream;

function webSocketConnection(url) {
    clearConnections();
    let ws = new WebSocket(url);
    connections.push(ws);

    ws.onclose = function () {
        let index = connections.indexOf(ws);
        if (index !== -1) {
            connections.splice(index, 1);
        }
    };

    return ws;
}

function clearConnections() {
    if (connections.length > 0) {
        connections.forEach(function (connection) {
            try {
                connection.close();
            } catch (e) {
                console.log(e);
            }
        })
    }
}