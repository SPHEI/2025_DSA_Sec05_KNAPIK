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
    const [names,setNames] = useState([{id: -1, name: '', email: '', phone: '', role_id: -1, id_2: -1, name_2: '', price: '', status: ''}])
    const [apartaments,setApartaments] = useState([{id: -1,name: '', street: '', building_number: '', building_name: '',flat_number:'',owner_id:-1 }])
    const pathname = usePathname();
    const router = useRouter();
    useEffect(() => {
        refresh()
    },[pathname])

    async function refresh()
    {
        try {
            var t = Cookies.get("token");
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
            //alert(JSON.stringify(data));
            if(data.message)
            {
            setError(data.message)
            }
            else
            {
            for(const a of data)
            {
                //alert(a.id)
                const res2 = await fetch('http://localhost:8080/tenant/info?token=' + t +"&id=" + a.id)
                const data2 = await res2.json()
                //alert(JSON.stringify(data2))
                if(data2.status)
                {
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

    const changeApartment = async(a_id: number, u_id: number, date: string) =>
    {
        console.log(a_id + " " + u_id + " " + date)
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

    async function evict(id: number)
    {
        alert(id)
    }

    if(ready)
    {
        if(error == 'none')
        {
            return (
                <main>
                    <div className="page-head w-[50%] min-w-[600px]">
                        <b className="text-4xl">Tenants</b> 
                        <button className="black-button" onClick={() =>{router.push("/accounts")}}>+ Add Tenants</button>
                    </div>
                    {names.map((text, index) => <TenantBox key={index} id={text.id} name={text.name} email={text.email} phone={text.phone} role_id={text.role_id} 
                    apartment_id = {text.id_2} apartment={text.name_2} rent={text.price} status="paid"
                    evict={evict} changeRent={changeRent} changeApartment={changeApartment}
                    apartments={apartaments}/>)}
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
