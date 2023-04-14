import './releases.css';
import 'bootstrap/dist/css/bootstrap.min.css';

import PotentialUpdatesList from '../components/PotentialUpdatesList';

import React, { useState, useEffect } from "react";
import { useParams } from 'react-router-dom';

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

        console.log(data)
        // pass as new single object instead of several params
        setProjects(data);

      } catch (error) {
        console.error(error);
      }
    }

    fetchData();
  }, []);



  // const tableRows = Object.keys(projects).map(key => (
  //   projects[key].map(item => (
  //     <tr key={item.Name}>
  //       <td>{key}</td>
  //       <td>{item.Name}</td>
  //     </tr>
  //   ))
  // ));

  return (
    <div>
      <PotentialUpdatesList projects={projects}/>
    </div>
  );
}

export default PotentialUpdates;