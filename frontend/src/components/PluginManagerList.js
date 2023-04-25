
import PluginManagerCard from './PluginManagerCard';
import React from "react";


function PluginManagerList({ projects }) {
  const pluginsArray = [];

  for (const key in projects) {
    pluginsArray.push({
      key,
      project: projects[key]
    });
  }

  const pluginCards = pluginsArray.map(plugin => (
    <PluginManagerCard key={plugin.key} project={plugin.project} />
  ));


  return (
    <div className="project-list">
      <div className="container-sm mt-5 ml-5 mr-5">
        <h3>Title</h3>
        <div className="table-responsive">
          <table className="table">
            <thead className="thead-light">
              <tr>
                <th>Project name</th>
                <th>Installed version</th>
                <th>Latest version</th>
                <th>Download status</th>
                <th>Force Download</th>
              </tr>
            </thead>

            {/* {projectCards.map(project => <PluginManagerCard key={project.key} project={project} />)} */}
            {projects === undefined
              ? <tbody><tr><td>No projects to display</td></tr></tbody>
              : pluginCards
            }
          </table>
        </div>
      </div>
    </div>
  );
}


export default PluginManagerList;