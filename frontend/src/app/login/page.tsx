"use client"
import "../globals.css";
import { useState, useEffect } from "react";
import { useRouter, usePathname } from 'next/navigation';
import bg from './bg.png';
import Cookies from "js-cookie";

function Login() {
    const [ready,setReady] = useState(false)
    const [error, setError] = useState('none')

    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')

    const router = useRouter();

    const pathname = usePathname();
    useEffect(() => {
        //Page setup goes here
        setReady(true);
    },[pathname])

    async function sendLogin(){
        try {
            const res = await fetch('http://localhost:8080/login',{
                method:'POST',
                body: JSON.stringify({ 
                    email,
                    password
                })
            });
            const data = await res.json();
            if(data.message)
            {
                alert(data.message)
            }
            else
            {
                Cookies.set("token", data.token)
                Cookies.set("role", data.role)
                router.push("/dashboard")
            }
        } catch (err: any) {
            setError(err.message)
        } finally{
            setReady(true);
        }
    }


    const inputStyle = "rounded-lg py-1 px-2 border-[1px] border-solid border-[var(--borders)]"
    
    if(ready)
    {
        if(error == 'none')
        {
            return (
                <main>
                    <div className="white-box w-[16%] py-4 relative top-50">
                        <div className="flex flex-col gap-6 py-1">
                            <div className="flex flex-col">
                                <b>Email</b>
                                <input className={inputStyle} placeholder="Email" onChange={(a)=>{setEmail(a.target.value)}}></input>
                            </div>
                            <div className="flex flex-col">
                                <b>Password</b>
                                <input className={inputStyle} type="password" placeholder="Password" onChange={(a)=>{setPassword(a.target.value)}}></input>
                            </div>
                            <button className="black-button" onClick={sendLogin}>Sign In</button>
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
