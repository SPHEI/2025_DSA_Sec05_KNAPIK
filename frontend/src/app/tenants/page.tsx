"use client"
import "../globals.css";
import { useState, useEffect } from "react";
import TenantBox from "../components/TenantBox";
//Admin only
function Dashboard() {
    const [ready,setReady] = useState(false)
    const [error, setError] = useState('none')
    const [names,setNames] = useState([''])

    useEffect(() => {
        const fetchData = async () => {
            try {
              const res = await fetch('http://localhost:8080/tenents');
              const data = await res.json();
              //alert(JSON.stringify(data));
              if(data.message)
              {
                setError(data.message)
              }
              else
              {
                setNames(data.names)
              }
            } catch (err: any) {
                setError(err.message)
            } finally{
                setReady(true);
            }
          };
          fetchData();
    },[])

    if(ready)
    {
        if(error == 'none')
        {
            return (
                <main>
                    <div className="page-head w-[50%]">
                        <b className="text-4xl">Tenants</b> 
                        <button className="black-button">+ Add Tenants</button>
                    </div>
                    {names.map((text, index) => <TenantBox key={index} name={text}/>)}
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

export default Dashboard;
