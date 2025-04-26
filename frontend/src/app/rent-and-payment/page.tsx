"use client"
import "../globals.css";
import { useState, useEffect } from "react";

//Shared by admin and tenants
function RentAndPayment() {
    const [ready,setReady] = useState(false)

    useEffect(() => {
        //Page setup goes here
        setReady(true);
    })

    if(ready)
    {
        return (
            <main>
                <h1>rents and payments page</h1>
            </main>
        );
    }
};

export default RentAndPayment;
