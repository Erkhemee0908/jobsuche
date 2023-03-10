// Define the URL of the JSON file
var url = "../Backend/jobList.json";

// Use the fetch API to load the JSON data
fetch(url)
  .then(response => response.json())
  .then(data => {
    // Create an HTML string to display the data
    var html = "<ul>";
    data.forEach(item => {
      html += "<li>" + item.titel + ": " + item.description + "</li>";
    });
    html += "</ul>";

    // Update the data-container element with the HTML
    document.getElementById("data-container").innerHTML = html;
  })
  .catch(error => {
    console.error(error);
  });
