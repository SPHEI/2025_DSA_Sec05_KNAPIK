"use client"
import "../globals.css";
import { useState, useEffect } from "react";


//Page is shared by all types of accounts
function Requests() {
    const [ready,setReady] = useState(false)

    useEffect(() => {
        //Page setup goes here
        setReady(true);
    })

    if(ready)
    {
        return (
            <main>
                <h1>requests page</h1>
            </main>
        );
    }
};

export default Requests;
