$('body').on('click', '.btn-open-raw-config', function (e) {
    e.preventDefault();
    let ws = webSocketConnection("ws://localhost:8888/actions/open/dir/");
    ws.onopen = function () {
        ws.send(JSON.stringify({"path": 'projects\\'}));
    };
});

$('body').on('click', '.btn-explore-config-dir', function (e) {
    e.preventDefault();
    let ws = webSocketConnection("ws://localhost:8888/actions/open/dir/");
    ws.onopen = function () {
        ws.send(JSON.stringify({"path": 'projects\\'}));
    };
});