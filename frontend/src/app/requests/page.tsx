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

  const [requests, setRequests] = useState([{id: -1, title : '',description: '', date_reported: '', status_id: -1, apartment_id: -1, name: ''}])
  const [repairs, setRepairs] = useState([{id: -1, title: '', fault_report_id: -1, date_assigned: '', date_completed: '', status_id: -1, subcontractor_id: -1, name: ''}])

  const [role, setRole] = useState('')
  const pathname = usePathname();
  const [subcontractors, setSubcontractors] = useState([{ID: -1, UserID: -1, Address: '', Nip: '', SpecialityID: -1, Name: ''}])
  

  function mapRepairsToRequest(faultID: number)
  {
    const a = repairs.filter(item => item.fault_report_id == faultID);
    var b: {
          ID: number,
          Title: string,
          DateAssigned: string,
          DateCompleted: string,
          StatusID: number,
          Subcontractor: string
        }[] = []
    a.forEach((c) => {
      b.push(
        {
          ID: c.id,
          Title: c.title,
          DateAssigned: c.date_assigned,
          DateCompleted: c.date_completed,
          StatusID: c.status_id,
          Subcontractor: c.name
        }
      )
    })
    return b;
  }

  useEffect(() => {
    refresh()
    setReady(true);
  }, [pathname]);

  async function refresh()
  {
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
                  //alert(JSON.stringify(data))
                  setRequests(data);
                }
            }
            catch(err: any)
            {
                setError(err.message);
            }
            try{
                const res = await fetch('http://localhost:8080/repair/list?token=' + t)
                const data = await res.json();
                if(data.message)
                {
                  setError(data.message)
                }
                else
                {
                  //alert(JSON.stringify(data))
                  setRepairs(data);
                }
            }
            catch(err: any)
            {
                setError(err.message);
            }
            if(a == "1")
            {
              try{
                const res = await fetch('http://localhost:8080/subcon/list?token=' + t)
                const data = await res.json();
                //alert(JSON.stringify(data))
                setSubcontractors(data);
              }
              catch(err: any)
              {
                  setError(err.message);
              }
            }
          }
  }

  if (ready) {
    if (error == "none") {
      if(role == "3")
      {
        return (
          <main>
            <div className="page-head w-[50%] min-w-[500px]">
              <b className="text-4xl">Assigned Repairs</b>
            </div>
            <div className="flex flex-col w-[50%] gap-5">
              {repairs != null ? repairs.map((a, index) => <RepairBox
              key={index} id={a.id} title={a.title} assigned_date={a.date_assigned} completed_date={a.date_completed} status={a.status_id} 
              subcontractor={a.name} subcontractors={subcontractors}
              refresh={refresh}/>)
              : <h1>No Repairs</h1>
            }
            </div>
          </main>
        );
      }
      else
      {
        return (
          <main>
            <div className="page-head w-[50%] min-w-[500px]">
              <b className="text-4xl">My requests</b>
            </div>
            {requests != null ? requests.map((a, index) => <RequestBox 
            key={index} id={a.id} title={a.title}description={a.description} date={a.date_reported} status={a.status_id}
            apartment_id={a.apartment_id} name={a.name} repairs={mapRepairsToRequest(a.id)} refresh={refresh}
            subcontractors={subcontractors}/>)
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
