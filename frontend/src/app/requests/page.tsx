"use client";
import "../globals.css";
import { useState, useEffect } from "react";
import RequestBox from "../components/RequestBox";
import RepairBox from "../components/RepairBox";
import { useRouter, usePathname } from 'next/navigation';
import Cookies from "js-cookie";

//Page is shared by all types of accounts
function Requests() {
  const [ready, setReady] = useState(false);
  const [error, setError] = useState("none");

  const [requests, setRequests] = useState([{id: -1, description: '', date_reported: '', status_id: -1, apartament_id: -1}])

  const [role, setRole] = useState('')
  const pathname = usePathname();
  useEffect(() => {
    const fetchData = async () => {
          //Page setup goes here
          var a = Cookies.get("role");
          if(a != null)
          {
            setRole(a)
            var t = Cookies.get("token");
            try{
            const res = await fetch('http://localhost:8080/faults/list?token=' + t)
                const data = await res.json();
                if(data.message)
                {
                    setError(data.message)
                }
                else
                {
                    setRequests(data.faults);
                }
            }
            catch(err: any)
            {
                setError(err.message);
            }
          }
        }
        fetchData();
    setReady(true);
  }, [pathname]);

  if (ready) {
    if (error == "none") {
      if(role == "3")
      {
        return (
          <main>
            <div className="page-head w-[50%]">
              <b className="text-4xl">Assigned Repairs</b>
            </div>
            <div className="flex flex-col w-[50%] gap-5">
              <RepairBox title="Fix cable" assigned_date="1-01-2024" completed_date="" status="in-progress" subcontractor="John" />
            </div>
          </main>
        );
      }
      else
      {
        return (
          <main>
            <div className="page-head w-[50%]">
              <b className="text-4xl">My requests</b>
            </div>
            {requests != null ? requests.map((a, index) => <RequestBox key={index} id={a.id} title="null" description={a.description} date={a.date_reported} status={a.status_id}/>)
            : <h1>No requests</h1>}
          </main>
        );
      }
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

export default Requests;
