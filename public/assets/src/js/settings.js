$('body').on('click', '.btn-open-raw-config', function (e) {
    e.preventDefault();
    let ws = webSocketConnection("ws://localhost:8888/settings/open-raw-config/");
    ws.onopen = function () {
        ws.send();
    };
});

$('body').on('click', '.btn-explore-config-dir', function (e) {
    e.preventDefault();
    let ws = webSocketConnection("ws://localhost:8888/settings/explore-config-dir/");
    ws.onopen = function () {
        ws.send();
    };
});