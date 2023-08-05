import Button from 'react-bootstrap/Button';
import { OverlayTrigger, Popover } from 'react-bootstrap';

function ProjecManagerCard({ project }) {

  function handleDoRescan(name, version){
    fetch('http://localhost:8080/plugin-manager/rescan', {
      method: 'POST',
      // mode: 'cors',
      body: JSON.stringify({
        name: name,
        version: version
      })
    })
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
    // TODO: fix class name
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
                <Button variant="outline-primary">Required by</Button>
              </OverlayTrigger>
              <OverlayTrigger trigger="click" placement="right" overlay={dependsOnPopover} rootClose>
                <Button variant="outline-primary">Depends on</Button>
              </OverlayTrigger>
              <OverlayTrigger trigger="click" placement="right" overlay={warningsPopover} rootClose>
                <Button variant="outline-primary">Warnings</Button>
              </OverlayTrigger>
            </div>

          </ul>
        </td>
        <td>
          <ul>
              <Button variant="danger" onClick={() => handleDelete(project.Name)}>Delete</Button>
          </ul>
        </td>
      </tr>
    </tbody>
  );
}

export default ProjecManagerCard;


