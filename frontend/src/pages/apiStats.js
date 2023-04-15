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
                <div class="stats-key">RateLimit</div>
                <div class="stats-value">11</div>
              </div>
              <div class="stats-row">
                <div class="stats-key">RateLimitRemaning</div>
                <div class="stats-value">11</div>
              </div>
              <div class="stats-row">
                <div class="stats-key">RateLimitReset</div>
                <div class="stats-value">1234</div>
              </div>
              <div class="stats-row">
                <div class="stats-key">RateLimitUsed</div>
                <div class="stats-value">1234</div>
              </div>
              <div class="stats-row">
                <div class="stats-key">WaitSlotSeconds</div>
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