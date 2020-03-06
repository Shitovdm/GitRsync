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
