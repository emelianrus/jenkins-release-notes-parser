import React, { useState } from 'react';
import AceEditor from 'react-ace';
import 'ace-builds/src-noconflict/mode-java';
import 'ace-builds/src-noconflict/mode-json';
import 'ace-builds/src-noconflict/theme-monokai';

function JsonEditor({ data }) {

  const [isModified, setIsModified] = useState(false);
  const [editorValue, setEditorValue] = useState(JSON.stringify(data, null, 2));
  const [errorMessage, setErrorMessage] = useState("");

  const handleSave = () => {
    try {
      JSON.parse(editorValue);
      setErrorMessage("");
    } catch (e) {
      setErrorMessage(e.message);
    }

    setIsModified(false);
  }

  const handleChange = (value) => {
    setEditorValue(value);
    setIsModified(true);
  }

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
        height="500px"
        width="500px"
      />
      <button onClick={handleSave}>Save</button>
      <p>Status: {isModified ? 'unsaved' : 'saved'}</p>
    </div>
  );
}

export default JsonEditor;