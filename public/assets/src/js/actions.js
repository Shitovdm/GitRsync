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
    let section = $(elem).data('section');
    let field = $(elem).data('field');
    let value = $(elem).val()
    if ($(elem).attr("type") === "checkbox") {
        value = $(elem).is(":checked")
    }
    if ($(elem).attr("type") === "text" && $(elem).data("field") === "MasterUser") {
        value = {
            "username": $("#commits_overriding__master_user_username").val(),
            "email": $("#commits_overriding__master_user_email").val()
        };
    }
    let ws = webSocketConnection("ws://localhost:8888/settings/save/");
    ws.onopen = function () {
        ws.send(JSON.stringify({"section": section, "field": field, "value": value}));
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

$('body').on('click', '.btn-add-committer-rule', function (e) {
    e.preventDefault();
    $("#add-committers-rules-row").remove()
    let html = $("#table-committers-rules").children("tbody").html()
    let itemsCount = $("#table-committers-rules").find("tr").length - 2;
    html += "<tr>" +
        "   <td class=\"btn-edit-committer-info\"" +
        "       data-toggle=\"modal\"" +
        "       data-target=\"#edit-committer-info-modal\"" +
        "       data-item=\"" + itemsCount + "\"" +
        "       data-type=\"old\"" +
        "       data-username=\"\"" +
        "       data-email=\"\">" +
        "       <b></b><br>" +
        "       <span></span>" +
        "   </td>" +
        "   <td style=\"cursor: default\">&rarr;</td>" +
        "   <td class=\"btn-edit-committer-info\"" +
        "       data-toggle=\"modal\"" +
        "       data-target=\"#edit-committer-info-modal\"" +
        "       data-item=\"" + itemsCount + "\"" +
        "       data-type=\"new\"" +
        "       data-username=\"\"" +
        "       data-email=\"\">" +
        "       <b></b><br>" +
        "       <span></span>" +
        "   </td>" +
        "   <td class=\"btn-remove-committer-rule\">&times;</td>" +
        "</tr>"

    html += "<tr id=\"add-committers-rules-row\">" +
        "<td colspan=\"4\" class=\"btn-add-committer-rule\">&#43;</td>" +
        "</tr>"
    $("#table-committers-rules").html("<tbody>" + html + "</tbody>")
});

$('body').on('click', '.btn-remove-committer-rule', function (e) {
    e.preventDefault();
    $(this).parent("tr").remove()
});

$('body').on('click', '.btn-edit-committer-info', function (e) {
    e.preventDefault();
    let form = $('#edit-committer-info-modal');
    form.find('input[name=type]').val($(this).data('type'));
    form.find('input[name=item]').val($(this).data('item'));
    form.find('input[name=username]').val($(this).data('username'));
    if ($(this).data('username') !== "") {
        form.find('input[name=username]').parent('div').addClass('is-filled');
    }
    form.find('input[name=email]').val($(this).data('email'));
    if ($(this).data('email') !== "") {
        form.find('input[name=email]').parent('div').addClass('is-filled');
    }
});

function SaveRules() {
    let section = "CommitsOverriding";
    let field = "CommittersRules";
    let value = ParseCommitersRules();
    let ws = webSocketConnection("ws://localhost:8888/settings/save/");
    ws.onopen = function () {
        ws.send(JSON.stringify({"section": section, "field": field, "value": value}));
    };
}

$('body').on('click', '.btn-save-committer-info', function (e) {
    e.preventDefault();
    let modal = $('#edit-committer-info-modal');
    let type = modal.find('input[name=type]').val();
    let item = modal.find('input[name=item]').val();
    let username = modal.find('input[name=username]').val();
    let email = modal.find('input[name=email]').val();

    let tableItem = $("#table-committers-rules").find("td[data-type=" + type + "]td[data-item=" + item + "]");
    tableItem.attr('data-username', username);
    tableItem.children("b").text(username)
    tableItem.attr('data-email', email);
    tableItem.children("span").text(email);

    SaveRules();
});

function ParseCommitersRules() {
    let table = $("#table-committers-rules").find("tr");
    let config = [];
    for (var i = 0; i < table.length - 1; i++) {
        let item = table[i];
        config.push({
            "old": {
                "username": $(item).find("td[data-type=old]").children("b").text(),
                "email": $(item).find("td[data-type=old]").children("span").text()
            },
            "new": {
                "username": $(item).find("td[data-type=new]").children("b").text(),
                "email": $(item).find("td[data-type=new]").children("span").text()
            }
        });
    }

    return config
}

$('body').on('change', '#commits_overriding__override_commits_with_one_author', function (e) {
    e.preventDefault();
    if ($(this).is(":checked")) {
        $(".master-user").removeClass("hidden");
        $(".committers-rules").addClass("hidden");
    } else {
        $(".master-user").addClass("hidden");
        $(".committers-rules").removeClass("hidden");
    }
});