$('body').on('click', '.btn-pull-source-repository', function (e) {
    e.preventDefault();
    ToggleAjaxPreloader();
    new PullSourceRepository($(this).data('uuid'), this);
});

$('body').on('click', '.btn-push-destination-repository', function (e) {
    e.preventDefault();
    ToggleAjaxPreloader();
    new PushDestinationRepository($(this).data('uuid'), this);
});

$('body').on('click', '.btn-sync-repositories', function (e) {
    e.preventDefault();
    ToggleAjaxPreloader();
    let uuid = $(this).data('uuid')
    let repoObj = this
    let ws = webSocketConnection("ws://localhost:8888/actions/pull/");
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
                setTimeout(function () {
                    ToggleAjaxPreloader();
                    new PushDestinationRepository(uuid, repoObj);
                }, 1000)
                break;
            case "error":
                SetRepositoryStatus(repoObj, "pull_failed");
                showNotification("error", "Sync is broken! Unable to sync repositories!");
                break;
        }
    };
});

$('body').on('click', '.btn-clear-repositories', function (e) {
    e.preventDefault();
    ToggleAjaxPreloader();
    new ClearRepositoryRuntimeData($(this).data('uuid'), this);
});

function ClearRepositoryRuntimeData(uuid, repoObj) {
    let ws = webSocketConnection("ws://localhost:8888/actions/clear/");
    SetRepositoryStatus(repoObj, "pending_clear");
    ws.onopen = function () {
        console.log("Cleaning repository runtime data...");
        ws.send(JSON.stringify({"uuid": uuid}));
    };
    ws.onmessage = function (msg) {
        ToggleAjaxPreloader()
        let body = JSON.parse(msg.data);
        showNotification(body["status"], body["message"]);
        switch (body["status"]) {
            case "success":
                SetRepositoryStatus(repoObj, "cleared");
                break;
            case "error":
                SetRepositoryStatus(repoObj, "clear_failed");
                break;
        }
    };
}

function PullSourceRepository(uuid, repoObj) {
    let ws = webSocketConnection("ws://localhost:8888/actions/pull/");
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
}

function PushDestinationRepository(uuid, repoObj) {
    let ws = webSocketConnection("ws://localhost:8888/actions/push/");
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
}

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