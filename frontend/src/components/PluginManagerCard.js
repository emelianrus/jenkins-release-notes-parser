import Button from 'react-bootstrap/Button';
import { OverlayTrigger, Popover } from 'react-bootstrap';

function ProjectCard({ project }) {

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

  const requiredByPopover = (
    <Popover id="popover-basic">
      <Popover.Body>
        {Object.keys(project.RequiredBy).length === 0 ? (
          <div>don't have any</div>
        ) : (
          Object.keys(project.RequiredBy).map((key) => (
            <div key={key}>{key}</div>
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
            <div key={key}>{key}</div>
          ))
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
            { project.Name }
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
              <OverlayTrigger trigger="focus" placement="right" overlay={requiredByPopover}>
                <Button variant="outline-primary">Required by</Button>
              </OverlayTrigger>
              <OverlayTrigger trigger="focus" placement="right" overlay={dependsOnPopover}>
                <Button variant="outline-primary">Depends on</Button>
              </OverlayTrigger>
            </div>

          </ul>
        </td>
        <td>
          <span id="plugin-name-ranged">
            <Button variant="primary" onClick={() => handleDoRescan(project.Name, project.Version)}>scan now</Button>
          </span>
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

export default ProjectCard;


