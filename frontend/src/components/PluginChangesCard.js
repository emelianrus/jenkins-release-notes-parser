
function PluginChangesCard({ project }) {

  let changeType = "unknown"
  if (project.Type === 1) {
    changeType = "new"
  } else if (project.Type === 2) {
    changeType = "update"
  } else if (project.Type === 3) {
    changeType = "the same"
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
            <span>{ project.CurrentVersion }</span>
          </ul>
        </td>
        <td>
          <ul>
            <span>{ project.NewVersion }</span>
          </ul>
        </td>
        <td>
          <ul>
            <span>{ changeType }</span>
          </ul>
        </td>
      </tr>
    </tbody>
  );
}

export default PluginChangesCard;


