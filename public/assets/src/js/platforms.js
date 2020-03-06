$('#platformForm').on('submit', function (e) {
    e.preventDefault();
    let form = $('#platformForm');
    let formData = getFormData(form);
    let ws = webSocketConnection("ws://localhost:8888/platforms/add/");
    ws.onopen = function () {
        ws.send(formData);
        reloadPageData();
    };
});

$('body').on('click', '.btn-remove-platform-modal', function (e) {
    $('.btn-remove-platform').attr('data-uuid', $(this).data('uuid'));
    $('.remove-platform-name').text($(this).data('name'));
});

$('body').on('click', '.btn-remove-platform', function (e) {
    e.preventDefault();
    let uuid = $(this).data('uuid');
    let formData = JSON.stringify({uuid: uuid});
    let ws = webSocketConnection("ws://localhost:8888/platforms/remove/");
    ws.onopen = function () {
        ws.send(formData);
        $('.btn-close-remove-platform').click();
        $('#platform_' + uuid).remove();
        reloadPageData();
    };
});

$('body').on('click', '.btn-edit-platform-modal', function (e) {
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


$('body').on('click', '.btn-edit-platform', function (e) {
    e.preventDefault();
    let form = $('#editPlatformForm');
    let formData = getFormData(form);
    let ws = webSocketConnection("ws://localhost:8888/platforms/edit/");
    ws.onopen = function () {
        ws.send(formData);
        $('.btn-close-edit-platform').click();
        reloadPageData();
    };
});

$(".platform-selector-act").change(function () {
    $("#" + $(this).attr('name') + "-prefix").text($(this.options[this.selectedIndex]).data('address'));
});
