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

    const [apartaments,setApartaments] = useState([{id: -1,name: '', street: '', building_number: '', building_name: '',flat_number:'',owner_id:-1, price: -1 }])

    const pathname = usePathname();

    const [showPopup,setShowPopup] = useState(false)

    const [name, setName] = useState('');
    const [street, setStreet] = useState('');
    const [number, setNumber] = useState('');
    const [bname, setBname] = useState('');
    const [fnumber,setFnumber] = useState('');

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
                setApartaments(data);
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
                    "aparment":{
                        "name": name,
                        "street": street,
                        "building_number" : number,
                        "building_name" : bname,
                        "flat_number" : fnumber,
                    }
                })
            });
            if(res.ok)
            {
                console.log("Apartment added succesfully.");
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

    const line = "flex flex-row gap-1";

    if(ready)
    {
        if(error == 'none')
        {
            return (
                <main>
                    <div className="page-head w-[50%] min-w-[450px]">
                        <b className="text-4xl">Apartments</b> 
                        <button className="black-button" onClick={() => {setShowPopup(true)}}>+ Add Apartment</button>
                    </div>
                    {apartaments.map((a,index) => <ApartmentBox key={index} id={a.id} name={a.name} street={a.street} building_number={a.building_number} building_name={a.building_name} flat_number={a.flat_number} owner_id={a.owner_id} price={a.price} changeRent={changeRent}/>)}
                    {showPopup && (
                    <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/50">
                        <div className="white-box w-[40%] py-4 rounded-lg relative min-w-[500px]">
                            <div className="flex flex-col gap-2 w-[100%]">
                                <b className="text-4xl">Add Apartment</b>
                                <div className={line}>
                                    <b className="w-[34%]">Name</b>
                                </div>
                                <div className={line}>
                                    <input className="input-box w-[34%]" placeholder="Name" value={name} onChange={(a) => {setName(a.target.value)}}/>
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
