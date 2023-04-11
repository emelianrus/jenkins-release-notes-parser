import React, { useState } from "react";

function JsonEditor({ data }) {
  const [jsonData, setJsonData] = useState(data);
  const [isEditing, setIsEditing] = useState(true);

  const handleEdit = () => {
    setIsEditing(true);
  };

  const handleSave = () => {
    setIsEditing(false);
    // Send the updated JSON data to the server or perform any other actions here
  };

  const handleCancel = () => {
    setJsonData(data);
    setIsEditing(false);
  };

  const handleInputChange = (event) => {
    setJsonData(JSON.parse(event.target.value));
  };

  return (
    <div>
      <pre>{JSON.stringify(jsonData, null, 2)}</pre>
      {isEditing ? (
        <>
          <textarea
            value={JSON.stringify(jsonData)}
            onChange={handleInputChange}
            rows={10}
            cols={50}
          />
          <button onClick={handleSave}>Save</button>
          <button onClick={handleCancel}>Cancel</button>
        </>
      ) : (
        <button onClick={handleEdit}>Edit</button>
      )}
    </div>
  );
}

export default JsonEditor;