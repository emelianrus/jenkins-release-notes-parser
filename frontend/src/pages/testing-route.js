
import React from "react";
import { useParams } from 'react-router-dom';


const TestingRoute = () => {
  const { owner, repo } = useParams();
  console.log(owner);
  console.log(repo);

  return (
    <div>
      <h1>
        TestingRoute PAGE CONTENT
      </h1>
    </div>
  );
};
  
export default TestingRoute;