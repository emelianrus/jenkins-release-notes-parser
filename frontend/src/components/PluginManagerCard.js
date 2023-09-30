import Button from 'react-bootstrap/Button';
import { OverlayTrigger, Popover } from 'react-bootstrap';
import React, { useState } from "react";
import Modal from 'react-bootstrap/Modal';


function PluginManagerCard({ project }) {

  const [manifestAttrs, setManifestAttrs] = useState({});
  const [showGetManifestAttrs, setShowGetManifestAttrs] = useState(false);
  const [getManifestMessage, setGetManifestMessage] = useState("");

  const handleShowGetManifestAttrs = () => {
    fetchManifestData()
    setShowGetManifestAttrs(true);
  }
  const handleCloseShowGetManifestAttrs = () => {
    setShowGetManifestAttrs(false);
  };

  async function fetchManifestData() {
    try {
      setGetManifestMessage("LOADING");
      const response = await fetch('http://localhost:8080/plugin-manager/get-manifest-attrs', {
        method: 'POST',
        // mode: 'cors',
        body: JSON.stringify({
          name: project.Name
        })
      })

      if (!response.ok) {
        let error = await response.json()
        console.log(error.error)
        setGetManifestMessage("Error during download hpi file of plugin 404");
        throw new Error('Network response was not ok');
      }
      const data = await response.json();
      setManifestAttrs(data);
      setGetManifestMessage("")
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  }

  function handleDelete(name){
    fetch('http://localhost:8080/plugin-manager/delete-plugin', {
      method: 'DELETE',
      // mode: 'cors',
      body: JSON.stringify({
        name: name
      })
    })
    // TODO: wrong using of components
    // need state change
    window.location.reload(false);
  }

  const popoverStyle = {
    width: 'auto',
    maxWidth: 'none',
    whiteSpace: 'nowrap',
  };

  const requiredByPopover = (

    <Popover id="popover-basic" style={popoverStyle}>
      <Popover.Body>
        {Object.keys(project.RequiredBy).length === 0 ? (
          <div>don't have any</div>
        ) : (
          Object.keys(project.RequiredBy).map((key,value) => (
            <div key={key}>{key}:{project.RequiredBy[key]}</div>
          ))
        )}
      </Popover.Body>
    </Popover>
  );

  const dependsOnPopover = (
    <Popover id="popover-basic">
      <Popover.Body>
        {Object.keys(project.Dependencies).length === 0 ? (
          <div>don't have any</div>
        ) : (
          Object.keys(project.Dependencies).map((key) => (
            <div key={key}>{key}:{project.Dependencies[key].Version}</div>
          ))
        )}
      </Popover.Body>
    </Popover>
  );

  const warningsPopover = (
    <Popover id="popover-basic">
      <Popover.Body>
        {!project.Warnings || project.Warnings.length === 0 ? (
          <div>Don't have any warnings.</div>
        ) : (
          <ul>
            {project.Warnings.map((warning) => (
              <li key={warning.Id}>
                <div>
                  <strong>{warning.Message}</strong><br /> Affects version {warning.Versions[0].LastVersion} and earlier <br />
                </div>
              </li>
            ))}
          </ul>
        )}
      </Popover.Body>
    </Popover>
  );


  return (
    <tbody>
      <tr id="server-plugins">
        <td>
          <span id="plugin-name-ranged">
            <a href={ project.GITUrl } style={{ color: 'inherit' }}>
              { project.Name }
            </a>
          </span>
        </td>
        <td>
          <ul>
            <li>
              <span>{ project.Version }</span>
            </li>
          </ul>
        </td>
        <td>
          <ul>
            <li>
              {/* from plugin site */}
              <span>{ project.LatestVersion}</span>
            </li>
          </ul>
        </td>
        <td>
          { project.RequiredCoreVersion }
        </td>
        <td>
          <ul>
            <div style={{ textAlign: "center"}}>
              <OverlayTrigger trigger="click" placement="right" overlay={requiredByPopover} rootClose>
                <Button variant="outline-primary" disabled={Object.keys(project.RequiredBy).length === 0}>Required by</Button>
              </OverlayTrigger>
              <OverlayTrigger trigger="click" placement="right" overlay={dependsOnPopover} rootClose>
                <Button variant="outline-primary" disabled={Object.keys(project.Dependencies).length === 0}>Depends on</Button>
              </OverlayTrigger>
              <OverlayTrigger trigger="click" placement="right" overlay={warningsPopover} rootClose>
                <Button variant="outline-primary" disabled={!project.Warnings || project.Warnings.length === 0}>Warnings</Button>
              </OverlayTrigger>
              <Button variant="outline-primary" onClick={handleShowGetManifestAttrs} >ShowManifest</Button>
            </div>

          </ul>
        </td>
        <td>
          <ul>
              <Button variant="danger" onClick={() => handleDelete(project.Name)}>Delete</Button>
          </ul>
        </td>
      </tr>
      <Modal show={showGetManifestAttrs} size="lg" onHide={handleCloseShowGetManifestAttrs}>
          <Modal.Header closeButton>
            <Modal.Title>Modal heading</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            {getManifestMessage !== "" ? (
              <div style={{ fontSize: "44px", fontWeight: "bold", textAlign: "center" }}>{getManifestMessage}</div>
              ) : (
                <>
                  {Object.entries(manifestAttrs).map(([key, value]) => (
                    <li key={key}><b>{key}</b>: {value}</li>
                  ))}
                </>
              )
            }
          </Modal.Body>
          <Modal.Footer>
            <Button onClick={handleCloseShowGetManifestAttrs}>Close</Button>
          </Modal.Footer>
        </Modal>

    </tbody>
  );
}

export default PluginManagerCard;



