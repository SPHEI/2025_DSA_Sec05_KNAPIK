import React from "react";

interface TenantProps {
    id: number
    name: String
    email: String
    phone: String
    role_id: number
  }

function TenantBox(props : TenantProps)
{
    return(
        <div className="white-box w-[50%] h-[200px]">
            <div className="flex flex-col">
                <h1>Name: {props.name}</h1>
                <h1>Email: {props.email}</h1>
                <h1>Phone: {props.phone}</h1>
                <h1>Apartment: 1 demo street</h1>
                <h1>Monthly Rent: 0$</h1>
                <h1>Rent Status: Paid</h1>
                <h1>Last Payment: 01.01.1971</h1>
            </div>
            <div className="flex flex-row justify-end items-center gap-10">
                <p className="status-box-green"></p>
                <div className="flex flex-col w-[150px] gap-2">
                    <button className="black-button">Edit</button>
                    <button className="black-button">Change Rent</button>
                    <button className="black-button">Delete</button>
                    <button className="black-button">View History</button>
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