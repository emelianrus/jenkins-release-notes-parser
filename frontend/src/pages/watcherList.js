
import React from "react";

import JsonEditor from "../components/JsonEditor";


const WatcherList = () => {
  const data = {
    "ace-editor": "1.1",
    "ansicolor": "1.0.2",
    "ant": "481.v7b_09e538fcca",
    "antisamy-markup-formatter": "2.7",
  };

  return (
    <div className="editor">
      <JsonEditor data={data} />
    </div>
  );
};

export default WatcherList;