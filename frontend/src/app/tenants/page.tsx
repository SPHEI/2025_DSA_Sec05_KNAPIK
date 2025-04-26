"use client"
import "../globals.css";
import { useState, useEffect } from "react";

//Admin only
function Dashboard() {
    const [ready,setReady] = useState(false)

    useEffect(() => {
        //Page setup goes here
        setReady(true);
    })

    if(ready)
    {
        return (
            <main>
                <h1>tenants page</h1>
            </main>
        );
    }
};

export default Dashboard;
