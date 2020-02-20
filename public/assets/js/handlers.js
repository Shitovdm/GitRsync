var interval;
var term;
var connections = [];
var stream;

function addProject() {
    let form = $('#projectForm');
    let formData = getFormData(form);
    if (typeof term === "undefined") {
        term = new ExecTerminal('terminal');
        term.UpdateTerminalFit();
    }

    let ws = webSocketConnection("ws://localhost:8888/projects/add/");
    ws.onopen = function()
    {
        ws.send(formData);
    };
    ws.onmessage = function(msg) {
        term.terminal.writeln(msg.data)
    };
}

function saveProjectSettings(settingsId) {
    let form = $('#'+settingsId);
    let formData = $(form).serialize();
    $.ajax({
        method: 'post',
        url: '/projects/settings/',
        data: formData,
        //contentType: "application/json; charset=utf-8",
        success: function (data) {
            load_settings(data.project, $('#'+data.element));
            $.notify({
                icon: "add_alert",
                message: "Настройки проекта успешно обновлены"
            }, {
                type: 'success',
                timer: 1000,
                placement: {
                    from: 'top',
                    align: 'center'
                }
            });
        },
    })
}

function load_settings(project, element) {
    let id = element.attr('id');
    $.ajax({
        method: 'get',
        url: '/projects/settings/',
        data: {project: project},
        success: function (data) {
            $('#collapse-'+id+' .card-body').html(data)
        },
    })
}

function startAll() {
    let launchC = [];
    let buildC = [];
    let formData = {};

    $('.launchC').each(function () {
        if ($(this).is(':checked')) {
            launchC.push($(this).data('project'));
        }
    });

    $('.buildC').each(function () {
        if ($(this).is(':checked')) {
            buildC.push($(this).data('project'));
        }
    });

    formData["launch"] = launchC;
    formData["build"] = buildC;
    if (typeof term === "undefined") {
        term = new ExecTerminal('terminal');
        term.UpdateTerminalFit();
    }

    let ws = webSocketConnection("ws://localhost:8888/launcher/");
    ws.onopen = function()
    {
        ws.send(JSON.stringify(formData));
    };
    ws.onmessage = function(msg) {
        term.terminal.writeln(msg.data)
    };
}

function stopAll() {
    let launchC = [];
    let formData = {};

    $('.launchC').each(function () {
        if ($(this).is(':checked')) {
            launchC.push($(this).data('project'));
        }
    });
    formData["stop"] = launchC;

    if (typeof term === "undefined") {
        term = new ExecTerminal('terminal');
        term.UpdateTerminalFit();
    }

    let ws = webSocketConnection("ws://localhost:8888/launcher/stopall/");
    ws.onopen = function()
    {
        ws.send(JSON.stringify(formData));
    };
    ws.onmessage = function(msg) {
        term.terminal.write(msg.data.replace(/\r/g, '\n\r'))
    };
}

function getServices() {
    $.ajax({
        method: 'post',
        url: '/docker/getServices/',
        dataType: "json",
        success: function (data) {
            $('.servicelist').html('');
            $.each(data, function (project, data) {
                let statuses = '';
                $.each(data, function (service, serviceData) {
                    let status = serviceData.status;
                    let term = serviceData.terminal;
                    statuses += '<span class="service-state badge badge-'+status+'" title="'+status+'"> ';
                    if (status === 'running') {
                        statuses += '<span onclick="stopService(\''+service+'\', \''+project+'\')" class="text-danger action" title="Остановить сервис"><i class="fas fa-stop"></i></span> ';
                        if (term !== "") {
                            statuses += '<span onclick="launchCmd(\''+service+'\', \''+project+'\')" class="text-primary action" title="Запустить shell"><i class="fas fa-terminal"></i></span> ';
                        }
                    }
                    statuses += '<span onclick="logsService(\''+service+'\', \''+project+'\')" class="text-primary action" title="Логи сервиса"><i class="fas fa-list-alt"></i></span> ';
                    statuses +=  service + '</span>';
                });
                $('#serviceList-'+project).html(statuses)
            });
        },
    })
}

function launchCmd(service, project) {
    $.ajax({
        method: 'get',
        url: '/docker/launch-cmd/',
        data: {service: service, project: project},
        success: function (data) {

        },
    })
}

function gitPull(project) {
    let formData = {};
    formData["project"] = project;
    if (typeof term === "undefined") {
        term = new ExecTerminal('terminal');
        term.UpdateTerminalFit();
    }

    let ws = webSocketConnection("ws://localhost:8888/launcher/gitpull/");
    ws.onopen = function()
    {
        ws.send(JSON.stringify(formData));
    };
    ws.onmessage = function(msg) {
        term.terminal.writeln(msg.data)
    };
}

function reloadConfig(project) {
    $.ajax({
        method: 'post',
        url: '/projects/reloadconfig/',
        data: {project: project},
        success: function (data) {
            load_settings(project, $('#'+data.element));
            $.notify({
                icon: "add_alert",
                message: "Настройки проекта успешно обновлены"
            }, {
                type: 'success',
                timer: 1000,
                placement: {
                    from: 'top',
                    align: 'center'
                }
            });
        },
    });
}

function stopService(serviceName, project) {
    let formData = {};
    formData["project"] = project;
    formData["service"] = serviceName;
    if (typeof term === "undefined") {
        term = new ExecTerminal('terminal');
        term.UpdateTerminalFit();
    }

    let ws = webSocketConnection("ws://localhost:8888/launcher/stopservice/");
    ws.onopen = function()
    {
        ws.send(JSON.stringify(formData));
    };
    ws.onmessage = function(msg) {
        term.terminal.writeln(msg.data);
    };
}

function logsService(serviceName) {
    if (typeof term === "undefined") {
        term = new ExecTerminal('terminal');
        term.UpdateTerminalFit();
    }
    term.terminal.clear();
    if (stream !== undefined) {
        stream.close();
    }
    var url = "/docker/getLogs/" + serviceName;
    stream = new EventSource(url);
    stream.addEventListener("end", function (e) {
        console.log("end");
        stream.close();
    });
    stream.addEventListener("message", function (msg) {
        console.log(msg.data);
        term.terminal.writeln(msg.data);
        //$(document).find('.log').append('<span class="message">$ ' + e.data + '</span><br/>')
    });
}

function webSocketConnection(url) {
    clearConnections();
    let ws = new WebSocket(url);
    connections.push(ws);

    ws.onclose = function() {
        let index = connections.indexOf(ws);
        if(index !== -1) {
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

if ($('#projectForm').length > 0) {
    $('#projectForm').on('submit', function(e) {
        e.preventDefault();
        addProject();
    });
}


if ($('#accordion').length > 0) {
    setTimeout(function () {
        $('#accordion .card').each(function() {
            load_settings($(this).data('project'), $(this));
        })
    }, 1000);
    $('#startAll').click(function () {
        startAll();
    });
    $('#stopAll').click(function () {
        stopAll();
    });
    $('.git-puller').click(function () {
        gitPull($(this).data('project'));
    });
    $('.config-reloader').click(function () {
        reloadConfig($(this).data('project'));
    });
    setInterval(function () {
        getServices();
    }, 5000);
}

function getFormData($form){
    let unIndexedArray = $form.serializeArray();
    let indexedArray = {};

    $.map(unIndexedArray, function(n, i){
        indexedArray[n['name']] = n['value'];
    });

    return JSON.stringify(indexedArray);
}

$(document).arrive(".selectpicker", function() {
    $(this).selectpicker();
});

$(document).arrive(".projectSettingsForm", function () {
    let settingsId = $(this).attr('id');
    $(this).on('submit', function(e) {
        e.preventDefault();
        saveProjectSettings(settingsId);
    });
});