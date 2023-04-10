import { Link } from 'react-router-dom';


function ProjectCard({ project }) {

    return (
      // TODO: fix class name
      <tbody>
        <tr id="server-plugins">
          <td>
              <span id="plugin-name-ranged">
                <Link to={`/test/${project.Owner}/${project.Name}/releases`}>
                  { project.Name }
                </Link>
              </span>
          </td>
          <td>
            <ul>
              <li>
                <span>Is Downloaded:</span> <span>{ project.IsDownloaded.toString() }</span>
              </li>
              <li>
                <span>Has Error:</span> <span>{ project.Error }</span>
              </li>
              <li>
                <span>Last Updated:</span> <span>{ project.LastUpdated }</span>
              </li>
            </ul>
          </td>
          <td >
            <span id="plugin-name-ranged">
              <button type="button" className="btn btn-primary rescan-btn">rescan</button>
            </span>
          </td>
        </tr>
      </tbody>
    );
}

export default ProjectCard;


