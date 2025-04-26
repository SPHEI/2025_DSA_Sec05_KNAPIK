"use client"
import "../globals.css";
import { useState, useEffect } from "react";

//Tenants Only
function SubmitIssue() {
    const [ready,setReady] = useState(false)

    useEffect(() => {
        //Page setup goes here
        setReady(true);
    })

    if(ready)
    {
        return (
            <main>
                <h1>issue submit page</h1>
            </main>
        );
    }
};

export default SubmitIssue;
