
import AceEditor from 'react-ace';
import 'ace-builds/src-noconflict/mode-java';
import 'ace-builds/src-noconflict/mode-json';
import 'ace-builds/src-noconflict/theme-monokai';
import React, { useState, useEffect } from "react";


const AddUpdatedPlugins = () => {

  const [isModified, setIsModified] = useState(false);

  // TODO: remove me
  const dummyValue = ``

  const [editorValue, setEditorValue] = useState(dummyValue);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const response = await fetch(`http://localhost:8080/add-updated-plugins/get-data`);
      const data = await response.json();
      const convertedData = Object.entries(data)
        .map(([key, value]) => `${key}:${value}`)
        .join('\n');
      setEditorValue(convertedData);
    } catch (error) {
      console.error(error);
    }
  };

  const handleSave = () => {
    handleClick()
    setIsModified(false);
  }

  const handleChange = (value) => {
    setEditorValue(value);
    setIsModified(true);
  }

  function handleClick() {
    fetch('http://localhost:8080/add-updated-plugins/edit-data', {
      method: 'POST',
      mode: 'cors',
      body: editorValue
    })
    window.location.replace('/plugin-manager');
  }

  return (
    <div className="editor">
      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        {/* <div id="error-message" style={{ backgroundColor: 'red', color: 'white' }}>{errorMessage}</div> */}
        <AceEditor
          mode="java"
          theme="monokai"
          onChange={(value) => {
            handleChange(value)
          }}
          name="my-editor"
          editorProps={{ $blockScrolling: false }}
          value={editorValue}
          height="800px"
          width="800px"
          fontFamily= "tahoma"
          fontSize= "14pt"
        />
        <button onClick={handleSave}>Save</button>
        <p>Status: {isModified ? 'unsaved' : 'saved'}</p>
      </div>
    </div>
  );
};

export default AddUpdatedPlugins;