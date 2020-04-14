$('body').on('click', '.btn-pull-source-repository', function (e) {
    e.preventDefault();
    let uuid = $(this).data('uuid');
    let ws = webSocketConnection("ws://localhost:8888/actions/pull/");
    ws.onopen = function () {
        console.log("Cloning or pulling source repository...");
        ws.send(JSON.stringify({"uuid": uuid}));
    };
    ws.onmessage = function (msg) {
        let body = JSON.parse(msg.data);
        showNotification(body["status"], body["message"])
    };
});

$('body').on('click', '.btn-push-destination-repository', function (e) {
    e.preventDefault();
    let uuid = $(this).data('uuid');
    let ws = webSocketConnection("ws://localhost:8888/actions/push/");
    ws.onopen = function () {
        console.log("Pushing to destination repository...");
        ws.send(JSON.stringify({"uuid": uuid}));
    };
    ws.onmessage = function (msg) {
        let body = JSON.parse(msg.data);
        showNotification(body["status"], body["message"])
    };
});

$('body').on('click', '.btn-block-repository', function (e) {
    e.preventDefault();
    let uuid = $(this).data('uuid');
    let formData = JSON.stringify({uuid: uuid});
    let request = new XMLHttpRequest();
    request.open("POST", "actions/block", true);
    request.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    request.addEventListener("readystatechange", () => {
        if(request.readyState === 4 && request.status === 200) {
            console.log(request.responseText);
        }
    });

    request.send(formData);
});