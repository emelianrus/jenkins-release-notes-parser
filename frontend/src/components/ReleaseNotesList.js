import './ReleaseNotesList.css';

import ReleaseNoteCard from './ReleaseNotesCard';
import React from "react";


function ReleaseNotesList({ projects }) {


  let projectList = [];


  let releaseNotes = []

  for (let project of projects) {
    if (project === undefined  || project.ReleaseNotes == null) {
      continue
    } else if (project.ReleaseNotes.length === 0){

      releaseNotes.push({
        ReleaseNotes: {
          "notusedname": {
            BodyHTML:"<b>RELEASES NOT FOUND</b>"
          }

        }
      })
      continue
    }

    releaseNotes.push(project)
  }

  if (releaseNotes.length === 0) {
    return <p><b>No updates found.</b></p>;
  }

  for (let project of releaseNotes) {
    for (const [key, value] of Object.entries(project.ReleaseNotes)) {
      projectList.push(
        <ReleaseNoteCard key={project.Name + value.Name} project={value} projectName={project.Name}/>
      )
    }
  }

  return (
    <div className="project-list">
      <div className="clearfix container-xl px-3 px-md-4 px-lg-5 mt-4">
        {projectList}
      </div>
    </div>
  );
}

export default ReleaseNotesList;