import "./apiStats.css"
import React, { useState, useEffect } from "react";

// TODO: implement this
// should be useful to track api status just for information
function APIStats() {
  const [apiStats, setApiStats] = useState([]);

  useEffect(() => {
    async function fetchData() {
      try {
        // TODO: https://api.github.com/repos/OWNER/REPO/releases
        // repos/:owner/:repo/releases
        // /project/:owner/:repo/releases
        const response = await fetch(`http://localhost:8080/api/stats`);
        const data = await response.json();

        // pass as new single object instead of several params
        setApiStats(data);
      } catch (error) {
        console.error(error);
      }
    }

    fetchData();
  }, []);

  return (
    <div className="container">
      <div className="row justify-content-center">
        <div className="col-md">
          <div className="stats-container">
            <div className="stats-column">
              <h2>github stats</h2>
              <div className="stats-row">
                <div className="stats-key">RateLimit</div>
                <div className="stats-value">11</div>
              </div>
              <div className="stats-row">
                <div className="stats-key">RateLimitRemaning</div>
                <div className="stats-value">11</div>
              </div>
              <div className="stats-row">
                <div className="stats-key">RateLimitReset</div>
                <div className="stats-value">1234</div>
              </div>
              <div className="stats-row">
                <div className="stats-key">RateLimitUsed</div>
                <div className="stats-value">1234</div>
              </div>
              <div className="stats-row">
                <div className="stats-key">WaitSlotSeconds</div>
                <div className="stats-value">1234</div>
              </div>
            </div>
            <div className="stats-column">
              <h2>plugin center stats</h2>
              <div className="stats-row">
                <div className="stats-key">last downloaded</div>
                <div className="stats-value">123</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default APIStats;