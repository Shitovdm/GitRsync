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

$('body').on('click', '.btn-info-repositories', function (e) {
    e.preventDefault();
    let uuid = $(this).data('uuid')
    let info = $("#repository-info-" + uuid)
    if (info.css("display") === "contents") {
        info.css("display", "none");
        info.html("");
        return
    }
    ToggleAjaxPreloader();
    let ws = webSocketConnection("ws://localhost:8888/actions/info/");
    ws.onopen = function () {
        ws.send(JSON.stringify({"uuid": uuid}));
    };
    ws.onmessage = function (msg) {
        ToggleAjaxPreloader()
        let body = JSON.parse(msg.data);
        showNotification(body["status"], body["message"]);
        switch (body["status"]) {
            case "success":
                let infoHtml = '<td colspan="5">\n' +
                    '<table class="repo-commit-container">\n'
                body["data"].forEach(function (commit) {
                    infoHtml += '<tr>\n' +
                        '<td>' + commit["Hash"].substr(0, 10) + "..." + commit["Hash"].substr(commit["Hash"].length - 5) + '</td>\n' +
                        '<td>' + commit["ParentHash"].substr(0, 10) + "..." + commit["ParentHash"].substr(commit["ParentHash"].length - 5) + '</td>\n' +
                        '<td>' + commit["Author"] + ' &lt;' + commit["AuthorEmail"] + '&gt;</td>\n' +
                        '<td>' + commit["Timestamp"] + '</td>\n' +
                        '<td>' + commit["Subject"] + '</td>\n' +
                        '</tr>\n'
                })
                infoHtml += '</table>\n' +
                    '</td>\n'
                info.html(infoHtml);
                info.css("display", "contents");
                break;
        }
    };
});

$('body').on('click', '.btn-activate-repositories', function (e) {
    e.preventDefault();
    ToggleAjaxPreloader();
    let uuid = $(this).data('uuid')
    let repoLine = $("#repository-" + uuid);
    let ws = webSocketConnection("ws://localhost:8888/actions/activate/");
    ws.onopen = function () {
        ws.send(JSON.stringify({"uuid": uuid}));
    };
    ws.onmessage = function (msg) {
        ToggleAjaxPreloader()
        let body = JSON.parse(msg.data);
        showNotification(body["status"], body["message"]);
        switch (body["status"]) {
            case "success":
                var repoLineHtml = repoLine.html();
                repoLine.remove();
                $("#active-repositories").append("<tr id='repository-" + uuid + "'>" + repoLineHtml + "</tr>");
                break;
        }
    };
});

$('body').on('click', '.btn-block-repositories', function (e) {
    e.preventDefault();
    ToggleAjaxPreloader();
    let uuid = $(this).data('uuid');
    let repoLine = $("#repository-" + uuid);
    let ws = webSocketConnection("ws://localhost:8888/actions/block/");
    ws.onopen = function () {
        ws.send(JSON.stringify({"uuid": uuid}));
    };
    ws.onmessage = function (msg) {
        ToggleAjaxPreloader()
        let body = JSON.parse(msg.data);
        showNotification(body["status"], body["message"]);
        switch (body["status"]) {
            case "success":
                var repoLineHtml = repoLine.html();
                repoLine.remove();
                $("#blocked-repositories").append("<tr id='repository-" + uuid + "'>" + repoLineHtml + "</tr>");
                break;
        }
    };
});

function SaveConfigField(elem) {
    ToggleAjaxPreloader();
    let section = $(elem).data('section');
    let field = $(elem).data('field');
    let value = $(elem).val()
    if($(elem).attr("type") === "checkbox"){
        value = $(elem).is(":checked")
    }
    console.log(section, field, value)
    let ws = webSocketConnection("ws://localhost:8888/settings/save/");
    ws.onopen = function () {
        ws.send(JSON.stringify({"section": section, "field": field, "value": value}));
        ToggleAjaxPreloader();
    };
}

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
