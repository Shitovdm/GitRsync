$('body').on('click', '.btn-pull-source-repository', function (e) {
    ToggleAjaxPreloader()
    e.preventDefault();
    let uuid = $(this).data('uuid');
    let ws = webSocketConnection("ws://localhost:8888/actions/pull/");
    let repoObj = this;
    SetRepositoryStatus(repoObj, "pending_pull");
    ws.onopen = function () {
        console.log("Cloning or pulling source repository...");
        ws.send(JSON.stringify({"uuid": uuid}));
    };
    ws.onmessage = function (msg) {
        ToggleAjaxPreloader()
        let body = JSON.parse(msg.data);
        showNotification(body["status"], body["message"]);
        switch (body["status"]) {
            case "success":
                SetRepositoryStatus(repoObj, "pulled");
                break;
            case "error":
                SetRepositoryStatus(repoObj, "pull_failed");
                break;
        }
    };
});

$('body').on('click', '.btn-push-destination-repository', function (e) {
    ToggleAjaxPreloader()
    e.preventDefault();
    let uuid = $(this).data('uuid');
    let ws = webSocketConnection("ws://localhost:8888/actions/push/");
    let repoObj = this;
    SetRepositoryStatus(repoObj, "pending_push");
    ws.onopen = function () {
        console.log("Pushing to destination repository...");
        ws.send(JSON.stringify({"uuid": uuid}));
    };
    ws.onmessage = function (msg) {
        ToggleAjaxPreloader()
        let body = JSON.parse(msg.data);
        showNotification(body["status"], body["message"]);
        switch (body["status"]) {
            case "success":
                SetRepositoryStatus(repoObj, "synced");
                break;
            case "error":
                SetRepositoryStatus(repoObj, "push_failed");
                break;
        }
    };
});

$('body').on('click', '.btn-block-repository', function (e) {
    ToggleAjaxPreloader()
    e.preventDefault();
    let uuid = $(this).data('uuid');
    let formData = JSON.stringify({uuid: uuid});
    let request = new XMLHttpRequest();
    request.open("POST", "actions/block", true);
    request.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    request.addEventListener("readystatechange", () => {
        if (request.readyState === 4 && request.status === 200) {
            console.log(request.responseText);
        }
    });

    request.send(formData);
});