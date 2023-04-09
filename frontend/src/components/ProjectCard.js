

function ProjectCard({ project }) {
    return (
      // TODO: fix class name
      <tbody>
        <tr id="server-plugins">
          <td>
              <span id="plugin-name-ranged">{ project.name }</span>
          </td>
          <td>
            <ul>
              <li>
                <span>is Downloaded:</span> <span>{ project.is_downloaded.toString() }</span>
              </li>
              <li>
                <span>Has error:</span> <span>{ project.has_error.toString() }</span>
              </li>
              <li>
                <span>Last updated:</span> <span>{ project.last_updated }</span>
              </li>
            </ul>
          </td>
          <td >
            <span id="plugin-name-ranged">
              <button type="button" class="btn btn-primary rescan-btn">rescan</button>
            </span>
          </td>
        </tr>
      </tbody>
    );
}

export default ProjectCard;


