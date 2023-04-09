import './ReleaseNotesList.css';

import ReleaseNoteCard from './ReleaseNoteCard';
import React, { useState } from "react";


function ReleaseNotesList({ projects, projectRepo, projectOwner}) {
  const [inputValue, setInputValue] = useState("");

  const handleInputChange = (event) => {
    setInputValue(event.target.value);
  };

  return (
    <div className="project-list">
      <div id="repository-container-header" className="pt-3 hide-full-screen">
        <div className="d-flex flex-wrap flex-justify-end px-md-4 px-lg-5">
          <div className="d-flex flex-wrap flex-items-center wb-break-word f3 text-normal">
            {projectOwner}
            <span className="mx-1 flex-self-stretch color-fg-muted">/</span>
            {projectRepo}
          </div>
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

        {projects.map(project => <ReleaseNoteCard key={project.name} project={project} />)}
      </div>
    </div>
  );
}


export default ReleaseNotesList;