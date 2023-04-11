import "./ModalWarn.css"

import React, { useState, useEffect } from "react";

function WarningModal() {
  const [showAlert, setShowAlert] = useState(true);
  const [message, setMessage] = useState("");

  useEffect(() => {
    async function fetchData() {
      try {
        const response = await fetch("http://localhost:8080/redis/status");
        const data = await response.json();
        if (data == "") {
          setShowAlert(false);
        } else {
          setShowAlert(true);
        }
        setMessage(data);

      } catch (error) {
        console.error(error);
      }
    }

    fetchData();
  }, []);

  return (
    <>
    {showAlert && (
      <div className="warning-window">
          <div className="warning-text">{message}</div>
      </div>
    )}
    </>
  );
}


export default React.memo(WarningModal);



