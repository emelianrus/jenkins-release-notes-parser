import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';

import NavBar from './components/NavBar';
import ProjectList from './components/ProjectList';

import React, { useState, useEffect } from "react";


function App() {
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
      <NavBar />
      <ProjectList projects={projects} />
    </div>
  );
}

export default App;