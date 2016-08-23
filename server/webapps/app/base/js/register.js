$(document).ready(function () {
    $(document).ajaxError(function (event, jqxhr, settings, thrownError) {
        if (thrownError == 'Unauthorized') {
            toastr.error("Please check your username and password")
        }
    });

    $("#btn-register").on('click', function (e) {
        e.preventDefault()
        var payload = {
            username: $("#username").val().split('@')[0],
            password: $("#password").val(),
            email: $("#email").val(),
            tenantid: 1,
            status: 'pending',
            roles: [],
            permissions: []
        };
        if ($("#registration-form").valid()) {
            $.post("/dashboard/users", JSON.stringify(payload), function (result) {
                Cookies.set("username", $("#username").val());
                toastr.info("Account activation is pending for approval. Please contact the System Administrator.")
            })
        }
    })
});