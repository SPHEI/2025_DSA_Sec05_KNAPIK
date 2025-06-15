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

  const [requests, setRequests] = useState([{ID: -1, Title : '',Description: '', DateReported: '', StatusID: -1, ApartmentID: -1, Name: ''}])
  const [repairs, setRepairs] = useState([{ID: -1, Title: '', FaultReportID: -1, DateAssigned: '', DateCompleted: {Time: '', Valid: false}, StatusID: -1, SubcontractorID: {Int64: -1, Valid: false}, Name: {String: '', Valid: false}}])

  const [role, setRole] = useState('')
  const pathname = usePathname();

  function mapRepairsToRequest(faultID: number)
  {
    const a = repairs.filter(item => item.FaultReportID == faultID);
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
          ID: c.ID,
          Title: c.Title,
          DateAssigned: c.DateAssigned,
          StatusID: c.StatusID
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
          }
  }

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
              {repairs != null ? repairs.map((a, index) => <RepairBox
              key={index} id={a.ID} title={a.Title} assigned_date={a.DateAssigned} completed_date={a.DateCompleted.Time} status={a.StatusID} subcontractor={a.Name.String}
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
            <div className="page-head w-[50%]">
              <b className="text-4xl">My requests</b>
            </div>
            {requests != null ? requests.map((a, index) => <RequestBox 
            key={index} id={a.ID} title={a.Title}description={a.Description} date={a.DateReported} status={a.StatusID}
            apartment_id={a.ApartmentID} name={a.Name} repairs={mapRepairsToRequest(a.ID)} refresh={refresh}/>)
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
