import 'bootstrap/dist/css/bootstrap.min.css';

import ProjectsList from '../components/ProjectsList';

import React, { useState, useEffect } from "react";

function Projects() {
  const [allProjects, setAllProjects] = useState([]);

  const [resultData, setResultData] = useState([]);

  const [watcherFilterDisabled, setWatcherFilterDisabled] = useState(true);
  const [allFilterDisabled, setAllFilterDisabled] = useState(false);


  useEffect(() => {
    async function fetchData() {
      try {
        const response = await fetch(`http://localhost:8080/projects`);

        const data = await response.json();

        data.sort((a, b) => a.Project.Name.localeCompare(b.Project.Name));
        if (data) {
          setAllProjects(data);

          // TODO: refact, for sure exist better way
          // by default show watcher list projects
          let resultData = data.filter(item => item.IsInWatcherList === true)
          setResultData(resultData)
        }
      } catch (error) {
        console.error(error);
      }
    }

    fetchData();
  }, []);

  // TODO: should be by default
  function showWatcherClick() {
    setWatcherFilterDisabled(true)
    setAllFilterDisabled(false)
    let res = allProjects.filter(item => item.IsInWatcherList === true)
    setResultData(res)
  }

  function showAllClick() {
    setAllFilterDisabled(true)
    setWatcherFilterDisabled(false)
    setResultData(allProjects)
  }

  return (
    <div>
      {/* buttons menu */}
      <div className="project-list">
        <div className="container-sm mt-5 ml-5 mr-5">
          <div className="row justify-content-end">
            <div className="col-auto">
              <button
              disabled={ allFilterDisabled ? true : false }
              onClick={showAllClick}>show all</button>
            </div>
            <div className="col-auto">
              <button
                disabled={ watcherFilterDisabled ? true : false }
                onClick={showWatcherClick}
                >show watcher list</button>
            </div>
          </div>
        </div>
      </div>
{/* buttons menu end */}
      <ProjectsList projects={resultData} />
    </div>
  );
}

export default Projects;