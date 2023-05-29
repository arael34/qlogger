$(document).on("DOMContentLoaded", () => {
    // Fetch auth token. If it exists, fetch data.
    const token = Cookies.get("LOGGER-AUTH");
    if (token)
        fetchLogData(token, null);

    $("#log-container").hide();
    // If the user needs to authenticate, listen for sign in.
    $("button#signin").on("click", () => fetchLogData(null, null));
    $("input#filter").on("change", (ev) => {
        fetchLogData(token, $(ev.target).val());
    });
});

function fetchLogData(_pw, filter) {
    let pw = _pw;
    if (!pw) {
        const input = $("input#password").val();
        if (typeof input !== "string" || input.length > 40) {
            displayError("invalid password");
            return;
        }
        pw = input;
    }
    
    const conn = new WebSocket(`${window.location}api/read/`.replace("http", "ws"), pw);
    conn.onopen = () => {
        $("#log-container").show();
        $("#log").empty();
        $("#auth").hide();
        Cookies.set("LOGGER-AUTH", pw, { expires: 1 });
    };
    conn.onerror = () => {
        displayError("couldn't connect to websocket, or not authorized.");
    };
    conn.onmessage = (ev) => {
        const parsedData = JSON.parse(ev.data);
        if (!filter || parsedData.Origin.includes(filter))
            displayData(parsedData);
    };
}

function displayData(item) {
    const log = $("#log");

    const element = $("<p>").text(`${item.Origin} @ ${item.TimeWritten}: ${item.Message}`);
    // Highlight warn and error levels accordingly.
    switch (item.Severity) {
        case 1:
            element.css("background-color", "yellow");
            break;
        case 2:
            element.css("background-color", "red");
    }
    log.prepend(element);
}

function displayError(message) { console.log(message); }
