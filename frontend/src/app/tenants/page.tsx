"use client"
import "../globals.css";
import { useState, useEffect } from "react";
import TenantBox from "../components/TenantBox";
import Cookies from "js-cookie";
import { useRouter, usePathname } from 'next/navigation';

//Admin only
function Tenants() {
    const [ready,setReady] = useState(false)
    const [error, setError] = useState('none')
    const [names,setNames] = useState([{id: -1, name: '', email: '', phone: '', role_id: -1, id_2: -1, name_2: '', price: '', renting_id: -1, status: ''}])
    const [apartaments,setApartaments] = useState([{id: -1,name: '', street: '', building_number: '', building_name: '',flat_number:'',owner_id:-1 }])
    const pathname = usePathname();
    const router = useRouter();
    useEffect(() => {
        refresh()
    },[pathname])

    const [sort, setSort] = useState('None')
    async function refresh()
    {
        try {
            var s = Cookies.get("tSort");
            if(s != null)
            {
                setSort(s)
                Cookies.remove("tSort")
            }
            var t = Cookies.get("token");
            await fetch('http://localhost:8080/test')
            await fetch('http://localhost:8080/payments/list?token=' + t)
            const res2 = await fetch('http://localhost:8080/apartament/list?token=' + t)
            const data2 = await res2.json();
            if(data2.message)
            {
                setError(data2.message)
            }
            else
            {
                setApartaments(data2);
            }
            const res = await fetch('http://localhost:8080/tenant/list?token=' + t)
            const data = await res.json();
            console.log(JSON.stringify(data));
            if(data.message)
            {
            setError(data.message)
            }
            else
            {
            for(const a of data)
            {
                //alert(a.id)
                const res2 = await fetch('http://localhost:8080/tenant/info?token=' + t +"&id=" + Number(a.id))
                const data2 = await res2.json()
                console.log(JSON.stringify(data2))
                if(data2!= null)
                {
                    a.id_2 = data2.apartment.id
                    a.name_2 = data2.apartment.name
                    a.price = data2.rent
                    a.renting_id = data2.renting_id
                    a.status = data2.status
                }
            }
            //alert(JSON.stringify(data))
            setNames(data)
            }
        } catch (err: any) {
            setError(err.message)
        } finally{
            setReady(true);
        }
    }

    const changeRent = async (id: number, newRent: number) => 
    {
        console.log(id + " " + newRent)
        var t = Cookies.get("token");
        try {
            const res = await fetch('http://localhost:8080/changerent',{
                method:'POST',
                body: JSON.stringify({ 
                    "token" : t,
                    "rent":{
                    "apartment_id" : id,
                    "price": newRent
                    }
                })
            });
            if(res.ok)
            {
                console.log("Rent changed succesfully.");
            }
            else
            {
                var data = await res.json()
                console.log(data.message)
            }
        } catch (err: any) {
            console.log(err.message)
        } finally{
            refresh()
        }
    }

    const changeApartment = async(a_id: number, u_id: number, date: string, r_id: number) =>
    {
        console.log(a_id + " " + u_id + " " + date)
        evict(r_id, date)
        var t = Cookies.get("token");
        const res2 = await fetch('http://localhost:8080/renting/start',{
            method:'POST',
            body: JSON.stringify({ 
                "token": t,
                "renting" : {
                    "apartment_id" : a_id,
                    "user_id" : u_id,
                    "start_date" : date + "T00:00:00Z"
                }
            })
        });
        if(res2.ok)
        {
            console.log("renting start succesful")
            refresh()
        }
        else
        {
            var data2 = await res2.json()
            alert(data2.message)
        }
    }

    async function evict(id: number, date: string)
    {
        console.log(id + " " + date)
        var t = Cookies.get("token");
        const res = await fetch('http://localhost:8080/renting/end',{
            method:'POST',
            body: JSON.stringify({ 
                "token": t,
                "end" : {
                    "end_date" : date + "T00:00:00Z",
                    "id" : id,
                }
            })
        });
        if(res.ok)
        {
            console.log("evict date succesful")
        }
        else
        {
            var data = await res.json()
            alert(data.message)
        }

        const res2 = await fetch('http://localhost:8080/renting/endStatus',{
            method:'POST',
            body: JSON.stringify({ 
                "token": t,
                "renting_id" :id 
            })
        });
        if(res2.ok)
        {
            console.log("evict status succesful")
        }
        else
        {
            var data2 = await res2.json()
            alert(data2.message)
        }
        refresh()
    }

    async function viewPayments(name: string)
    {
        Cookies.set("tFilter", name)
        router.push("/rent-and-payment")
    }
    
    const priorityPaid: Record<string, number> = {
        "Paid" : 3,
        "Pending" : 2,
        "Overdue" : 1
    };

    const priorityPending: Record<string, number> = {
        "Paid" : 2,
        "Pending" : 3,
        "Overdue" : 1
    };
    
    const priorityOverdue: Record<string, number> = {
        "Paid" : 1,
        "Pending" : 2,
        "Overdue" : 3
    };

    if(ready)
    {
        if(error == 'none')
        {
            return (
                <main>
                    <div className="page-head w-[50%] min-w-[600px]">
                        <b className="text-4xl">Tenants</b> 
                        <div className="flex flex-row gap-1 min-w-[300px]">
                            <h1 className="text-2xl">Sort:</h1>
                            <select className="input-box w-[50%]" value={sort} onChange={(a) => {setSort(a.target.value)}}>
                                <option value="None">None</option>
                                <option value="Paid">Paid</option>
                                <option value="Pending">Pending</option>
                                <option value="Overdue">Overdue</option>
                            </select>
                        <button className="black-button" onClick={() =>{router.push("/accounts")}}>+ Add Tenants</button>
                        </div>
                    </div>
                    {names.length > 0  ? (sort == 'None' ? names : [...names].sort((a, b) => {
                        const priorityA = sort == "Paid" ? priorityPaid[a.status] : sort == "Pending" ? priorityPending[a.status] : priorityOverdue[a.status];
                        const priorityB = sort == "Paid" ? priorityPaid[b.status] : sort == "Pending" ? priorityPending[b.status] : priorityOverdue[b.status];
                        return priorityB - priorityA;
                    })).map((text, index) => <TenantBox key={index} id={text.id} name={text.name} email={text.email} phone={text.phone} role_id={text.role_id} 
                    apartment_id = {text.id_2} apartment={text.name_2} rent={text.price} status={text.status} renting_id={text.renting_id}
                    evict={evict} changeRent={changeRent} changeApartment={changeApartment} viewPayments={viewPayments}
                    apartments={apartaments}/>)
                    :<h1>No tenants</h1>}
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

export default Tenants;
