import React from "react";
import { useState, useEffect } from "react";
import Cookies from "js-cookie";

interface ApartmentProps {
    id: number
    name: String
    street: String
    building_number: String
    building_name: String
    flat_number: String
    owner_id: number
    rent: number
    refresh: Function
  }

function ApartmentBox(props : ApartmentProps)
{
    const [showPopup,setShowPopup] = useState(false)
    const [newRent, setNewRent] = useState('')
    const formatNumber = (value: string) => {
        const digits = value.replace(/[^\d.]/g, '');
        return digits;
    }
    function numberChange(e: React.ChangeEvent<HTMLInputElement>){
        const a = formatNumber(e.target.value);
        setNewRent(a);
    }

    async function changeRent()
    {
        var t = Cookies.get("token");
        try {
            const res = await fetch('http://localhost:8080/changerent',{
                method:'POST',
                body: JSON.stringify({ 
                    "token" : t,
                    "apartament_id" : props.id,
                    "rent": Number(newRent)
                })
            });
            if(res.ok)
            {
                alert("Rent changed succesfully.");
            }
            else
            {
                var data = await res.json()
                alert(data.message)
            }
            setShowPopup(false)
        } catch (err: any) {
            //alert(err.message)
            setShowPopup(false)
        } finally{
            props.refresh()
        }
    }

    return(
        <div className="flex justify-center w-[100%]">
            <div className="white-box w-[50%] h-[200px]">
                <div className="flex flex-col">
                    <h1>Name: {props.name}</h1>
                    <h1>Street: {props.street}</h1>
                    <h1>Building Number: {props.building_number}</h1>
                    <h1>Building Name: {props.building_name}</h1>
                    <h1>Flat Number: {props.flat_number}</h1>
                    {props.rent != -1 ? <h1>Rent: {props.rent}</h1> : <h1>Rent: Not Set</h1>}
                </div>
                <button className="black-button" onClick={()=>{setShowPopup(true)}}>Change Rent</button>
            </div>
            {showPopup && (
                <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/50">
                    <div className="white-box w-[40%] py-4 rounded-lg relative">
                        <div className="flex flex-col gap-2 w-[100%]">
                            <b className="text-4xl">Change Rent</b>
                            <div className="flex flex-row gap-1">
                                <input className="input-box w-[34%]" placeholder="New rent" value={newRent} onChange={numberChange}/>
                                <button className="black-button w-[34%]" onClick={changeRent}>Change</button>
                            </div>
                        </div>
                        <button onClick={() => setShowPopup(false)}className="absolute top-4 right-4 text-xl font-bold cursor-pointer">x</button>
                    </div>
                </div>
            )}
        </div>
    )
}

ApartmentBox.defaultProps = {
    id: -1,
    name: "Default Name",
    street: "Default Street",
    building_number: "Default Building Number",
    building_name: "Default Building Name",
    flat_number: "Default Flat Number",
    owner_id: -1,
    rent: -1
}
export default ApartmentBox;