function getFormData($form) {
    let unIndexedArray = $form.serializeArray();
    let indexedArray = {};

    $.map(unIndexedArray, function (n, i) {
        indexedArray[n['name']] = n['value'];
    });

    return JSON.stringify(indexedArray);
}

function reloadPageData() {
    location.reload();
    return false;
}

$(document).arrive(".selectpicker", function () {
    $(this).selectpicker();
});

function showNotification(type, message){
    let icon = "";
    switch (type) {
        case "success":
            icon = "add_alert";
            break;
        case "error":
            icon = "error";
            break;
        case "warning":
            icon = "warning";
            break;
        case "info":
            icon = "info";
            break;
    }
    $.notify({
        icon: icon,
        message: message

    },{
        type: type,
        timer: 3000,
        placement: {
            from: 'top',
            align: 'right'
        }
    });
}