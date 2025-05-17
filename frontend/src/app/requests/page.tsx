"use client";
import "../globals.css";
import { useState, useEffect } from "react";
import RequestBox from "../components/RequestBox";

//Page is shared by all types of accounts
function Requests() {
  const [ready, setReady] = useState(false);
  const [error, setError] = useState("none");
  useEffect(() => {
    //Page setup goes here
    setReady(true);
  }, []);

  if (ready) {
    if (error == "none") {
      return (
        <main>
          <div className="page-head w-[50%]">
            <b className="text-4xl">My requests</b>
          </div>
          <RequestBox
            title="Water leak"
            date="10 May 2025"
            status="Solved"
            colorClass="status-box-green"
          />
          <RequestBox
            title="Water leak"
            date="10 May 2025"
            status="Solved"
            colorClass="status-box-red"
          />
          <RequestBox
            title="Water leak"
            date="10 May 2025"
            status="Solved"
            colorClass="status-box-yellow"
          />
        </main>
      );
    } else {
      return (
        <main>
          <b>An error has occured:</b>
          <h1>{error}</h1>
        </main>
      );
    }
  } else {
    return (
      <main>
        <h1>Loading...</h1>
      </main>
    );
  }
}

export default Requests;
