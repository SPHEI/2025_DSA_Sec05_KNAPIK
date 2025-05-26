"use client"
import "../globals.css";
import { useState, useEffect } from "react";
import { useRouter, usePathname } from 'next/navigation';

//Admin only
function Reports() {
    const [ready,setReady] = useState(false)
    const [error, setError] = useState('none')
    const pathname = usePathname();
    useEffect(() => {
        //Page setup goes here
        setReady(true);
    },[pathname])

    if(ready)
    {
        if(error == 'none')
        {
            return (
                <main>
                    <h1>reports page</h1>
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

export default Reports;
