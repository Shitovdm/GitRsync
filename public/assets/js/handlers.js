var interval;
var term;
var connections = [];
var stream;

function getFormData($form){
    let unIndexedArray = $form.serializeArray();
    let indexedArray = {};

    $.map(unIndexedArray, function(n, i){
        indexedArray[n['name']] = n['value'];
    });

    return JSON.stringify(indexedArray);
}

function reloadPageData() {
    location.reload();
    return false;
}

$('#platformForm').on('submit', function(e)
{
    e.preventDefault();
    let form = $('#platformForm');
    let formData = getFormData(form);
    let ws = webSocketConnection("ws://localhost:8888/platforms/add/");
    ws.onopen = function()
    {
        ws.send(formData);
        reloadPageData();
    };
    ws.onmessage = function(msg) {
        term.terminal.writeln(msg.data)
    };
});

$('body').on('click', '.btn-remove-platform-modal', function (e)
{
    $('.btn-remove-platform').attr('data-uuid', $(this).data('uuid'));
    $('.remove-platform-name').text($(this).data('name'));
});

$('body').on('click', '.btn-remove-platform', function (e)
{
    e.preventDefault();
    let uuid = $(this).data('uuid');
    let formData = JSON.stringify({uuid: uuid});
    let ws = webSocketConnection("ws://localhost:8888/platforms/remove/");
    ws.onopen = function()
    {
        ws.send(formData);
        $('.btn-close-remove-platform').click();
        $('#platform_' + uuid).remove();
        reloadPageData();
    };
    ws.onmessage = function(msg) {
        term.terminal.writeln(msg.data)
    };
});

$('body').on('click', '.btn-edit-platform-modal', function (e)
{
    $('.btn-edit-platform').attr('data-uuid', $(this).data('uuid'));
    let form = $('#editPlatformForm');
    form.find('input[name=uuid]').val($(this).data('uuid'));
    form.find('input[name=name]').parent('div').addClass('is-filled');
    form.find('input[name=name]').val($(this).data('name'));
    form.find('input[name=address]').parent('div').addClass('is-filled');
    form.find('input[name=address]').val($(this).data('address'));
    form.find('input[name=username]').parent('div').addClass('is-filled');
    form.find('input[name=username]').val($(this).data('username'));
    form.find('input[name=password]').parent('div').addClass('is-filled');
    form.find('input[name=password]').val($(this).data('password'));
});

$('body').on('click', '.btn-edit-platform', function (e)
{
    e.preventDefault();
    let form = $('#editPlatformForm');
    let formData = getFormData(form);
    console.log(formData);
    let ws = webSocketConnection("ws://localhost:8888/platforms/edit/");
    ws.onopen = function()
    {
        ws.send(formData);
        $('.btn-close-edit-platform').click();
        reloadPageData();
    };
    ws.onmessage = function(msg) {
        term.terminal.writeln(msg.data)
    };
});

$('#repositoryForm').on('submit', function(e)
{
    e.preventDefault();
    let form = $('#repositoryForm');
    let formData = getFormData(form);
    let ws = webSocketConnection("ws://localhost:8888/repositories/add/");
    ws.onopen = function()
    {
        ws.send(formData);
        reloadPageData();
    };
    ws.onmessage = function(msg) {
        term.terminal.writeln(msg.data)
    };
});


















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


$(document).arrive(".selectpicker", function() {
    $(this).selectpicker();
});
