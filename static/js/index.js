$(document).on("DOMContentLoaded", () => {
    fetchLogData(null);

    $("#log-container").hide();
    // If the user needs to authenticate, listen for sign in.
    $("button#signin").on("click", () => fetchLogData(null));
    $("input#filter").on("change", (ev) => {
        fetchLogData($(ev.target).val());
    });
});

function fetchLogData(filter) {
    const pw = Cookies.get("LOGGER-AUTH") || $("input#password").val();
    // Filter out potential bad passwords.
    if (!pw || typeof pw !== "string" || pw.length > 40) {
        displayError("invalid password");
        return;
    }
   
    // This is horrible, but it works.
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
    const element = $("<p>").text(
        `${item.Origin}/${item.Category} @ ${item.TimeWritten}: ${item.Message}`
    );
    // Highlight warn and error levels accordingly.
    switch (item.Severity) {
        case 1:
            element.css("background-color", "yellow");
            break;
        case 2:
            element.css("background-color", "orange");
            break;
        case 2:
            element.css("background-color", "red");
    }
    log.prepend(element);
    // Remove the last element if there are more than 50
    if (log.children().length > 50)
        log.children().last().remove();
}

function displayError(message) { console.error(message); }
