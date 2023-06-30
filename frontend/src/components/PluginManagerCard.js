import Button from 'react-bootstrap/Button';

function ProjectCard({ project }) {

  function handleDoRescan(name, version){
    console.log(name);
    console.log(version);
    fetch('http://localhost:8080/plugin-manager/rescan', {
      method: 'POST',
      mode: 'cors',
      body: JSON.stringify({
        name: name,
        version: version
      })
    })
  }

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
              <span>PLUGIN VERSION latest</span>
            </li>
          </ul>
        </td>
        <td>
          <ul>
            <li>
              <span>Is dependency of:</span> <span>BUTTON</span>
            </li>
            <li>
              <span>Depends on:</span> <span>BUTTON</span>
            </li>
          </ul>
        </td>
        <td>
          <span id="plugin-name-ranged">
            <Button variant="primary" onClick={() => handleDoRescan(project.Name, project.Version)}>scan now</Button>
          </span>
        </td>
      </tr>
    </tbody>
  );
}

export default ProjectCard;


