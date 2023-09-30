import 'bootstrap/dist/css/bootstrap.min.css';

import React, { useState, useEffect } from "react";
import PluginChangesCard from '../components/PluginChangesCard';
import ReleaseNotesList from '../components/ReleaseNotesList';
import Button from 'react-bootstrap/Button';

function PluginChanges() {
  const [projects, setProjects] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [backendStatus, setBackendStatus] = useState("");

  useEffect(() => {
    fetchDiffUpdated();
  }, []);

  const fetchDiffUpdated = async () => {
    try {
      const response = await fetch(`http://localhost:8080/plugin-manager/get-fixed-deps-diff`);
      const data = await response.json();

      const sortedData = data.sort((a, b) => {
        const nameA = a.Name.toLowerCase();
        const nameB = b.Name.toLowerCase();

        if (nameA < nameB) {
          return -1;
        }
        if (nameA > nameB) {
          return 1;
        }

        return 0;
      });

      setProjects(sortedData);
      setIsLoading(false);
    } catch (error) {
      console.error(error);
      setBackendStatus(error.message);
    }
  };

  const handleGetTxtFile = async () => {
    try {
      const response = await fetch('http://localhost:8080/plugin-changes/download-file', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/octet-stream',
        },
      });

      if (response.ok) {
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = 'file.txt';
        link.click();
        window.URL.revokeObjectURL(url);
      } else {
        console.error('Failed to download file:', response.status);
      }
    } catch (error) {
      console.error('Error while downloading file:', error);
    }
  };

  function handleDoApply() {
    const result = projects.map(project => `${project.Name}:${project.NewVersion}`).join('\n');
    fetch('http://localhost:8080/add-plugin-list/add-plugins', {
      method: 'POST',
      mode: 'cors',
      body: result
    })

    window.location.replace('/plugin-manager');
  }

  async function checkDiffsWithUpgrade() {
    try {
      await fetch(`http://localhost:8080/plugin-manager/check-deps-with-update`);
      fetchDiffUpdated();
    } catch (error) {
      console.error(error);
    }
  }

  const pluginsArray = [];
  for (const key in projects) {
    pluginsArray.push({
      key,
      project: projects[key]
    });
  }

  const pluginCards = pluginsArray.map(plugin => (
    <PluginChangesCard key={plugin.key} project={plugin.project} />
  ));
  const footerStyle = {
    marginTop: '20px', // Adjust the margin as needed
    backgroundColor: '#f0f0f0', // Set your desired background color
    padding: '10px', // Set padding if needed
    textAlign: 'center',
  };
  return (
    <div>
      <div className="project-list">
      <div className="container-sm mt-5 ml-5 mr-5">
        <h3>Plugin changes</h3>
        <div className="table-responsive">
          <Button variant="primary" style={{ marginRight: '10px' }} onClick={fetchDiffUpdated}>
            Show diff between updated
          </Button>
          <Button variant="primary" style={{ marginRight: '10px' }} onClick={checkDiffsWithUpgrade}>
            Get deps plugins (with core update)
          </Button>
          <Button variant="primary" disabled>
            Get deps plugins (without core update)
          </Button>
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

                  {projects === undefined
                    ? <tbody><tr><td>No projects to display</td></tr></tbody>
                    : pluginCards
                  }

                </table>

                <div style={{ display: 'flex', justifyContent: 'center' }}>
                  <Button variant="primary" onClick={handleGetTxtFile}>Get Txt File</Button>
                  <div style={{ width: '10px' }}></div> {/* Margin */}
                  <Button variant="primary" onClick={handleDoApply}>Apply to plugin manager list</Button>
                  <div style={{ width: '10px' }}></div> {/* Margin */}
                  {/* <Button variant="warning" onClick={forceRescan}>Force rescan</Button> */}
                </div>

                <b>RELEASE NOTES</b>
                <ReleaseNotesList projects={projects}/>
              </>
            )}
          </div>
        </div>

        <footer style={footerStyle}>
          © 2023 Jenkins Plugin Manager
        </footer>
      </div>
    </div>
  );
}

export default PluginChanges;