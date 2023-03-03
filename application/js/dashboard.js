
var allUsersCount = 0;
var activeUsersCount = 0;
var activeCallsCount = 0;

const RECENT = "recent";
const CALL = "call";

function addActivity(activity, message, uniqueid = "", classStatus = "") {
    $("#log-" + activity + "-activities").prepend(
        $("#template-log-recent-activities").html().
            replace("%timestamp%", (new Date).toLocaleTimeString([], {hour12:false})).
            replace("%message%", message).replace("%id%", uniqueid).replace("%status%", classStatus)
    );

    if($("#log-" + activity + "-activities").find('.activity').length > 12) {
        $("#log-" + activity + "-activities").find('.activity:last').remove();
    }
}

$(document).ready(function() {
    var ws = new WebSocket("wss://" + window.location.host + "/api/dashboard/connect");

    ws.onopen = function() {
        addActivity(RECENT, "", "", "establish");
        addActivity(CALL, "", "", "establish");
    };

    ws.onmessage = function(event) {
        var pair = String(event.data).split("~");
        switch(pair[0]) {
            case "ALL_USERS":
                allUsersCount = Number(pair[1]);
                $("#all-users").html(allUsersCount);
                break;
            case "ACTIVE_USERS":
                activeUsersCount = Number(pair[1]);
                $("#active-users").html(activeUsersCount);
                break;
            case "ACTIVE_USERS_ALTER":
                activeUsersCount += Number(pair[1]);
                $("#active-users").html(activeUsersCount);
                break;
            case "ACTIVE_CALLS":
                activeCallsCount = Number(pair[1]);
                $("#active-calls").html(activeCallsCount);
                break;
            case "ACTIVE_CALLS_ALTER":
                activeCallsCount += Number(pair[1]);
                $("#active-calls").html(activeCallsCount)
                break;
            case "RECENT_ACTIVITIES":
                addActivity(RECENT, pair[1], pair.length > 2 ? pair[2] : "");
                break;
            case "CALL_ACTIVITIES":
                addActivity(CALL, pair[1], pair.length > 2 ? pair[2] : "");
                break;
            default:
        }
    };

    ws.onclose = function() {
        addActivity(RECENT, "", "", "lost");
        addActivity(CALL, "", "", "lost");
    };
    
});