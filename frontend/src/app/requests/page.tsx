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
  const [subcontractors, setSubcontractors] = useState([{id: -1, user_id: -1, address: '', nip: '', speciality_id: -1, name: ''}])

  const [sort, setSort] = useState('None')
  

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
                if(data != null)
                {
                  if(data.message)
                  {
                    console.log(data.message)
                    setRequests([])
                  }
                  else
                  {
                    //alert(JSON.stringify(data))
                    setRequests(data);
                  }
                }
                else
                {
                  setRequests([])
                }
            }
            catch(err: any)
            {
                setError(err.message);
            }
            try{
                const res = await fetch('http://localhost:8080/repair/list?token=' + t)
                const data = await res.json();
                if(data != null)
                {
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
                else
                {
                  setRepairs([])
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

  async function changeRepairSubcon(sub: number, id: number)
  {
    //alert(sub + " " + id)
    var t = Cookies.get("token");
    const res = await fetch('http://localhost:8080/repair/contractor',{
              method:'POST',
              body: JSON.stringify({ 
                  "token" : t,
                  "contractor" : {
                    "subcontractor_id" : sub,
                    "id" : id
                    }
              })
          });
          if(res.ok)
          {
              //alert("Subcontractor changed succesfully.");
          }
          else
          {
              var data = await res.json()
              alert(data.message)
          }
          refresh()
  }

  async function changeRepairStatus(id: number, s: string)
  {
    
    //alert(id + " " + s)
    var t = Cookies.get("token");
    var d = new Date()
    var dd = String(d.getDate()).padStart(2,'0')
    var m = String(d.getMonth() + 1).padStart(2,'0')
    var y = String(d.getFullYear())

    var date = y + "-" + m + "-" + dd
    var t = Cookies.get("token");
    const res = await fetch('http://localhost:8080/repair/data',{
              method:'POST',
              body: JSON.stringify({ 
                  "token" : t,
                  "repair" : {
                    "name": s,
                    "date_completed" : date + "T15:04:05Z",
                    "id": id
                    }
              })
          });
    if(res.ok)
    {
        //alert("Status changed succesfully.");
    }
    else
    {
        var data = await res.json()
        alert(data.message)
    }
    refresh()
  }

  async function changeFaultStatus(i: number, id: number)
  {
    //alert(id)
    var t = Cookies.get("token");
    if(i == 1)
    {
      i = 2
    }
    else
    {
      i = 1
    }
    const res = await fetch('http://localhost:8080/faults/status',{
              method:'POST',
              body: JSON.stringify({ 
                  "token" : t,
                  "fault" : {
                    "status_id" : i,
                    "id" : id
                  }
              })
          });
          if(res.ok)
          {
            //alert("a")
          }
          else
          {
              var data = await res.json()
              alert(data.message)
          }
        refresh()
  }

  async function addRepair(id: number, tit: string)
  {
    var t = Cookies.get("token");
    var d = new Date()
    var dd = String(d.getDate()).padStart(2,'0')
    var m = String(d.getMonth() + 1).padStart(2,'0')
    var y = String(d.getFullYear())

    var date = y + "-" + m + "-" + dd
    const res = await fetch('http://localhost:8080/repair/add',{
              method:'POST',
              body: JSON.stringify({ 
                  "token" : t,
                  "repair" : {
                    "title" : tit,
                    "fault_report_id" : id,
                    "date_assigned" : date + "T15:04:05Z"
                  }
              })
          });
          if(res.ok)
          {
              //alert("Repair created succesfully.");
          }
          else
          {
              var data = await res.json()
              alert(data.message)
          }
          refresh()
  }

  const priorityInProgress: Record<number, number> = {
        3 : 2,
        2:  3,
        1 : 1,
    };

  if (ready) {
    if (error == "none") {
      if(role == "3")
      {
        return (
          <main>
            <div className="page-head w-[50%] min-w-[500px]">
              <b className="text-4xl">Assigned Repairs</b>
              <div className="flex flex-row gap-1">
                <h1 className="text-2xl">Sort:</h1>
                <select className="input-box w-[70%]" value={sort} onChange={(a) => {setSort(a.target.value)}}>
                    <option value="None">None</option>
                    <option value="Completed">Completed</option>
                    <option value="In-Progress">In-Progress</option>
                    <option value="Pending">Pending</option>
                </select>
              </div>
            </div>
            <div className="flex flex-col w-[50%] gap-5">
              {repairs.length > 0 ? (sort == 'None' ? repairs : [...repairs].sort((a, b) => {
                  const priorityA = sort == "Completed" ? a.status_id : sort == "Pending" ? 4 - a.status_id : priorityInProgress[a.status_id];
                  const priorityB = sort == "Completed" ? b.status_id : sort == "Pending" ? 4 - b.status_id : priorityInProgress[b.status_id];
                  return priorityB - priorityA;
              })).map((a, index) => <RepairBox
              key={index} id={a.id} title={a.title} assigned_date={a.date_assigned} completed_date={a.date_completed} status={a.status_id} 
              subcontractor={a.name} subcontractors={subcontractors}
              changeRepairSubcon={changeRepairSubcon} changeRepairStatus={changeRepairStatus}/>)
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
              <div className="flex flex-row gap-1">
                <h1 className="text-2xl">Sort:</h1>
                <select className="input-box w-[100%]" value={sort} onChange={(a) => {setSort(a.target.value)}}>
                    <option value="None">None</option>
                    <option value="Open">Open</option>
                    <option value="Closed">Closed</option>
                </select>
              </div>
            </div>
            {requests.length > 0  ? (sort == 'None' ? requests : [...requests].sort((a, b) => {
              const priorityA = a.status_id;
              const priorityB = b.status_id;
              return sort == 'Open' ? priorityA - priorityB : priorityB - priorityA;
            })).map((a, index) => <RequestBox 
            key={index} id={a.id} title={a.title}description={a.description} date={a.date_reported} status={a.status_id}
            apartment_id={a.apartment_id} name={a.name} repairs={mapRepairsToRequest(a.id)} 
            changeRepairSubcon={changeRepairSubcon} changeRepairStatus={changeRepairStatus}
            changeFaultStatus={changeFaultStatus} addRepair={addRepair}
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
