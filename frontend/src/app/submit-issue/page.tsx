"use client"
import "../globals.css";
import { useState, useEffect } from "react";

//Tenants Only
function SubmitIssue() {
    const [ready,setReady] = useState(false)
    const [error, setError] = useState('none')

    useEffect(() => {
        //Page setup goes here
        setReady(true);
    },[])

    if(ready)
    {
        if(error == 'none')
        {
            return (
                <main>
                    <h1>issue submit page</h1>
                </main>
            );
        }
        else
        {
            return (
                <main>
                    <b>An error has occured:</b>
                    <h1>{error}</h1>
                </main>
            );
        }
    }
    else
    {
        return(
            <main>
                <h1>Loading...</h1>
            </main>
        )
    }
};

export default SubmitIssue;
