function updateServerList() {
  var serverlist;
  $.ajax({
    url: "/server-list",
    type: "GET",
    success: (data) => {
      console.log(data);
      document.getElementById("servers-online").innerHTML = data;
    }    
  })
}
