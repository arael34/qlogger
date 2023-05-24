$(document).ready(() => {
    // Fetch log data
    $.ajax({
      url: `${window.location}/api/read`,
      method: "GET",
      dataType: "json",
      success: (response) => {
        displayData(response);
      },
      error: (err) => {
        console.log(`Error fetching data: ${err.message}`);
      }
    });

    function displayData(data) {
      const log = $('#log');
      log.empty();

      // Loop through data fetched from backend
      $.each(data, (_, item) => {
        const element = $("<p>").text(item.message);

        // Highlight warn and error levels accordingly.
        switch (item.level) {
            case 1:
                element.css("background-color", "yellow");
                break;
            case 2:
                element.css("background-color", "red");
        }
        log.append(element);
      });
    }
  });
