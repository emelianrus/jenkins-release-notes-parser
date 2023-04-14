
import ProjectCard from './ProjectsCard';
import React from "react";


function ProjectsList({ projects }) {

  return (
    <div className="project-list">
      <div className="container-sm mt-5 ml-5 mr-5">
        <h3>Title</h3>
        <div className="table-responsive">
          <table className="table">
            <thead className="thead-light">
              <tr>
                <th>Project name</th>
                <th>Download status</th>
                <th>Force Download</th>
              </tr>
            </thead>

            {projects === undefined
              ? <tbody><tr><td>No projects to display</td></tr></tbody>
              : projects.map(project => <ProjectCard key={project.Name} project={project} />)
            }
          </table>
        </div>
      </div>
    </div>
  );
}


export default ProjectsList;