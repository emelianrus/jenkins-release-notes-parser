console.log("loaded servers.js");


$(document).ready(function () {
  // add a click event listener to the plugins list (the parent element)
  $(".delete-server-btn").on("click", function () {
    // get the ID of the item to delete
    var index = $(this).data("server-name");
    // create a payload object with the ID of the item
    var payload = {
      jenkinsName: index,
    };
    console.log(payload);
    var liBlock = $(this).closest("li");
    // send an AJAX request to the server
    $.ajax({
      url: "/delete-jenkins-server", // replace this with the actual URL of the delete endpoint
      type: "POST",
      contentType: "application/json",
      data: JSON.stringify(payload),
      success: function (response) {
        liBlock.remove();
      },
      error: function (jqXHR, textStatus, errorThrown) {
        // handle any errors that occur during the request here
        // console.error(textStatus, errorThrown);
      },
    });
  });
  // add a click event listener to the plugins list (the parent element)
  $(".delete-btn").on("click", function () {
    // get the ID of the item to delete
    var index = $(this).data("index");
    var res = index.split(":");

    // create a payload object with the ID of the item
    var payload = {
      jenkinsName: res[0],
      pluginName: res[1],
    };
    console.log(payload);
    var liBlock = $(this).closest("li");
    // send an AJAX request to the server
    $.ajax({
      url: "/delete-plugin", // replace this with the actual URL of the delete endpoint
      type: "POST",
      contentType: "application/json",
      data: JSON.stringify(payload),
      success: function (response) {
        liBlock.remove();
      },
      error: function (jqXHR, textStatus, errorThrown) {
        // handle any errors that occur during the request here
        // console.error(textStatus, errorThrown);
      },
    });
  });
});

$(document).on("click", ".add-new-plugin-btn", function () {
  var serverName = $(this).data("server-name");
  console.log(serverName);
  $("#addPluginModal").find("#jenkins-server-name").val(serverName);
});

$(document).ready(function () {
  // add new field
  $("#addPluginModal").on("click", ".add-field-btn", function () {
    var field = `
      <div id="plugin-fields">
        <div class="form-row plugin-field">
        <div class="form-group col">
          <input type="text" class="form-control" name="pluginName" placeholder="Name">
        </div>
        <div class="form-group col">
          <input type="text" class="form-control" name="pluginVersion" placeholder="Version">
        </div>
        <div class="form-group col-auto">
          <button type="button" class="btn btn-danger remove-field-btn">Remove</button>
        </div>
        </div>
      </div>
    `;
    $("#plugin-fields").append(field);
  });

  // remove field
  $("#addPluginModal").on("click", ".remove-field-btn", function () {
    $(this).closest(".form-row").remove();
  });

  $("#add-server-modal-submit").on("click", function () {
    var formData = $("#add-new-server-form").serializeArray();
    var pluginsJson = JSON.stringify(formData);
    console.log(pluginsJson)
    $.ajax({
      type: "POST",
      url: "your-server-page1.php",
      data: pluginsJson,
      success: function (response) {
        // handle success response
      },
      error: function (xhr, status, error) {
        // handle error response
      },
    });
    $("#addServerModal").modal("hide");
  });



  $("#add-plugin-submit").on("click", function (event) {
    var formData = $("#add-plugin-form").serializeArray();
    var pluginsJson = JSON.stringify(formData);
    console.log(pluginsJson)
    $.ajax({
      type: "POST",
      url: "11122",
      data: pluginsJson,
      success: function (response) {
        // handle success response
      },
      error: function (xhr, status, error) {
        // handle error response
      },
    });
    $("#addPluginModal").modal("hide");
  });

  $("#change-version-submit").on("click", function () {
    var formData = $("#change-version-form").serializeArray();
    var pluginsJson = JSON.stringify(formData);
    console.log(pluginsJson)
    $.ajax({
      type: "POST",
      url: "your-asd",
      data: pluginsJson,
      success: function (response) {
        // handle success response
      },
      error: function (xhr, status, error) {
        // handle error response
      },
    });
    $("#changeVersionModal").modal("hide");
  });

});

$(document).on("click", ".change-version-btn", function () {
  var version = $(this).closest("li").find("span").text().split(":")[1];
  $("#changeVersionModal").find("#plugin-version").val(version);
});




// submit form




