$(document).on("DOMContentLoaded", () => {
    // Fetch auth token. If it exists, fetch data.
    const token = Cookies.get("LOGGER-AUTH");
    if (token)
        fetchLogData(token);
    // If the user needs to authenticate, listen for sign in.
    $("button#signin").on("click", () => fetchLogData(null));
});
function fetchLogData(_pw) {
    let pw = _pw;
    if (!pw) {
        const input = $("input#password").val();
        if (typeof input !== "string" || input.length > 40) {
            displayError();
            return;
        }
        pw = input;
    }
    
    const conn = new WebSocket(`${window.location}api/read/`.replace("http", "ws"), pw);
    conn.onopen = () => {
        const log = $("#log");
        log.empty();
        $("#auth").hide();
        // Display initial data.
    };
    conn.onerror = () => {
        // Display error.
    };
    conn.onmessage = (ev) => {
        // Display message on frontend
    };
}

function displayData(data) {
    const log = $("#log");
    // Loop through data fetched from backend
    $.each(data, (_, item) => {
        const element = $("<p>").text(`${item.Origin} @ ${item.TimeWritten}: ${item.Message}`);
        // Highlight warn and error levels accordingly.
        switch (item.Severity) {
            case 1:
                element.css("background-color", "yellow");
                break;
            case 2:
                element.css("background-color", "red");
        }
        log.append(element);
    });
}
function displayError() {
    console.log();
}
