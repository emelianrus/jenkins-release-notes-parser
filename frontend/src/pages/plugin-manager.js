import './releases.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import Button from 'react-bootstrap/Button';
import Modal from 'react-bootstrap/Modal';
import PluginManagerList from '../components/PluginManagerList';
import Form from 'react-bootstrap/Form';
import React, { useState, useEffect } from "react";

function PluginManager() {

  const [plugins, setPlugins] = useState([]);
  const [jenkinsCoreVersion, setJenkinsCoreVersion] = useState("");

  let data = [
    {
      IsInWatcherList: true,
      Project: {
        Name: "kubernetes",
        Owner: "jenkinsci",
        Error: "",
        IsDownloaded: true,
        LastUpdated: '"18 April 2023 17:08:04"',
        ReleaseNotes: null,
      },
    },
    {
      IsInWatcherList: true,
      Project: {
        Name: "antisamy-markup-formatter",
        Owner: "jenkinsci",
        Error: "",
        IsDownloaded: true,
        LastUpdated: '"18 April 2023 17:08:04"',
        ReleaseNotes: null,
      },
    },
    {
      IsInWatcherList: true,
      Project: {
        Name: "ant",
        Owner: "jenkinsci",
        Error: "",
        IsDownloaded: true,
        LastUpdated: '"18 April 2023 17:08:04"',
        ReleaseNotes: null,
      },
    },
  ];

  const [showCoreVersion, setShowCoreVersion] = useState(false);
  const [showAddNewPlugin, setAddNewPlugin] = useState(false);

  const [coreVersion, setCoreVersion] = useState("2.222.2");

  const handleClose = () => setShowCoreVersion(false);
  const handleEditCoreVersion = () => setShowCoreVersion(true);

  const handleCloseAddNewPlugin = () => setAddNewPlugin(false);
  const handleAddNewPlugin = () => setAddNewPlugin(true);


  useEffect(() => {
    async function fetchData() {
      try {
        // TODO: https://api.github.com/repos/OWNER/REPO/releases
        // repos/:owner/:repo/releases
        // /project/:owner/:repo/releases
        const response = await fetch(`http://localhost:8080/plugin-manager/get-data`);
        const data = await response.json();

        // pass as new single object instead of several params
        setPlugins(data.Plugins);
        setCoreVersion(data.CoreVersion)

      } catch (error) {
        console.error(error);
      }
    }

    fetchData();
  }, []);



  const handleSubmit = (event) => {
    event.preventDefault();
    const newCoreVersion = event.target.elements.coreVersion.value;
    setCoreVersion(newCoreVersion);
    handleClose();
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
            <Button variant="primary" onClick={handleEditCoreVersion} >Change version</Button>
            </div>
          </div>
        </div>
      </div>
{/* buttons menu end */}

      <PluginManagerList projects={plugins} />


      <div className="container-sm mt-5 ml-5 mr-">
        <Button variant="outline-primary" onClick={handleAddNewPlugin}>Add new plugin</Button>
      </div>


      {/* ADD NEW PLUGIN MODAL WINDOW */}
      <Modal show={showAddNewPlugin} onHide={handleCloseAddNewPlugin}>
        <Modal.Header closeButton>
          <Modal.Title>Edit Jenkins core version</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form>
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
                autoFocus
              />
            </Form.Group>
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={handleCloseAddNewPlugin}>
            Close
          </Button>
          <Button variant="primary" onClick={handleCloseAddNewPlugin}>
            Save Changes
          </Button>
        </Modal.Footer>
      </Modal>


      {/* EDIT CORE VERSION MODAL WINDOW */}
      <Modal show={showCoreVersion} onHide={handleClose}>
        <Modal.Header closeButton>
          <Modal.Title>Edit Jenkins core version</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form onSubmit={handleSubmit}>
            <Form.Group className="mb-3" controlId="exampleForm.ControlInput1">
              <Form.Control
                type="text"
                name="coreVersion"
                defaultValue={coreVersion}
                autoFocus
              />
            </Form.Group>
            <Modal.Footer>
              <Button variant="secondary" onClick={handleClose}>
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