import './releases.css';
import 'bootstrap/dist/css/bootstrap.min.css';

import ReleaseNotesList from '../components/ReleaseNotesList';

import React, { useState, useEffect } from "react";

function Releases() {
  const [projects, setProjects] = useState([]);
  const [projectName, setProjectName] = useState("");
  const [projectGroup, setProjectGroup] = useState("");

  useEffect(() => {
    async function fetchData() {
      try {
        const response = await fetch("http://localhost:8080/project/jenkins-plugin-name/release-notes");
        const data = await response.json();

        // pass as new single object instead of several params
        setProjects(data.ReleaseNotes);
        setProjectName(data.Name);
        setProjectGroup(data.ProjectGroup);

      } catch (error) {
        console.error(error);
      }
    }

    fetchData();
  }, []);

  return (
    <div>
      <ReleaseNotesList projects={projects} projectName={projectName} projectGroup={projectGroup} />
    </div>
  );
}

export default Releases;