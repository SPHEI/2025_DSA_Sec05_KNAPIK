"use client";
import "../globals.css";
import { useState, useEffect } from "react";
import { useRouter, usePathname } from 'next/navigation';
import Cookies from "js-cookie";

//Tenants Only
function SubmitIssue() {
  const [ready, setReady] = useState(false);
  const [error, setError] = useState("none");
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");

  const [apartaments,setApartaments] = useState([{id: -1,name: '', street: '', building_number: '', building_name: '',flat_number:'',owner_id:-1, rent: -1 }])
  const [apartment, setApartment] = useState(1);
  const [role, setRole] = useState('')

  async function handleSubmit()
  {
    var t = Cookies.get("token");
    const res = await fetch('http://localhost:8080/faults/add',{
                method:'POST',
                body: JSON.stringify({ 
                    "token" : t,
                    "fault" : {
                      "title" : title,
                      "description": description,
                      "status_id" : 1,
                      "apartment_id" : apartment
                    }
                })
            });
            if(res.ok)
            {
                alert("Issue submitted succesfully.");
            }
            else
            {
                var data = await res.json()
                alert(data.message)
                console.log(data.message)
            }
  };

  const pathname = usePathname();
  useEffect(() => {
    refresh()
  }, [pathname]);

  async function refresh()
  {
    var a = String(Cookies.get("role"))
    setRole(a)
    var t = Cookies.get("token");
    if(a === '1')
    {
        try{
        const res = await fetch('http://localhost:8080/apartament/list?token=' + t)
            const data = await res.json();
            if(data.message)
            {
                setError(data.message)
            }
            else
            {
                setApartaments(data);
            }
        }
        catch(err: any)
        {
            setError(err.message);
        }
    }
    else
    {
      try{
        const res = await fetch('http://localhost:8080/tenant/info?token=' + t)
            const data = await res.json();
            if(data.message)
            {
                setError(data.message)
            }
            else
            {
              //alert(data.apartament_id)
                setApartment(data.apartament_id);
            }
        }
        catch(err: any)
        {
            setError(err.message);
        }
    }
    setReady(true);
  }
  if (ready) {
    
    if (error == "none") {
      return (
        <main className="flex flex-col items-center justify-start gap-6 mt-10">
          <div className="page-head w-[50%]">
            <b className="text-4xl">Please describe your issue here</b>
          </div>

          <input
            type="text"
            placeholder="Title....."
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="input-box w-[50%]"
          />

          {role === '1' && (
            <div className="w-[50%]">
            <b>Apartament</b>
              <select className="input-box w-[100%]" value={apartment} onChange={(a) => {setApartment(Number(a.target.value))}}>
                  {apartaments.map((a,index) => (<option key={index} value={a.id}>{a.name}</option>))}
              </select>
            </div>
        )}

          <textarea
            placeholder="Describe the issue in detail..."
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className="input-box w-[50%] h-[300px]"
          />

          <button onClick={handleSubmit} className="black-button">
            Submit Issue
          </button>
        </main>
      );
    } else {
      return (
        <main>
          <b>An error has occured:</b>
          <h1>{error}</h1>
        </main>
      );
    }
  } else {
    return (
      <main>
        <h1>Loading...</h1>
      </main>
    );
  }
}

export default SubmitIssue;
