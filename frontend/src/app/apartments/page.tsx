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

    const [apartments, setApartments] = useState([''])

    const pathname = usePathname();

    const [showPopup,setShowPopup] = useState(false)


    const [name, setName] = useState('');
    const [street, setStreet] = useState('');
    const [number, setNumber] = useState('');
    const [bname, setBname] = useState('');
    const [fnumber,setFnumber] = useState('');
    const [email, setEmail] = useState('');

    const formatNumber = (value: string) => {
        const digits = value.replace(/\D/g, '');
        return digits;
    }
    function numberChange(e: React.ChangeEvent<HTMLInputElement>){
        const a = formatNumber(e.target.value);
        setNumber(a);
    }


    useEffect(() => {
        refresh()
    },[pathname])
    
    async function refresh()
    {
        try {
            var t = Cookies.get("token");
            const res = await fetch('http://localhost:8080/apartaments',{
            method: 'POST',
            body: JSON.stringify({
                "token" : t 
            })
            });
            const data = await res.json();
            //alert(JSON.stringify(data));
            if(data.message)
            {
            setError(data.message)
            }
            else
            {
            setApartments(data.apartaments)
            }
        } catch (err: any) {
            setError(err.message)
        } finally{
            setReady(true);
        }
    }

    async function addApartment()
    {
        var t = Cookies.get("token");
        // alert(JSON.stringify({
        //     "token" : t,
        //     "owner_email" : email,
        //     name,
        //     street,
        //     "building_number" : number,
        //     "building_name" : bname,
        //     "flat_number" : fnumber 
        // }));
        try {
            const res = await fetch('http://localhost:8080/addapartament',{
                method:'POST',
                body: JSON.stringify({ 
                    "token" : t,
                    "owner_email" : email,
                    name,
                    street,
                    "building_number" : number,
                    "building_name" : bname,
                    "flat_number" : fnumber 
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
                    {apartments.map((a,index) => <ApartmentBox key={index} name={a}/>)}
                    {showPopup && (
                    <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/50">
                        <div className="white-box w-[40%] h-[40%] rounded-lg relative">
                            <div className="flex flex-col gap-2 w-[100%]">
                                <b className="text-4xl">Add Apartment</b>
                                <div className={line}>
                                    <b className="w-[34%]">Name</b>
                                    <b>Owner E-mail</b>
                                </div>
                                <div className={line}>
                                    <input className="input-box" placeholder="Name" value={name} onChange={(a) => {setName(a.target.value)}}/>
                                    <input className="input-box" placeholder="E-mail" value={email} onChange={(a) => {setEmail(a.target.value)}}/>
                                </div>
                                <div className={line}>
                                    <b className="w-[34%]">Street</b>
                                    <b>Building Number</b>
                                </div>
                                <div className={line}>
                                    <input className="input-box" placeholder="Street" value={street} onChange={(a) => {setStreet(a.target.value)}} />
                                    <input className="input-box" placeholder="Number" value={number} onChange={numberChange}/>
                                </div>
                                <div className={line}>
                                    <b className="w-[34%]">Building Name</b>
                                    <b>Flat Number</b>
                                </div>
                                <div className={line}>
                                    <input className="input-box" placeholder="Building name" value={bname} onChange={(a) => {setBname(a.target.value)}} />
                                    <input className="input-box" placeholder="Flat number" value={fnumber} onChange={(a) => {setFnumber(a.target.value)}} />
                                </div>
                                <button className="black-button w-[34%]" onClick={addApartment}>Add</button>
                            </div>
                            <button onClick={() => setShowPopup(false)}className="absolute top-4 right-4 text-xl font-bold cursor-pointer">x</button>
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
