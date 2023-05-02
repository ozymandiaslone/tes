function updateQueryStatus() {
  var uuid = document.getElementById("uuid").getAttribute("uuid-data");
  $.ajax({
    type: "GET",
    url: "/response/" + uuid + "/sync",
    success: function(data) {
      var response = $.parseJSON(data);
      if (response.hasOwnProperty("error")) {
        console.log(response.error);
      } else {
        document.getElementById("response").innerHTML = response.response;
      }
    },
    error: function(data) {
      console.log(data);
    }
  });
}

// Call updateQueryStatus every 3 seconds
setInterval(updateQueryStatus, 3000);

