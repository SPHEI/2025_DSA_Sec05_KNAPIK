"use client"
import "../globals.css";
import { useState, useEffect } from "react";

//Admin only
function Reports() {
    const [ready,setReady] = useState(false)

    useEffect(() => {
        //Page setup goes here
        setReady(true);
    })

    if(ready)
    {
        return (
            <main>
                <h1>reports page</h1>
            </main>
        );
    }
};

export default Reports;
