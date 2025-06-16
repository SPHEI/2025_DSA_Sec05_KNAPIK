"use client";
import "./globals.css";
import { useState, useEffect } from "react";
import bghp from "./assets/background.jpg";

//App start - Checks auth and redirects either to login or dashboard
function App() {
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
          <img
            src={bghp.src}
            alt="Background"
            className="absolute top-0 left-0 w-screen h-screen z-[-1]"
          />

          <div className="relative z-10 flex flex-col items-center justify-start pt-60 w-full text-center px-4">
            <h1 className="text-white text-7xl font-bold mb-4 drop-shadow-lg">
              The amazing building-management app
            </h1>
            <p className="text-white text-2xl drop-shadow-md max-w-3xl mx-auto">
              Everyday we bring humanity closer to the end of everyday with our products of today.
            </p>
          </div>
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

export default App;
