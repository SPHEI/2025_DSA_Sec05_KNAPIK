"use client"
import "../globals.css";
import { useState, useEffect } from "react";
import PaymentBox from "../components/PaymentBox";
import placeholder from "./placeholder.png"
import card from "./card.png"
import wrench from "./wrench.png"
import calendar from "./calendar.png"
import { useRouter, usePathname } from 'next/navigation';
//Page is shared by all types of accounts
function Dashboard() {
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
                    <div className="page-head w-[50%]">
                        <b className="text-4xl">Dashboard</b> 
                    </div>
                    <div className="white-box w-[50%] h-[200px]">
                        <div className="flex flex-row items-center justify-center relative right-20 gap-8">
                            <img src={placeholder.src} width={175} />
                            <div className="flex flex-col">
                                <b className="text-2xl">Tenant Info</b>
                                <h1 className="text-xl">Name: Jan Kowalski</h1>
                                <h1 className="text-xl">Apartmet: 15 Street</h1>
                            </div>
                        </div>
                    </div>
                    <div className="flex flex-row gap-4 w-[50%]">
                        <div className="white-box h-[150px] w-[100%]">
                            <div className="flex flex-col items-center justify-center">
                                <img src={card.src} width={40} />
                                <b>Total Rent Paid</b>
                                <h1>2400$</h1>
                            </div>
                        </div>
                        <div className="white-box h-[150px] w-[100%]">
                            <div className="flex flex-col items-center justify-center">
                                <img src={wrench.src} width={40} />
                                <b>Open Requests</b>
                                <h1>2</h1>
                            </div>
                        </div>
                        <div className="white-box h-[150px] w-[100%]">
                            <div className="flex flex-col items-center justify-center">
                                <img src={calendar.src} width={40} />
                                <b>Next Rent Due</b>
                                <h1>28.06.2025</h1>
                            </div>
                        </div>
                    </div>
                    <div className="white-box w-[50%] py-4">
                        <div className="flex flex-col items-left justify-start w-full h-full gap-2">
                            <b className="text-xl">Recent Operations</b>
                            <div className="flex flex-col gap-1">
                                <PaymentBox date={"April 30 2025"} type={"Loss"} amount={2500}/>
                                <PaymentBox date={"April 30 2025"} type={"Income"} amount={2500}/>
                                <PaymentBox date={"April 30 2025"} type={""} amount={2500}/>
                                <PaymentBox date={"April 30 2025"} type={""} amount={2500}/>
                            </div>
                        </div>
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
