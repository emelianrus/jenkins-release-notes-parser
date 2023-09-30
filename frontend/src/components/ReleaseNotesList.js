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

        setBackendStatus("Loading");
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
        const uniqueData = data.filter(note => !releaseNotes.some(existingNote => existingNote.version === note.version));
        // Update releaseNotes based on previous state
        setReleaseNotes(prevReleaseNotes => [...prevReleaseNotes, ...uniqueData]);
        setBackendStatus("");
      } catch (error) {
        console.error(error);
      }
    }
  }

  let projectList = [];

  let resultNotes = [];
  for (let project of releaseNotes) {

    if (project === undefined || project.ReleaseNotes == null) {
      continue
    } else if (project.ReleaseNotes.length === 0){

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