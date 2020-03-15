$('body').on('click', '.btn-pull-source-repository', function (e) {

    e.preventDefault();
    let uuid = $(this).data('uuid');
    let ws = webSocketConnection("ws://localhost:8888/actions/pull/");
    ws.onopen = function () {
        console.log("Cloning...")
        ws.send(JSON.stringify({"uuid": uuid}));
    };
    ws.onmessage = function (msg) {
        console.log()
        if(msg.data === "done") {
            ws.close()
            console.log("cloning done!!")
        }
    };
});

$('body').on('click', '.btn-block-repository', function (e) {
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