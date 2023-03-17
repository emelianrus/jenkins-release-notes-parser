console.log("loaded servers.js");

$(document).on("click", ".add-new-plugin-btn", function () {
  var serverName = $(this).data("server-name");
  console.log(serverName);
  $("#addPluginModal").find("#jenkins-server-name").val(serverName);
});

$(document).on("click", ".change-version-btn", function () {
  var version = $(this).closest("li").find("span").text().split(":")[1];
  $("#changeVersionModal").find("#plugin-version").val(version);
});


$(document).on("click", ".change-version-btn", function () {
  var serverName = $(this).data("server-name");
  $("#changeVersionModal").find("#jenkins-server-name").val(serverName);
  var pluginName =  $(this).data("plugin-name");
  $("#changeVersionModal").find("#jenkins-plugin-name").val(pluginName);
});


$(document).ready(function () {
  // add a click event listener to the plugins list (the parent element)
  $(".delete-server-btn").on("click", function () {
    var serverName = $(this).data("server-name");

    console.log(serverName);
    var block = $(this).closest("tr");

    // send an AJAX request to the server
    $.ajax({
      url: "/delete-server", // replace this with the actual URL of the delete endpoint
      type: "POST",
      contentType: "application/json",
      data: JSON.stringify(serverName),
      success: function () {
        block.remove();
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

    var jsonObject = {};
    // Loop through the JSON array and add each key-value pair to the JSON object
    for (let i = 0; i < formData.length; i++) {
        const item = formData[i];
        jsonObject[item.name] = item.value;
    }

    var payload = {
      jenkinsName: jsonObject["serverName"],
      coreVersion: jsonObject["coreVersion"]
    }

    $.ajax({
      type: "POST",
      url: "/add-server",
      data: JSON.stringify(payload),
      success: function () {
        // handle success response
        // TODO: add server live
        window.location.reload();
      },
      error: function (xhr, status, error) {
        // handle error response
      },
    });
    $("#addServerModal").modal("hide");
  });

  $("#add-plugin-submit").on("click", function () {
    var formData = $("#add-plugin-form").serializeArray();
    const formDataConverted = formData.reduce((obj, { name, value }) => {
      if (name === "jenkinsServerName") {
        obj.jenkinsName = value;
      } else if (name === "pluginName" || name === "pluginVersion") {
        if (!obj.plugins) {
          obj.plugins = [];
        }
        obj.plugins.push({ name, value });
      }
      return obj;
    }, {});

    var payload = {
      jenkinsName: formDataConverted["jenkinsName"],
      plugins: formDataConverted["plugins"]
    };

    $.ajax({
      type: "POST",
      url: "/add-new-plugin",
      data: JSON.stringify(payload),
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

    const formDataConverted = formData.reduce((obj, { name, value }) => {
      obj[name] = value;
      return obj;
    }, {});

    var payload = {
      jenkinsName: formDataConverted["jenkinsServerName"],
      pluginName: formDataConverted["jenkinsPluginName"],
      newPluginVersion: formDataConverted["pluginVersion"]
    };

    $.ajax({
      type: "POST",
      url: "/change-plugin-version",
      data: JSON.stringify(payload),
      success: function (response) {
        // handle success response
        // TODO: change table version live
        window.location.reload();
      },
      error: function (xhr, status, error) {
        // handle error response
        console.log(error);
      },
    });
    $("#changeVersionModal").modal("hide");
  });

});

