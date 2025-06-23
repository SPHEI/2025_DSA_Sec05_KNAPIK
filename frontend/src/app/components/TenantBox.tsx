import React from "react";
import { useState, useEffect } from "react";

interface TenantProps {
    id: number
    name: String
    email: String
    phone: String
    role_id: number
    apartment_id: number
    apartment: string
    rent: string
    status: string
    evict: Function
    renting_id: number

    changeRent: (id: number, newRent: number) => void;
    changeApartment: (a_id: number, u_id: number, date: string) => void;

    apartments: {id: number,name: string, street: string, building_number: string, building_name: string,flat_number:string,owner_id:number }[]
  }

function TenantBox(props : TenantProps)
{
    const [showChange, setShowChange] = useState(false)
    const [showChange2, setShowChange2] = useState(false)

    const [newRent, setNewRent] = useState('')

    const [apartment, setApartment] = useState(1);
    const [date, setDate] = useState('');

    const formatNumber = (value: string) => {
        const digits = value.replace(/[^\d.]/g, '');
        return digits;
    }
    function numberChange(e: React.ChangeEvent<HTMLInputElement>){
        const a = formatNumber(e.target.value);
        setNewRent(a);
    }

    function callChangeRent()
    {
        props.changeRent(props.apartment_id, Number(newRent))
        setShowChange2(false)
    }

    function callChangeApartment()
    {
        props.changeApartment(apartment, props.id, date)
        setShowChange(false)
    }

    const line = "flex flex-row gap-1";
    return(
        <div className="w-[50%]">
            <div className="white-box w-[100%] h-[200px] min-w-[600px]">
                <div className="flex flex-col">
                    <h1>Name: {props.name}</h1>
                    <h1>Email: {props.email}</h1>
                    <h1>Phone: {props.phone}</h1>
                    <h1>Apartment: {props.apartment != '' ? props.apartment : "None"}</h1>
                    {props.apartment != '' && <h1>Monthly Rent: {Number(props.rent) > 0 ? props.rent : "Not Set"}</h1>}
                    {props.apartment != '' && <h1>Rent Status: {props.status}</h1>}
                </div>
                <div className="flex flex-row justify-end items-center gap-10">
                    <p className={props.status == "Overdue" ? "status-box-red" : props.status == "Pending" ? "status-box-yellow" : "status-box-green"}></p>
                    <div className="flex flex-col w-[160px] gap-2">
                        <div></div>
                        <button className="black-button" onClick={()=>{setShowChange(true)}}>{props.apartment != '' ? "Change Apartment" : "Assign Apartment"}</button>
                        {props.apartment != '' && <button className="black-button" onClick={()=>{setShowChange2(true)}}>Change Rent</button>}
                        {props.apartment != '' && <button className="black-button" onClick={()=>{props.evict(props.renting_id)}}>Evict</button>}
                    </div>
                </div>
            </div>
            {showChange && (
                <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/50">
                    <div className="white-box w-[40%] py-4 rounded-lg relative min-w-[400px]">
                        <div className="flex flex-col gap-2 w-[100%] ">
                            <div>
                                <div className={line}>
                                    <b className="w-[26%]">Apartment</b>
                                    <b className="w-[26%]">Start Date</b>
                                </div>
                                <div className={line}>
                                    <select className="input-box w-[26%]" value={apartment} onChange={(a) => {setApartment(Number(a.target.value))}}>
                                        {props.apartments.map((a,index) => (<option key={index} value={a.id}>{a.name}</option>))}
                                    </select>
                                    <input
                                    type="date"
                                    className="input-box w-[26%]"
                                    value={date}
                                    onChange={(e) => setDate(e.target.value)}
                                    />
                                    <button className="black-button" onClick={callChangeApartment}> Change </button>
                                </div>
                            </div>
                        </div>
                        <button onClick={() => setShowChange(false)}className="absolute top-4 right-4 text-xl font-bold cursor-pointer">x</button>
                    </div>
                </div>
            )}
            {showChange2 && (
                <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/50">
                    <div className="white-box w-[40%] py-4 rounded-lg relative min-w-[400px]">
                        <div className="flex flex-col gap-2 w-[100%] ">
                            <b className="text-4xl">Change Rent</b>
                            <div className="flex flex-row gap-1">
                                <input className="input-box w-[34%]" placeholder="New rent" value={newRent} onChange={numberChange}/>
                                <button className="black-button w-[34%]" onClick={callChangeRent}>Change</button>
                            </div>
                        </div>
                        <button onClick={() => setShowChange2(false)}className="absolute top-4 right-4 text-xl font-bold cursor-pointer">x</button>
                    </div>
                </div>
            )}
        </div>
    )
}

TenantBox.defaultProps = {
    id: -1,
    name: "Default Name",
    email: "default@default.def",
    phone: "999",
    role_id: -1
}
export default TenantBox;