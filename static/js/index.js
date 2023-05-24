function fetchLogData() {
  const pw = $("input#password").val();

  // Fetch log data
  $.ajax({
    url: `${window.location}/api/read`,
    method: "GET",
    headers: {"Authorization": pw},
    dataType: "json",
    success: (response) => {
      console.log(response);
      displayData(response);
    },
    error: (err) => {
      console.log(`Error fetching data: ${err.message}`);
    }
  });

  function displayData(data) {
    const log = $('#log');
    log.empty();
    $("#auth").hide();

    // Loop through data fetched from backend
    $.each(data, (_, item) => {
      const element = $("<p>").text(item.Message);

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
}

$("button#signin").on("click", fetchLogData);
