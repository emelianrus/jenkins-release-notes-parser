
import 'bootstrap/dist/css/bootstrap.min.css';

import React, { useState, useEffect } from "react";
import PluginChangesCard from '../components/PluginChangesCard';
import PotentialUpdatesList from '../components/PotentialUpdatesList';
function PluginChanges() {

  const [pluginsDiff, setPluginsDiff] = useState([]);
  const [projects, setProjects] = useState([]);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {

    // try {
    //   const response = await fetch(`http://localhost:8080/plugin-manager/check-deps`);
    //   const data = await response.json();
    //   console.log(data)
    // } catch (error) {
    //   console.error(error);
    // }

    try {
      const response = await fetch(`http://localhost:8080/plugin-manager/get-fixed-deps-diff`);
      const data = await response.json();
      setPluginsDiff(data);
    } catch (error) {
      console.error(error);
    }

    try {
      // TODO: https://api.github.com/repos/OWNER/REPO/releases
      // repos/:owner/:repo/releases
      // /project/:owner/:repo/releases
      const response = await fetch(`http://localhost:8080/potential-updates`);
      const data = await response.json();

      // pass as new single object instead of several params
      setProjects(data);
    } catch (error) {
      console.error(error);
    }

  };

  const pluginsArray = [];
  for (const key in pluginsDiff) {
    pluginsArray.push({
      key,
      project: pluginsDiff[key]
    });
  }

  const pluginCards = pluginsArray.map(plugin => (
    <PluginChangesCard key={plugin.key} project={plugin.project} />
  ));
  return (
    <div>
      <div className="project-list">
      <div className="container-sm mt-5 ml-5 mr-5">
        <h3>Title</h3>
        <div className="table-responsive">
          <table className="table">
            <thead className="thead-light">
              <tr>
                <th>Project</th>
                <th>From version</th>
                <th>To version</th>
                <th>Type</th>
              </tr>
            </thead>

            {/* {projectCards.map(project => <PluginManagerCard key={project.key} project={project} />)} */}
            {pluginsDiff === undefined
              ? <tbody><tr><td>No projects to display</td></tr></tbody>
              : pluginCards
            }

          </table>
          <b>RELEASE NOTES</b>
          <PotentialUpdatesList projects={projects}/>
        </div>
      </div>
    </div>
    </div>
  );
}

export default PluginChanges;