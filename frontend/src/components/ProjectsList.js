
import ProjectCard from './ProjectsCard';
import React, { useState } from "react";


function ProjectsList({ projects }) {

  const [inputValue, setInputValue] = useState("");

  const handleInputChange = (event) => {
    setInputValue(event.target.value);
  };

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

            {projects.map(project => <ProjectCard key={project.name} project={project} />)}



          </table>
        </div>
      </div>
    </div>
  );
}


export default ProjectsList;