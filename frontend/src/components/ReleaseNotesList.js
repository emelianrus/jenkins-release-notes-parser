import './ReleaseNotesList.css';

import ReleaseNoteCard from './ReleaseNotesCard';
import React, { useState, useEffect } from "react";
import Button from 'react-bootstrap/Button';

function ReleaseNotesList({ projects }) {

  const [releaseNotes, setReleaseNotes] = useState([]);
  // const [releaseNotesFetched, setReleaseNotesFetched] = useState(false);
  const [backendStatus, setBackendStatus] = useState("");

  const fetchReleaseNotes = async () => {
    for (const project of projects) {

      try {
        const controller = new AbortController();
        const signal = controller.signal;

        setBackendStatus("Loading, you may hit github rate limit(50 request - unauthorized, 5000 requests - authorized per hour) so need to wait 1h+");
        const response = await fetch(`http://localhost:8080/plugin-manager/get-release-notes-diff`, {
          method: 'POST',
          mode: 'cors',
          body: JSON.stringify({
              name: project.Name,
              versions: project.DiffVersions
          }),
          signal, // Pass the signal to the fetch options
        });
        const data = await response.json();
        console.log(data)
        const isUniqueName = !releaseNotes.some(existingData => existingData.Name === data.Name);

        if (isUniqueName) {
          // Update releaseNotes based on previous state
          setReleaseNotes(prevReleaseNotes => [...prevReleaseNotes, data]);
        } else {
          console.log('Data with the same name already exists in releaseNotes.');
        }
        setBackendStatus("");
      } catch (error) {
        console.error(error);
      }
    }
  }
  let resultNotes = [];
  for (let project of releaseNotes) {


    if (project.ReleaseNotes == null) {
      continue
    }
    else if (project.ReleaseNotes.length === 0){

      resultNotes.push({
        ReleaseNotes: {
          "notusedname": {
            BodyHTML:"<b>RELEASES NOT FOUND</b>"
          }

        }
      })
      continue
    }

    resultNotes.push(project)
  }


  // if (resultNotes.length === 0) {
  //   return <p><b>No updates found.</b></p>;
  // }
  let projectList = [];
  for (let project of resultNotes) {
    for (const [key, value] of Object.entries(project.ReleaseNotes)) {
      projectList.push(
        <ReleaseNoteCard key={project.Name + value.Name} project={value} projectName={project.Name}/>
      )
    }
  }

  return (
    <div className="project-list">
      <div style={{ display: 'flex', justifyContent: 'center' }}>
        <Button variant="primary" onClick={fetchReleaseNotes}>Fetch Release Notes</Button>
      </div>
      <div className="clearfix container-xl px-3 px-md-4 px-lg-5 mt-4">
        {projectList}
      </div>
      <div style={{ display: 'flex', justifyContent: 'center' }}>
        <b>{backendStatus}</b>
      </div>
    </div>

  );
}

export default ReleaseNotesList;