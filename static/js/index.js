$(document).ready(() => {
  // Fetch auth token. If it exists, fetch data.
  const token = Cookies.get("LOGGER-AUTH");
  if (token) fetchLogData(token);

  // If the user needs to authenticate, listen for sign in.
  $("button#signin").on("click", () => fetchLogData(null));
});

function fetchLogData(_pw) {
  let pw = _pw;
  if (!pw) pw = $("input#password").val();
  if (pw.length > 50) return;

  // Fetch log data
  $.ajax({
    url: `${window.location}/api/read`,
    method: "GET",
    headers: {"Authorization": pw},
    dataType: "json",
    success: (response) => {
      // Display log data and set auth token.
      displayData(response);
      Cookies.set("LOGGER-AUTH", pw, { expires: 1 });
    },
    error: (err) => {
      console.log(`Error fetching data: ${err.message}`);
    }
  });
}

function displayData(data) {
  const log = $('#log');
  log.empty();
  $("#auth").hide();

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
