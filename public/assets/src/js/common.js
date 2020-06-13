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

function showNotification(type, message) {
    let icon = "";
    let notifyType = "info";
    switch (type) {
        case "success":
            notifyType = "success";
            icon = "add_alert";
            break;
        case "error":
            notifyType = "danger";
            icon = "error";
            break;
        case "warning":
            notifyType = "warning";
            icon = "warning";
            break;
        case "info":
            notifyType = "info";
            icon = "info";
            break;
    }
    $.notify({
        icon: icon,
        message: message

    }, {
        type: notifyType,
        timer: 4000,
        placement: {
            from: 'top',
            align: 'right'
        }
    });
}

function ToggleAjaxPreloader() {
    if ($("#preloader").css("display") === "none") {
        $("#preloader").show(100)
    } else {
        $("#preloader").hide(100)
    }
}
