import './ReleaseNotesList.css';

import ReleaseNoteCard from './ReleaseNotesCard';
import React from "react";


function ReleaseNotesList({ projects }) {
  if (projects.length === 0) {
    return <p><b>No updates found.</b></p>;
  }
  let projectList = [];


  for (const project of projects) {
    if (project.ReleaseNotes === undefined || project.ReleaseNotes === null) {
      continue
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