$(document).ready(() => {
    $.ajax({
      url: "asdfsdaklf",
      method: "GET",
      dataType: "json",
      success: (response) => {
        displayData(response);
      },
      error: (err) => {
        displayData(fakeData);
        console.log(`Error fetching data: ${err.message}`);
      }
    });

    function displayData(data) {
      const log = $('#log');

      log.empty();

      $.each(data, (_, item) => {
        const element = $("<p>").text(item.message);
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
