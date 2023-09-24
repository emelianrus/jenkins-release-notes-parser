import 'bootstrap/dist/css/bootstrap.min.css';
import Button from 'react-bootstrap/Button';
import Modal from 'react-bootstrap/Modal';
import PluginManagerList from '../components/PluginManagerList';
import Form from 'react-bootstrap/Form';
import React, { useState, useEffect } from "react";
import { Link } from 'react-router-dom';

function PluginManager() {

  const [plugins, setPlugins] = useState([]);

  const [showCoreVersion, setShowCoreVersion] = useState(false);
  const [showAddNewPlugin, setAddNewPlugin] = useState(false);
  // status should represent status of all plugins
  // are deps were resolved, are all deps added, are all plugins downloaded
  // const [showSyncStatus, setSyncStatus] = useState("Not in sync");

  const [coreVersion, setCoreVersion] = useState("2.222.2");

  const handleCloseCoreVersion = () => setShowCoreVersion(false);
  const handleEditCoreVersion = () => setShowCoreVersion(true);

  const handleCloseAddNewPlugin = () => {
    setAddNewPlugin(false);
    setAddPluginError('');
  };
  const handleAddNewPlugin = () => setAddNewPlugin(true);

  const [addPluginError, setAddPluginError] = useState("");

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const response = await fetch(`http://localhost:8080/plugin-manager/get-data`);
      const data = await response.json();

      setPlugins(data.Plugins);
      setCoreVersion(data.CoreVersion);
    } catch (error) {
      console.error(error);
    }
  };

  const handleSaveAddNewPlugin = async () => {
    try {
      const pluginName = document.querySelector('input[type="pluginName"]').value;
      const pluginVersion = document.querySelector('input[type="pluginVersion"]').value;

      const requestData = {
        name: pluginName,
        version: pluginVersion
      };

      const response = await fetch('http://localhost:8080/plugin-manager/add-new-plugin', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
      });

      const responseData = await response.json();

      if (!response.ok) {
        setAddPluginError(responseData.message);
        throw new Error(responseData.message);
      }

      // Reload plugins after adding a new plugin
      setAddPluginError("");
      fetchData();
      handleCloseAddNewPlugin(false);
    } catch (error) {
      console.error(error);
    }
  };

  const handleEditCoreSubmit = (event) => {
    event.preventDefault();
    const newCoreVersion = event.target.elements.coreVersion.value;

    fetch('http://localhost:8080/plugin-manager/edit-core-version', {
      method: 'POST',
      // mode: 'cors',
      body: JSON.stringify({
        name: newCoreVersion
      })
    })

    setCoreVersion(newCoreVersion);
    handleCloseCoreVersion();
  };

  const handleGetTxtFile = async () => {
    try {
      const response = await fetch('http://localhost:8080/plugin-manager/download-file', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/octet-stream',
        },
      });

      if (response.ok) {
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = 'file.txt'; // Set the desired file name
        link.click();
        window.URL.revokeObjectURL(url);
      } else {
        console.error('Failed to download file:', response.status);
      }
    } catch (error) {
      console.error('Error while downloading file:', error);
    }
  };



  return (
    <div>
      {/* buttons menu */}

      <div className="project-list">
        <div className="container-sm mt-5 ml-5 mr-5">
          <div className="row justify-content-end">
            <div className="col-auto">
              Jenkins core version: {coreVersion}
            </div>
            <div className="col-auto">
              <Button variant="primary" onClick={handleEditCoreVersion}>Change version</Button>
              {/* <div>
                Status: {showSyncStatus}
              </div> */}
            </div>
          </div>
        </div>
      </div>
{/* buttons menu end */}

      <PluginManagerList projects={plugins} />

      <div className="container-sm mt-5 ml-5 d-flex">
        <div className="mr-2">
          <Link to="/add-plugin-list">
            <Button variant="outline-primary">Edit plugin manager list</Button>
          </Link>
        </div>
        <div className="mr-2">
          <Button variant="outline-primary" onClick={handleAddNewPlugin}>Add one plugin</Button>
          <Button variant="outline-primary" onClick={handleGetTxtFile} disabled={ Object.keys(plugins).length === 0 } >Get Txt File</Button>
        </div>
      </div>
      {/* ADD NEW PLUGIN MODAL WINDOW */}
      <Modal show={showAddNewPlugin} onHide={handleCloseAddNewPlugin}>
        <Modal.Header closeButton>
          <Modal.Title>Edit Jenkins core version</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form>
            {addPluginError && <div style={{ color: 'red' }}>{addPluginError}</div>}
            <Form.Group className="mb-3" controlId="exampleForm.ControlInput2">
            <Form.Label>Plugin name</Form.Label>
              <Form.Control
                type="pluginName"
                placeholder=""
                autoFocus
              />

            <Form.Label>Plugin version</Form.Label>
              <Form.Control
                type="pluginVersion"
                placeholder=""
              />
            </Form.Group>
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={handleCloseAddNewPlugin}>
            Close
          </Button>
          <Button variant="primary" onClick={handleSaveAddNewPlugin}>
            Save Changes
          </Button>
        </Modal.Footer>
      </Modal>


      {/* EDIT CORE VERSION MODAL WINDOW */}
      <Modal show={showCoreVersion} onHide={handleCloseCoreVersion}>
        <Modal.Header closeButton>
          <Modal.Title>Edit Jenkins core version</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form onSubmit={handleEditCoreSubmit}>
            <Form.Group className="mb-3" controlId="exampleForm.ControlInput1">
              <Form.Control
                type="text"
                name="coreVersion"
                defaultValue={coreVersion}
                autoFocus
              />
            </Form.Group>
            <Modal.Footer>
              <Button variant="secondary" onClick={handleCloseCoreVersion}>
                Close
              </Button>
              <Button variant="primary" type="submit">
                Save Changes
              </Button>
            </Modal.Footer>
          </Form>
        </Modal.Body>
      </Modal>

    </div>
  );
}

export default PluginManager;