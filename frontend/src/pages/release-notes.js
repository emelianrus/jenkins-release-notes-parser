import './releases.css';
import 'bootstrap/dist/css/bootstrap.min.css';

import ReleaseNotesList from '../components/ReleaseNotesList';

import React, { useState, useEffect } from "react";
import { useParams } from 'react-router-dom';

function Releases() {
  const { owner, repo } = useParams();

  const [projects, setProjects] = useState([]);
  const [projectRepo, setProjectRepo] = useState("");
  const [projectOwner, setProjectOwner] = useState("");

  useEffect(() => {
    async function fetchData() {
      try {
        // TODO: https://api.github.com/repos/OWNER/REPO/releases
        // repos/:owner/:repo/releases
        // /project/:owner/:repo/releases
        const response = await fetch(`http://localhost:8080/project/${owner}/${repo}/releases`);
        const data = await response.json();

        // pass as new single object instead of several params
        setProjects(data.ReleaseNotes);
        setProjectRepo(repo);
        setProjectOwner(owner);

      } catch (error) {
        console.error(error);
      }
    }

    fetchData();
  }, []);

  return (
    <div>
      <ReleaseNotesList projects={projects} projectRepo={projectRepo} projectOwner={projectOwner} />
    </div>
  );
}

export default Releases;