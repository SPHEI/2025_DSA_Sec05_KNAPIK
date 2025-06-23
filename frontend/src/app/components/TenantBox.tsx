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
  }

function TenantBox(props : TenantProps)
{
    const [showChange, setShowChange] = useState(false)
    const [showChange2, setShowChange2] = useState(false)
    const [showPayments, setShowPayments] = useState(false)
    return(
        <div className="white-box w-[50%] h-[200px] min-w-[600px]">
            <div className="flex flex-col">
                <h1>Name: {props.name}</h1>
                <h1>Email: {props.email}</h1>
                <h1>Phone: {props.phone}</h1>
                <h1>Apartment: {props.apartment != '' ? props.apartment : "None"}</h1>
                {props.apartment != '' && <h1>Monthly Rent: {Number(props.rent) > 0 ? props.rent : "Not Set"}</h1>}
                {props.apartment != '' && <h1>Rent Status: {props.status}</h1>}
            </div>
            <div className="flex flex-row justify-end items-center gap-10">
                <p className="status-box-green"></p>
                <div className="flex flex-col w-[160px] gap-2">
                    <button className="black-button" onClick={()=>{setShowChange(true)}}>Change Apartment</button>
                    <button className="black-button" onClick={()=>{setShowChange2(true)}}>Change Rent</button>
                    <button className="black-button" onClick={()=>{props.evict(props.apartment_id)}}>Evict</button>
                    <button className="black-button" onClick={()=>{setShowPayments(true)}}>View Payments</button>
                </div>
            </div>
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