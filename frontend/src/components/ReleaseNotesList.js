import './ReleaseNotesList.css';

import ReleaseNoteCard from './ReleaseNotesCard';
import React from "react";


function ReleaseNotesList({ projects }) {
  if (projects.length === 0) {
    return <p><b>No updates found.</b></p>;
  }
  let projectList = [];


  for (let project of projects) {
    if (project.ReleaseNotes == null) {
      continue
    } else if (project.ReleaseNotes.length === 0){
      project.ReleaseNotes = {
        "notusedname": {
          BodyHTML:"<b>RELEASES NOT FOUND</b>"
        }
      }
    }
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