import 'bootstrap/dist/css/bootstrap.min.css';

import ProjectsList from '../components/ProjectsList';

import React, { useState, useEffect } from "react";

function Projects() {
  const [projects, setProjects] = useState([]);

  useEffect(() => {
    async function fetchData() {
      try {
        // TODO: https://api.github.com/repos/OWNER/REPO/releases
        // repos/:owner/:repo/releases
        // /project/:owner/:repo/releases
        const response = await fetch(`http://localhost:8080/projects`);

        const data = await response.json();
        // pass as new single object instead of several params
        if (data) {
          setProjects(data);
        }
        data.sort((a, b) => a.Name.localeCompare(b.Name));

      } catch (error) {
        console.error(error);
      }
    }

    fetchData();
  }, []);

  return (
    <div>
      <ProjectsList projects={projects} />
    </div>
  );
}

export default Projects;