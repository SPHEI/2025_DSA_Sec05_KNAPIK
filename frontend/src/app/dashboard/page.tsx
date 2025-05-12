"use client"
import "../globals.css";
import { useState, useEffect } from "react";

//Page is shared by all types of accounts
function Dashboard() {
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
                    <div className="page-head w-[50%]">
                        <b className="text-4xl">Dashboard</b> 
                    </div>
                    <div className="white-box w-[50%] h-[200px]">

                    </div>
                    <div className="flex flex-row gap-4 w-[50%]">
                        <div className="white-box h-[200px] w-[100%]">
                            <div className="flex flex-col">
                                <i>image goes here</i>
                                <b>Total Rent Paid</b>
                                <h1>2400$</h1>
                            </div>
                        </div>
                        <div className="white-box h-[200px] w-[100%]">
                            <div className="flex flex-col">
                                <i>image goes here</i>
                                <b>Open Requests</b>
                                <h1>2</h1>
                            </div>
                        </div>
                        <div className="white-box h-[200px] w-[100%]">
                            <div className="flex flex-col">
                                <i>image goes here</i>
                                <b>Next Rent Due</b>
                                <h1>28.06.2025</h1>
                            </div>
                        </div>
                    </div>
                    <div className="white-box w-[50%] h-[200px]">

                    </div>
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
