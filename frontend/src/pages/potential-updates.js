import './releases.css';
import 'bootstrap/dist/css/bootstrap.min.css';

import PotentialUpdatesList from '../components/PotentialUpdatesList';

import React, { useState, useEffect } from "react";

function PotentialUpdates() {

  const [projects, setProjects] = useState([]);

  useEffect(() => {
    async function fetchData() {
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
    }

    fetchData();
  }, []);



  return (
    <div>
      <PotentialUpdatesList projects={projects}/>
    </div>
  );
}

export default PotentialUpdates;