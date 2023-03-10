// Define the URL of the JSON file
var url = "http://127.0.0.1:5500/Backend/jobList.json";

// Use the fetch API to load the JSON data
fetch(url)
  .then(response => response.json())
  .then(data => {
    // Create an HTML string to display the data
    var html = "<ol>";
    data.forEach(item => {
      html += "<li><b>" + item.titel + ": </b><br><br>" + item.description + "</li><br><br>";
    });
    html += "</ol>";

    // Update the data-container element with the HTML
    document.getElementById("data-container").innerHTML = html;
  })
  .catch(error => {
    console.error(error);
  });

  
  