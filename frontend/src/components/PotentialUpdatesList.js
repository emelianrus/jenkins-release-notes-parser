import './PotentialUpdatesList.css';

import PotentialUpdatesCard from './PotentialUpdatesCard';
import React from "react";


function PotentialUpdatesList({ projects }) {
  if (projects.length === 0) {
    return <p><b>No updates found.</b></p>;
  }
  let projectList = [];

  // const projectList = [];
  for (const project of projects) {
    console.log(project)
    // console.log(project.ReleaseNotes)
    if (project.ReleaseNotes === undefined || project.ReleaseNotes === null) {
      continue
    }
    for (const [key, value] of Object.entries(project.ReleaseNotes)) {

      projectList.push(
        <PotentialUpdatesCard key={project.Name + value.Name} project={value} projectName={project.Name}/>
      )
    }

  //   projectList = Object.entries(project).map(([projectName, project]) => (
  //     <PotentialUpdatesCard
  //       key={projectName}
  //       project={project}
  //       projectName={projectName}
  //     />
  //   ));

  }

  return (
    <div className="project-list">
      <div className="clearfix container-xl px-3 px-md-4 px-lg-5 mt-4">
        {projectList}
      </div>
    </div>
  );
}


// {Object.keys(projects).map((key) => (
//   <React.Fragment key={key}>
//     {projects[key].map((item) => (
//       <PotentialUpdatesCard key={item.Name} project={item} projectName={key}/>
//     ))}
//   </React.Fragment>
// ))}


export default PotentialUpdatesList;