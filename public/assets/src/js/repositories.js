$('#repositoryForm').on('submit', function (e) {
    e.preventDefault();
    let form = $('#repositoryForm');
    let formData = getFormData(form);
    let ws = webSocketConnection("ws://localhost:8888/repositories/add/");
    ws.onopen = function () {
        ws.send(formData);
        reloadPageData();
    };
});

$('body').on('click', '.btn-remove-repository-modal', function (e) {
    $('.btn-remove-repository').attr('data-uuid', $(this).data('uuid'));
    $('.remove-repository-name').text($(this).data('name'));
});

$('body').on('click', '.btn-remove-repository', function (e) {
    e.preventDefault();
    let uuid = $(this).data('uuid');
    let formData = JSON.stringify({uuid: uuid});
    let ws = webSocketConnection("ws://localhost:8888/repositories/remove/");
    ws.onopen = function () {
        ws.send(formData);
        $('.btn-close-remove-repository').click();
        $('#repository_' + uuid).remove();
        reloadPageData();
    };
});

$('body').on('click', '.btn-edit-repository-modal', function (e) {
    $('.btn-edit-repository').attr('data-uuid', $(this).data('uuid'));
    let form = $('#editRepositoryForm');
    form.find('input[name=uuid]').val($(this).data('uuid'));
    form.find('input[name=name]').parent('div').addClass('is-filled');
    form.find('input[name=name]').val($(this).data('name'));
    form.find('select[name=spu]').val($(this).data('spu'));
    form.find('input[name=spp]').parent('div').addClass('is-filled');
    form.find('input[name=spp]').val($(this).data('spp'));
    form.find('select[name=dpu]').val($(this).data('dpu'));
    form.find('input[name=dpp]').parent('div').addClass('is-filled');
    form.find('input[name=dpp]').val($(this).data('dpp'));
    $('.selectpicker').selectpicker('refresh')
});


$('body').on('click', '.btn-edit-repository', function (e) {
    e.preventDefault();
    let form = $('#editRepositoryForm');
    let formData = getFormData(form);
    let ws = webSocketConnection("ws://localhost:8888/repositories/edit/");
    ws.onopen = function () {
        ws.send(formData);
        $('.btn-close-edit-repository').click();
        reloadPageData();
    };
});