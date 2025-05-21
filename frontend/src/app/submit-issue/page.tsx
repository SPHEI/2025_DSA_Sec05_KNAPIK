"use client";
import "../globals.css";
import { useState, useEffect } from "react";

//Tenants Only
function SubmitIssue() {
  const [ready, setReady] = useState(false);
  const [error, setError] = useState("none");
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");

  const handleSubmit = () => {
    //handle api

    console.log("Submitted:", { title, description });
    alert("Issue submitted!");
  };

  useEffect(() => {
    //Page setup goes here
    setReady(true);
  }, []);

  if (ready) {
    if (error == "none") {
      return (
        <main className="flex flex-col items-center justify-start gap-6 mt-10">
          <div className="page-head w-[50%]">
            <b className="text-4xl">Please describe your issue here</b>
          </div>

          <input
            type="text"
            placeholder="Title....."
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="input-box w-[50%]"
          />

          <textarea
            placeholder="Describe the issue in detail..."
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className="input-box w-[50%] h-[300px]"
          />

          <button onClick={handleSubmit} className="black-button">
            Submit Issue
          </button>
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

export default SubmitIssue;
