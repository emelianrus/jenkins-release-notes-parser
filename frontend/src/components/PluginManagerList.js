
import PluginManagerCard from './PluginManagerCard';
import React from "react";


function PluginManagerList({ projects }) {

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

            {projects === undefined
              ? <tbody><tr><td>No projects to display</td></tr></tbody>
              : projects.map(project => <PluginManagerCard key={project.Project.Name} project={project.Project} />)
            }
          </table>
        </div>
      </div>
    </div>
  );
}


export default PluginManagerList;