"use client"
import "../globals.css";
import { useState, useEffect } from "react";


function Login() {
    const [ready,setReady] = useState(false)

    useEffect(() => {
        //Page setup goes here
        setReady(true);
    })

    if(ready)
    {
        return (
            <main>
                <h1>login page</h1>
            </main>
        );
    }
};

export default Login;
