
function PluginChangesCard({ project }) {

  let changeType = "unknown"
  let backgroundColor = "#ffffff"

  if (project.Type === 1) {
    changeType = "new"
    backgroundColor = "#e7f8e5"
  } else if (project.Type === 2) {
    changeType = "update"
    backgroundColor = "#fee8a3"
  } else if (project.Type === 3) {
    changeType = "the same"
    backgroundColor = "#ffffff"
  } else if (project.Type === 5) {
    changeType = "removed"
    backgroundColor = "#f77c7c"
  }

  const rowStyles = {
    backgroundColor: backgroundColor
  };

  return (
    <tbody>
      <tr id="server-plugins" style={rowStyles}>
        <td>
          <span id="plugin-name-ranged">
            <a href={ project.HTMLURL } style={{ color: 'inherit' }}>
              { project.Name }
            </a>
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


