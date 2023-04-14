import "./apiStats.css"
import React from "react";

// TODO: implement this
// should be useful to track api status just for information
function APIStats() {


  return (
    <div class="container">
      <div class="row justify-content-center">
        <div class="col-md">
          <div class="stats-container">
            <div class="stats-column">
              <h2>github stats</h2>
              <div class="stats-row">
                <div class="stats-key">api calls used</div>
                <div class="stats-value">11</div>
              </div>
              <div class="stats-row">
                <div class="stats-key">api calls total</div>
                <div class="stats-value">11</div>
              </div>
              <div class="stats-row">
                <div class="stats-key">api reset time</div>
                <div class="stats-value">1234</div>
              </div>
            </div>
            <div class="stats-column">
              <h2>plugin center stats</h2>
              <div class="stats-row">
                <div class="stats-key">last downloaded</div>
                <div class="stats-value">123</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default APIStats;