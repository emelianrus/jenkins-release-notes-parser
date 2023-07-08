import 'bootstrap/dist/css/bootstrap.min.css';

import React, { useState, useEffect } from "react";
import PluginChangesCard from '../components/PluginChangesCard';
import ReleaseNotesList from '../components/ReleaseNotesList';

function PluginChanges() {
  const [pluginsDiff, setPluginsDiff] = useState([]);
  const [projects, setProjects] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [backendStatus, setBackendStatus] = useState("Loading...");

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      await fetch(`http://localhost:8080/plugin-manager/check-deps`);
    } catch (error) {
      console.error(error);
      setBackendStatus(error.message);
    }

    try {
      const response = await fetch(`http://localhost:8080/plugin-manager/get-fixed-deps-diff`);
      const data = await response.json();
      setPluginsDiff(data);
      setProjects(data);
      setIsLoading(false);
    } catch (error) {
      console.error(error);
      setBackendStatus(error.message);
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
        <h3>Plugin changes</h3>
        <div className="table-responsive">
        {isLoading ? (
          <div style={{ fontSize: "44px", fontWeight: "bold", textAlign: "center" }}>{backendStatus}</div>
          ) : (
            <>
              <table className="table">
                <thead className="thead-light">
                  <tr>
                    <th>Project</th>
                    <th>From version</th>
                    <th>To version</th>
                    <th>Type</th>
                  </tr>
                </thead>

                {pluginsDiff === undefined
                  ? <tbody><tr><td>No projects to display</td></tr></tbody>
                  : pluginCards
                }

              </table>
              <b>RELEASE NOTES</b>
              <ReleaseNotesList projects={projects}/>
            </>
          )}
        </div>
      </div>
    </div>
    </div>
  );
}

export default PluginChanges;