import AceEditor from 'react-ace';
import 'ace-builds/src-noconflict/mode-java';
import 'ace-builds/src-noconflict/mode-json';
import 'ace-builds/src-noconflict/theme-monokai';
import React, { useState, useEffect } from "react";

function JsonEditor() {

  const [isModified, setIsModified] = useState(false);
  const [editorValue, setEditorValue] = useState(JSON.stringify());
  const [errorMessage, setErrorMessage] = useState("");


  useEffect(() => {
    async function fetchData() {
      try {
        // TODO: https://api.github.com/repos/OWNER/REPO/releases
        // repos/:owner/:repo/releases
        // /project/:owner/:repo/releases
        const response = await fetch(`http://localhost:8080/watcher-list`);

        const data = await response.json();
        // pass as new single object instead of several params
        if (data) {
          setEditorValue(JSON.stringify(data, null, 2));
        }

      } catch (error) {
        console.error(error);
      }
    }

    fetchData();
  }, []);


  const handleSave = () => {
    try {
      JSON.parse(editorValue);
      setErrorMessage("");
      handleClick()
    } catch (e) {
      setErrorMessage(e.message);
    }

    setIsModified(false);
  }

  const handleChange = (value) => {
    setEditorValue(value);
    setIsModified(true);
  }

  function handleClick() {
    fetch('http://localhost:8080/watcher-list', {  // Enter your IP address here
      method: 'POST',
      mode: 'cors',
      body: editorValue
    })
  }

  // https://github.com/securingsincity/react-ace/blob/master/docs/Ace.md#available-props
  return (
    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
      <div id="error-message" style={{ backgroundColor: 'red', color: 'white' }}>{errorMessage}</div>
      <AceEditor
        mode="java"
        theme="monokai"
        onChange={(value) => {
          handleChange(value)
        }}
        name="my-editor"
        editorProps={{ $blockScrolling: true }}
        value={editorValue}
        height="800px"
        width="800px"
        fontFamily= "tahoma"
        fontSize= "14pt"
        // tabSize="2"
      />
      <button onClick={handleSave}>Save</button>
      <p>Status: {isModified ? 'unsaved' : 'saved'}</p>
    </div>
  );
}

export default JsonEditor;