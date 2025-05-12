"use client"
import "../globals.css";
import { useState, useEffect } from "react";
import { useRouter } from 'next/navigation'
import bg from './bg.png';

function Login() {
    const [ready,setReady] = useState(false)
    const [error, setError] = useState('none')

    const router = useRouter();

    useEffect(() => {
        //Page setup goes here
        setReady(true);
    },[])


    const inputStyle = "rounded-lg py-1 px-2 border-[1px] border-solid border-[var(--borders)]"
    
    if(ready)
    {
        if(error == 'none')
        {
            return (
                <main>
                    <div className="white-box w-[14%] py-4 relative top-25">
                        <div className="flex flex-col gap-6 py-1">
                            <div className="flex flex-col">
                                <b>Email</b>
                                <input className={inputStyle} placeholder="Email"></input>
                            </div>
                            <div className="flex flex-col">
                                <b>Password</b>
                                <input className={inputStyle} placeholder="Password"></input>
                            </div>
                            <button className="black-button" onClick={() => router.push("/dashboard")}>Sign In</button>
                        </div>
                    </div>
                    <img src={bg.src} className="absolute top-0 left-0 w-screen h-screen z-[-1]"/>
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

export default Login;
