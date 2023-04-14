import './PotentialUpdatesList.css';

import PotentialUpdatesCard from './PotentialUpdatesCard';
import React from "react";


function PotentialUpdatesList({ projects }) {
  if (typeof projects === 'object' && Object.keys(projects).length === 0) {
    return <p><b>No updates found.</b></p>;
  }
  return (
    <div className="project-list">

      <div className="clearfix container-xl px-3 px-md-4 px-lg-5 mt-4">
        {Object.keys(projects).map((key) => (
          <React.Fragment key={key}>
            {projects[key].map((item) => (
              <PotentialUpdatesCard key={item.Name} project={item} projectName={key}/>
            ))}
          </React.Fragment>
        ))}
      </div>
    </div>
  );
}


export default PotentialUpdatesList;