import './ProjectList.css';

import ProjectCard from './ProjectCard';
import React, { useState } from "react";


function ProjectList({ projects }) {

  const [inputValue, setInputValue] = useState("");

  const handleInputChange = (event) => {
    setInputValue(event.target.value);
  };

  return (
    <div className="project-list">
      <div id="repository-container-header" className="pt-3 hide-full-screen">
        <div class="d-flex flex-wrap flex-justify-end px-md-4 px-lg-5">
          <span>
            REPO_GROUP/PROJECT_NAME
          </span>
        </div>
      </div>

      <div className="clearfix container-xl px-3 px-md-4 px-lg-5 mt-4">
        <div className="d-flex flex-justify-center">
          <div className="d-flex flex-column flex-sm-row flex-wrap mb-3 pb-3 col-11 justify-content-sm-end border-bottom">
            <div className="d-flex flex-column flex-md-row">
              <div>
                <form className="position-relative ml-md-2" action="/#" acceptCharset="UTF-8" method="get">
                  <input id="release-filter" type="search" name="q" className="form-control subnav-search-input width-full"
                    placeholder="Find a release" value={inputValue} onChange={handleInputChange} aria-label="Find a release"/>
                </form>
              </div>
            </div>
          </div>
        </div>
        {/* top menu releases find edit ^^ */}

        {projects.map(project => <ProjectCard key={project.name} project={project} />)}
      </div>
    </div>
  );
}


export default ProjectList;