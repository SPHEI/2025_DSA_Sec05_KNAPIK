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

    useEffect(() => {
        const fetchData = async () => {
            
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
          };
          fetchData();
    },[pathname])

    if(ready)
    {
        if(error == 'none')
        {
            return (
                <main>
                    <div className="page-head w-[50%]">
                        <b className="text-4xl">Apartments</b> 
                        <button className="black-button">+ Add Apartment</button>
                    </div>
                    {apartments.map((a,index) => <ApartmentBox key={index} name={a}/>)}
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
