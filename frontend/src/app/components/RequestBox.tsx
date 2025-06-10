import React from "react";
import { useState, useEffect } from "react";
import Cookies from "js-cookie";
import { useRouter, usePathname } from 'next/navigation';

import RepairBox from "./RepairBox";


interface RequestProps {
  id: number;
  title: string;
  description: string;
  date: string;
  status: number;
  apartment_id: number;
  name: string;
  repairs : {
    ID: number;
    Title: string;
    DateAssigned: string;
    DateCompleted: string;
    StatusID: number;
    Subcontractor: string;
  }[]
  refresh: Function;
}

function RequestBox(props: RequestProps) {
  const [showPopup,setShowPopup] = useState(false)
  const [showPopup2,setShowPopup2] = useState(false)

  const [title, setTitle] = useState('')
  const [subcontractor, setSubcontractor] = useState(1)

  const [subcontractors, setSubcontractors] = useState([{ID: -1, UserID: -1, Address: '', Nip: '', SpecialityID: -1, Name: ''}])

  const [role, setRole] = useState('')

  const pathname = usePathname();
  async function changeStatus()
  {
    var t = Cookies.get("token");
    var i = 1;
    if(props.status == 1)
    {
      i = 2;
    }
    const res = await fetch('http://localhost:8080/faults/status',{
              method:'POST',
              body: JSON.stringify({ 
                  "token" : t,
                  "fault" : {
                    StatusID : i,
                    ID : props.id
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
        props.refresh()
  }
useEffect(() => {
        const fetchData = async () => {
          //Page setup goes here
          var a = Cookies.get("role");
          if(a != null)
          {
            setRole(a)
            if(a == "1")
            {
              var t = Cookies.get("token");
              const res = await fetch('http://localhost:8080/subcon/list?token=' + t)
              const data = await res.json();
              //alert(JSON.stringify(data))
              setSubcontractors(data);
            }
          }
        }
        fetchData();
    },[pathname])

  async function addRepair()
  {
    var t = Cookies.get("token");
    const res = await fetch('http://localhost:8080/repair/add',{
              method:'POST',
              body: JSON.stringify({ 
                  "token" : t,
                  "repair" : {
                    Title : title,
                    FaultReportID : props.id,
                    DateAssigned : "2006-01-02T15:04:05Z"
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
          props.refresh()
          setShowPopup2(false)
  }

  const line = "flex flex-row gap-1";
  return (
    <div className="white-box w-[50%] h-[200px]">
      <div className="flex flex-row items-center justify-between gap-8 w-full px-8 ">
        <div className="flex flex-col w-[50%]">
          <b className="text-xl">Title: {props.title}</b>
          <h1 className="text-xl">Status: {props.status == 1 ? "Open" : "Closed"}</h1>
          <button className="black-button w-[100%]" onClick={()=>{setShowPopup(true)}}>View Details</button>
        </div>
        <div className={props.status == 1 ? "status-box-yellow" : "status-box-green"}></div>
      </div>
      {showPopup && (
                <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/50 overflow-auto">
                    <div className="white-box w-[60%] py-4 rounded-lg relative max-h-[80vh] overflow-y-auto">
                        <div className={"flex flex-col gap-2 w-[100%] py-4 relative top-" + String(props.repairs.length * 50)}>
                            <b className="text-4xl">Request Details</b>
                            <h1>Title</h1>
                            <h1 className="input-box">{props.title}</h1>
                            <h1>Description</h1>
                            <h1 className="input-box">{props.description}</h1>
                            <h1>Related Apartment: {props.name}</h1>
                            <h1>Date Submitted: {props.date.split("T")[0]}</h1>
                            <h1>Status: {props.status == 1 ? "Open" : "Closed"}</h1>
                            <button className="black-button" onClick={changeStatus}>{props.status == 1 ? "Close Request" : "Reopen Request"}</button>
                            <div></div><div></div><div></div>
                            <div className="page-head w-[100%]">
                              <b className="text-4xl">Associated Repairs:</b>
                              {role === "1" && <button className="black-button w-[50%]" onClick={() => {setShowPopup2(true)}}>+ Add Repair</button>}
                            </div>
                            {props.repairs != null ? props.repairs.map((a, index) => <RepairBox
                                key={index} id={a.ID} title={a.Title} assigned_date={a.DateAssigned} completed_date={a.DateCompleted} status={a.StatusID} subcontractor={a.Subcontractor}
                                refresh={props.refresh}/>)
                                : <h1>No Repairs</h1>
                            }
                        </div>
                        <button onClick={() => setShowPopup(false)}className="absolute top-4 right-4 text-xl font-bold cursor-pointer">x</button>
                    </div>
                </div>
      )}
      {showPopup2 && (
                <div className="fixed inset-0 z-40 flex items-center justify-center bg-black/50">
                    <div className="white-box w-[50%] py-4 rounded-lg relative max-h-[80vh]">
                        <div className="flex flex-col gap-2 w-[100%] py-4">
                            <b className="text-4xl">Add Repair</b>
                            <div className={line}>
                                    <b className="w-[34%]">Title</b>
                                    {/* <b className="w-[34%]">Subcontractor</b> */}
                                </div>
                                <div className={line}>
                                    <input className="input-box w-[34%]" placeholder="Name" value={title} onChange={(a) => {setTitle(a.target.value)}}/>
                                    {/* <select className="input-box w-[34%]" value={subcontractor}onChange={(a) => {setSubcontractor(Number(a.target.value))}}>
                                        {subcontractors.map((a,index) => (<option key={index} value={a.ID}>{a.Name}</option>))}
                                    </select> */}
                                </div>
                                <button className="black-button w-[34%]" onClick={addRepair}>Add</button>
                        </div>
                        <button onClick={() => setShowPopup2(false)}className="absolute top-4 right-4 text-xl font-bold cursor-pointer">x</button>
                    </div>
                </div>
      )}
    </div>
  );
}

RequestBox.defaultProps = {
  id: -1,
  title: "Electricity",
  description: "blablabla",
  date: "30 April 2137",
  status: -1,
  apartment_id: -1,
  name: 'tets'
};

export default RequestBox;
