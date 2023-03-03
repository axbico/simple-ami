
function addUserState(message, uniqueid = "", classStatus = "", sign = ".") {
    $("#users-inject").prepend(
        $("#template-user-state").html().
            replace("%userid%", message).replace("%message%", message).
            replace("%id%", uniqueid).replace("%sign%", sign).replace("%status%", classStatus)
    );
}

$(document).ready(function() {
    var ws = new WebSocket("wss://" + window.location.host + "/api/users/connect");

    ws.onopen = function() {
        $("#users-connection-status").toggleClass("establish", true);
        $("#users-connection-status").toggleClass("lost", false);
        
        $("#users-inject").html("")
    };

    ws.onmessage = function(event) {
        var pair = String(event.data).split("~");
        
        switch(pair[0]) {
            case "USERS_CONTACT_ADD": case "USERS_CONTACT_REMOVE": case "USERS_CONTACT_UNKNOWN":
                $("#userid-" + pair[1]).remove();
            default:
        }
        addUserState(
            pair[1], 
            pair.length > 2 ? pair[2] : "", 
            pair.length > 2 ? "establish" : "", 
            pair.length > 2 ? "<" : "."
        );
    };

    ws.onclose = function() {
        $("#users-connection-status").toggleClass("establish", false);
        $("#users-connection-status").toggleClass("lost", true);
    };

});