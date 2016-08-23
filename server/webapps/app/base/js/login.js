$(document).ready(function () {
    Metronic.init(); // init metronic core componets
    Layout.init(); // init layout
    QuickSidebar.init(); // init quick sidebar
    Demo.init(); // init demo features

    $(document).ajaxError(function (event, jqxhr, settings, thrownError) {
        if (thrownError == 'Unauthorized') {
            toastr.error("Please check your username and password")
        }
    });
    $("#login-form").validate({
        rules: {
            password: {
                required: true
            },
            username: {
                required: true,
                email: true
            }
        }
    });

    $("#login-btn").on('click', function (e) {
        e.preventDefault();
        if ($("#login-form").valid()) {
            userRealm = $("#username").val().split('@');
            var payload = {username: userRealm[0], password: $("#password").val(), tenantdomain: userRealm[1]};
            $.post(loginAPIUrl, JSON.stringify(payload), function (result) {
                Cookies.set("username", userRealm[0]);
                Cookies.set("tenantid", result.tenantid);
                Cookies.set("tenantdomain", userRealm[1]);
                Cookies.set("jwt", result.token);
                window.location.href = document.referrer
            });
        }
    })
});