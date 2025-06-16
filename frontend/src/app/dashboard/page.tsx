"use client"
import "../globals.css";
import { useState, useEffect } from "react";
import PaymentBox from "../components/PaymentBox";
import placeholder from "./placeholder.png"
import card from "./card.png"
import wrench from "./wrench.png"
import calendar from "./calendar.png"
import { useRouter, usePathname } from 'next/navigation';
import Cookies from "js-cookie";
//Page is shared by all types of accounts
function Dashboard() {
    const [ready,setReady] = useState(false)
    const [error, setError] = useState('none')

    const [userInfo, setUserInfo] = useState({id: -1, name: '', phone: '', role_id: -1})
    const [tenantInfo, setTenantInfo] = useState({apartment_id: -1, rent: -1})
    const [subconInfo, setSubconInfo] = useState({address: '', nip: '', speciality_id: -1})

    const [requests, setRequests] = useState([{id: -1, title : '',description: '', date_reported: '', status_id: -1, apartment_id: -1, name: ''}])
    const [repairs, setRepairs] = useState([{id: -1, title: '', fault_report_id: -1, date_assigned: '', date_completed: '', status_id: -1, subcontractor_id: -1, name: ''}])

    const router = useRouter();
    const pathname = usePathname();
    useEffect(() => {
        refresh()
    },[pathname])
    
    async function refresh()
    {
        var t = Cookies.get("token");
        try{
            const res = await fetch('http://localhost:8080/info?token=' + t)
            const data = await res.json();
            if(data.message)
            {
                setError(data.message)
            }
            else
            {
                setUserInfo(data)
                if(data.role_id != 3)
                {
                    const res3 = await fetch('http://localhost:8080/faults/list?token=' + t)
                    const data3 = await res3.json();
                    if(data3.message)
                    {
                        setError(data3.message)
                    }
                    else
                    {
                        //alert(JSON.stringify(data))
                        setRequests(data3);
                    }
                }
                else
                {
                    const res4 = await fetch('http://localhost:8080/repair/list?token=' + t)
                    const data4 = await res4.json();
                    if(data4.message)
                    {
                        setError(data4.message)
                    }
                    else
                    {
                        //alert(JSON.stringify(data))
                        setRepairs(data4);
                    }
                }
                if(data.role_id == 2)
                {
                    const res2 = await fetch('http://localhost:8080/tenant/info?token=' + t)
                    const data2 = await res2.json();
                    if(data2.message)
                    {
                        setError(data2.message)
                    }
                    else
                    {
                        setTenantInfo(data2)
                    }
                }
                else if(data.role_id == 3)
                {
                    const res2 = await fetch('http://localhost:8080/subcon/info?token=' + t)
                    const data2 = await res2.json();
                    if(data2.message)
                    {
                        setError(data2.message)
                    }
                    else
                    {
                        setSubconInfo(data2)
                    }
                }
            }
        }
        catch(err: any)
        {
            setError(err.message);
        }
        setReady(true);
    }
    if(ready)
    {
        if(error == 'none')
        {
            return (
                <main>
                    <div className="page-head w-[50%]">
                        <b className="text-4xl">Dashboard</b> 
                    </div>
                    <div className="white-box w-[50%] h-[200px] min-w-[600px]">
                        <div className="flex flex-row items-center justify-center relative right-20 gap-8">
                            <img src={placeholder.src} width={175} />
                            {userInfo.role_id === 1 &&
                                <div className="flex flex-col">
                                    <b className="text-2xl">Admin Info</b>
                                    <h1 className="text-xl">Name: {userInfo.name}</h1>
                                </div>
                            }
                            {userInfo.role_id === 2 &&
                                <div className="flex flex-col">
                                    <b className="text-2xl">Tenant Info</b>
                                    <h1 className="text-xl">Name: {userInfo.name}</h1>
                                    <h1 className="text-xl">Apartment: {tenantInfo.apartment_id}</h1>
                                </div>
                            }
                            {userInfo.role_id === 3 &&
                                <div className="flex flex-col">
                                    <b className="text-2xl">Subcontractor Info</b>
                                    <h1 className="text-xl">Name: {userInfo.name}</h1>
                                    <h1 className="text-xl">Address: {subconInfo.address}</h1>
                                    <h1 className="text-xl">Nip: {subconInfo.nip}</h1>
                                </div>
                            }
                        </div>
                    </div>
                    <div className="flex flex-row gap-4 w-[50%] min-w-[600px]">
                        {userInfo.role_id != 3 &&
                        <div className="white-box h-[150px] w-[100%]">
                            <div className="flex flex-col items-center justify-center">
                                <img src={card.src} width={40} />
                                <b>Total Rent Paid</b>
                                <h1>2400$</h1>
                            </div>
                        </div>
                        }
                        <button className="cursor-pointer h-[100%] w-[150%]" onClick={() => router.push("/requests")}>
                            <div className="white-box h-[150px] w-[100%]">
                                <div className="flex flex-col items-center justify-center">
                                        <img src={wrench.src} width={40} />
                                        <b>{userInfo.role_id === 1 ? "Requests" : userInfo.role_id === 2 ? "My Requests" : "Assigned Repairs"}</b>
                                        <h1>{userInfo.role_id != 3 ? requests.length : repairs.length}</h1>
                                </div>
                            </div>
                        </button>
                        {userInfo.role_id != 3 &&
                            <div className="white-box h-[150px] w-[100%]">
                                <div className="flex flex-col items-center justify-center">
                                    <img src={calendar.src} width={40} />
                                    <b>Next Rent Due</b>
                                    <h1>28.06.2025</h1>
                                </div>
                            </div>
                        }
                    </div>
                    {userInfo.role_id != 3 &&
                        <div className="white-box w-[50%] min-w-[600px] py-4">
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
                    }
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
