import './releases.css';
import 'bootstrap/dist/css/bootstrap.min.css';

import ReleaseNotesList from '../components/ReleaseNotesList';

import React, { useState, useEffect } from "react";

function Releases() {
  const [projects, setProjects] = useState([]);

  useEffect(() => {
    async function fetchData() {
      try {
        const response = await fetch("http://localhost:8080/projects");
        const data = await response.json();
        setProjects(data);
      } catch (error) {
        console.error(error);
      }
    }

    fetchData();
  }, []);

  return (
    <div>
      <ReleaseNotesList projects={projects} />
    </div>
  );
}

export default Releases;