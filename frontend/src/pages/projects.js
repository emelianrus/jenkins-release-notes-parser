import 'bootstrap/dist/css/bootstrap.min.css';

import ProjectsList from '../components/ProjectsList';

import React, { useState, useEffect } from "react";

function Projects() {
  const [projects, setProjects] = useState([]);

  useEffect(() => {
    setProjects([{
      name: "projectName1",
      owner: "ownername1",
      is_downloaded: true,
      has_error: false,
      last_updated: "Feb12 12:55"
    },
    {
      name: "projectName2",
      owner: "ownername2",
      is_downloaded: false,
      has_error: false,
      last_updated: "Aug12 1:14"
    }])
    // async function fetchData() {
    //   try {
    //     const response = await fetch("http://localhost:8080/projects");
    //     const data = await response.json();
    //     setProjects(data);
    //   } catch (error) {
    //     console.error(error);
    //   }
    // }

    // fetchData();
  }, []);

  return (
    <div>
      <ProjectsList projects={projects} />
    </div>
  );
}

export default Projects;