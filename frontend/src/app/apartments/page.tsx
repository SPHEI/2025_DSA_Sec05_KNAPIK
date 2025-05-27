"use client"
import "../globals.css";
import { useState, useEffect } from "react";
import Cookies from "js-cookie";
import ApartmentBox from "../components/ApartmentBox";
import { useRouter, usePathname } from 'next/navigation';

//Admin only
function Apartments() {
    const [ready,setReady] = useState(false)
    const [error, setError] = useState('none')

    const [apartaments,setApartaments] = useState([{id: -1,name: '', street: '', building_number: '', building_name: '',flat_number:'',owner_id:-1, rent: -1 }])
    const [owners, setOwners] = useState([{id: -1, name:'',email:'',phone:''}])

    const pathname = usePathname();

    const [showPopup,setShowPopup] = useState(false)

    const [showPopup2,setShowPopup2] = useState(false)

    const [name, setName] = useState('');
    const [street, setStreet] = useState('');
    const [number, setNumber] = useState('');
    const [bname, setBname] = useState('');
    const [fnumber,setFnumber] = useState('');
    const [ownerId, setOwnerId] = useState(1);

    const formatNumber = (value: string) => {
        const digits = value.replace(/\D/g, '');
        return digits;
    }
    function numberChange(e: React.ChangeEvent<HTMLInputElement>){
        const a = formatNumber(e.target.value);
        setNumber(a);
    }

    const [newOwnerName,setNewOwnerName] = useState('')
    const [newOwnerEmail,setNewOwnerEmail] = useState('')
    const [newOwnerPhone,setNewOwnerPhone] = useState('')

    function formatPhone(value: string) {
        let cleaned = value.replace(/[^\d+\s]/g, '');
        if (cleaned.startsWith('+')) {
            cleaned = '+' + cleaned.slice(1).replace(/[+]/g, '');
        } else {
            cleaned = cleaned.replace(/[+]/g, '');
        }
        return cleaned;
    };
    function phoneChange(e: React.ChangeEvent<HTMLInputElement>){
        const formatted = formatPhone(e.target.value);
        setNewOwnerPhone(formatted);
    };


    useEffect(() => {
        refresh()
    },[pathname])
    
    async function refresh()
    {
        var t = Cookies.get("token");
        try{
            const res = await fetch('http://localhost:8080/apartament/list?token=' + t)
            const data = await res.json();
            //alert(JSON.stringify(data))
            if(data.message)
            {
                setError(data.message)
            }
            else
            {
                setApartaments(data.apartaments);
            }
            const res2 = await fetch('http://localhost:8080/owner/list?token=' + t)
            const data2 = await res2.json();
            //alert(JSON.stringify(data2))
            if(data2.message)
            {
                setError(data2.message)
            }
            else
            {
                setOwners(data2.owners);
            }
        }catch(err: any){
            setError(err.message);
        }finally{
            setReady(true);
        }
    }

    async function addApartment()
    {
        var t = Cookies.get("token");
    //    alert(JSON.stringify({ 
    //                 "token" : t,
    //                 name,
    //                 street,
    //                 "building_number" : number,
    //                 "building_name" : bname ,
    //                 "flat_number" : fnumber,
    //                 "owner_id": ownerId
    //             }))
        try {
            const res = await fetch('http://localhost:8080/apartament/add',{
                method:'POST',
                body: JSON.stringify({ 
                    "token" : t,
                    name,
                    street,
                    "building_number" : number,
                    "building_name" : bname,
                    "flat_number" : fnumber,
                    "owner_id": ownerId
                })
            });
            if(res.ok)
            {
                alert("Apartment added succesfully.");
            }
            else
            {
                var data = await res.json()
                alert(data.message)
            }
            setShowPopup(false)
            refresh()
        } catch (err: any) {
            setError(err.message)
        } finally{
            setReady(true);
        }
    }

    async function addOwner()
    {
        var t = Cookies.get("token");
        try {
            const res = await fetch('http://localhost:8080/owner/add',{
                method:'POST',
                body: JSON.stringify({ 
                    "token" : t,
                    "name": newOwnerName,
                    "email": newOwnerEmail,
                    "phone": newOwnerPhone
                })
            });
            if(res.ok)
            {
                alert("Owner added succesfully.");
            }
            else
            {
                var data = await res.json()
                alert(data.message)
            }
            setShowPopup2(false)
            refresh()
        } catch (err: any) {
            setError(err.message)
        } finally{
            setReady(true);
        }
    }

    const line = "flex flex-row gap-1";

    if(ready)
    {
        if(error == 'none')
        {
            return (
                <main>
                    <div className="page-head w-[50%]">
                        <b className="text-4xl">Apartments</b> 
                        <button className="black-button" onClick={() => {setShowPopup(true)}}>+ Add Apartment</button>
                    </div>
                    {apartaments.map((a,index) => <ApartmentBox key={index} id={a.id} name={a.name} street={a.street} building_number={a.building_number} building_name={a.building_name} flat_number={a.flat_number} owner_id={a.owner_id} rent={a.rent} refresh={refresh}/>)}
                    {showPopup && (
                    <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/50">
                        <div className="white-box w-[40%] py-4 rounded-lg relative">
                            <div className="flex flex-col gap-2 w-[100%]">
                                <b className="text-4xl">Add Apartment</b>
                                <div className={line}>
                                    <b className="w-[34%]">Name</b>
                                    <b className="w-[34%]">Owner</b>
                                </div>
                                <div className={line}>
                                    <input className="input-box w-[34%]" placeholder="Name" value={name} onChange={(a) => {setName(a.target.value)}}/>
                                    <select className="input-box w-[34%]" value={ownerId}onChange={(a) => {setOwnerId(Number(a.target.value))}}>
                                        {owners.map((a,index) => (<option key={index} value={a.id}>{a.email}</option>))}
                                    </select>
                                    <button className="black-button w-[34%]" onClick={() => {setShowPopup2(true)}}>Add Owner</button>
                                </div>
                                <div className={line}>
                                    <b className="w-[34%]">Street</b>
                                    <b>Building Number</b>
                                </div>
                                <div className={line}>
                                    <input className="input-box w-[34%]" placeholder="Street" value={street} onChange={(a) => {setStreet(a.target.value)}} />
                                    <input className="input-box w-[34%]" placeholder="Number" value={number} onChange={numberChange}/>
                                </div>
                                <div className={line}>
                                    <b className="w-[34%]">Building Name</b>
                                    <b>Flat Number</b>
                                </div>
                                <div className={line}>
                                    <input className="input-box w-[34%]" placeholder="Building name" value={bname} onChange={(a) => {setBname(a.target.value)}} />
                                    <input className="input-box w-[34%]" placeholder="Flat number" value={fnumber} onChange={(a) => {setFnumber(a.target.value)}} />
                                </div>
                                <button className="black-button w-[34%]" onClick={addApartment}>Add</button>
                            </div>
                            <button onClick={() => setShowPopup(false)}className="absolute top-4 right-4 text-xl font-bold cursor-pointer">x</button>
                        </div>
                    </div>
                    )}
                    {showPopup2 && (
                    <div className="fixed inset-0 z-31 flex items-center justify-center bg-black/50">
                        <div className="white-box w-[40%] py-4 rounded-lg relative">
                            <div className="flex flex-col gap-2 w-[100%]">
                                <b className="text-4xl">Add Owner</b>
                                    <div className={line}>
                                        <b className="w-[34%]">Name</b>
                                        <b className="w-[34%]">Email</b>
                                        <b className="w-[34%]">Phone</b>
                                    </div>
                                    <div className={line}>
                                        <input className="input-box w-[34%]" placeholder="Name" value={newOwnerName} onChange={(a) => {setNewOwnerName(a.target.value)}} />
                                        <input className="input-box w-[34%]" placeholder="Email" value={newOwnerEmail} onChange={(a) => {setNewOwnerEmail(a.target.value)}}/>
                                        <input className="input-box w-[34%]" placeholder="Phone" value={newOwnerPhone} onChange={phoneChange}/>
                                    </div>
                                <button className="black-button w-[34%]" onClick={addOwner}>Add</button>
                                <button onClick={() => setShowPopup2(false)}className="absolute top-4 right-4 text-xl font-bold cursor-pointer">x</button>
                            </div>
                        </div>
                    </div>
                    )}
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

export default Apartments;
