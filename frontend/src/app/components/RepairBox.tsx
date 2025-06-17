import React from "react";
import { useState, useEffect } from "react";
import Cookies from "js-cookie";
import { useRouter, usePathname } from 'next/navigation';

interface RepairProps {
  id: number
  title: string;
  assigned_date: string;
  completed_date: string;
  status: number;
  subcontractor: string;
  subcontractors: {id: number, user_id: number, address: string, nip: string, speciality_id: number, name: string}[];
  changeRepairSubcon: (sub: number, id: number) => void;
  changeRepairStatus: (id: number, s: string) => void;
}

function RepairBox(props: RepairProps) {
  const [showPopup,setShowPopup] = useState(false)
  const [showPopup2,setShowPopup2] = useState(false)
  const [showPopup3,setShowPopup3] = useState(false)

  const [subcontractor, setSubcontractor] = useState(1)

  const [state, setState] = useState('pending')


  const [role, setRole] = useState('')

  const pathname = usePathname();
  useEffect(() => {
          var a = Cookies.get("role");
          if(a != null)
          {
            setRole(a)
          }
        },[pathname])
  return (
    <div className="white-box w-[100%] h-[200px]">
      <div className="flex flex-row items-center justify-between gap-8 w-full px-8 ">
        <div className="flex flex-col w-[50%]">
          <b className="text-xl">Title: {props.title}</b>
          <h1 className="text-xl">Status: {props.status == 1 ? "Pending" : props.status == 2 ? "In-Progress" : "Completed"}</h1>
          <button className="black-button w-[100%]" onClick={()=>{setShowPopup(true)}}>View Details</button>
        </div>
        <div className={props.status == 1 ? "status-box-red" : props.status == 2 ? "status-box-yellow" : "status-box-green"}></div>
      </div>
      {showPopup && (
                <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/50">
                    <div className="white-box w-[30%] py-4 rounded-lg relative">
                        <div className="flex flex-col gap-2 w-[100%]">
                            <b className="text-4xl">Repair Details</b>
                            <h1>Title</h1>
                            <h1 className="input-box">{props.title}</h1>
                            <h1>Date Assigned: {props.assigned_date.split("T")[0]}</h1>
                            <h1>Date Completed: {props.status == 3 ? props.completed_date.split("T")[0] : "-"}</h1>
                            <h1>Status: {props.status == 1 ? "Pending" : props.status == 2 ? "In-Progress" : "Completed"}</h1>
                            <h1>Assigned to: {props.subcontractor}</h1>
                            {role === "1" && <button className="black-button" onClick={() => setShowPopup2(true)}>Change Subcontractor</button>}
                            {(role === "1" || role === "3") && <button className="black-button" onClick={() => setShowPopup3(true)}>Change Status</button>}
                        </div>
                        <button onClick={() => setShowPopup(false)}className="absolute top-4 right-4 text-xl font-bold cursor-pointer">x</button>
                    </div>
                </div>
      )}
      {showPopup2 && (
                <div className="fixed inset-0 z-40 flex items-center justify-center bg-black/50">
                    <div className="white-box w-[20%] py-4 rounded-lg relative">
                        <div className="flex flex-col gap-2 w-[100%]">
                            <b className="text-4xl">Change Subcontractor</b>
                            <b className=" w-[100%]">Subcontractor</b>
                            <select className="input-box  w-[100%]" value={subcontractor}onChange={(a) => {setSubcontractor(Number(a.target.value))}}>
                                {props.subcontractors.map((a,index) => (<option key={index} value={a.id}>{a.name}</option>))}
                            </select>
                            <button className="black-button  w-[100%]" onClick={()=>{props.changeRepairSubcon(subcontractor, props.id) ; setShowPopup2(false)}}>Change</button>
                        </div>
                        <button onClick={() => setShowPopup2(false)}className="absolute top-4 right-4 text-xl font-bold cursor-pointer">x</button>
                    </div>
                </div>
      )}
      {showPopup3 && (
                <div className="fixed inset-0 z-40 flex items-center justify-center bg-black/50">
                    <div className="white-box w-[20%] py-4 rounded-lg relative">
                        <div className="flex flex-col gap-2 w-[100%]">
                            <b className="text-4xl">Change Status</b>
                            <div className="flex flex-row gap-1">
                                <input type="radio" value= "pending" checked={state === "pending"} onChange={(a) => {setState(a.target.value)}}></input>
                                <h1>Pending</h1>
                            </div>
                            <div className="flex flex-row gap-1">
                                <input type="radio" value= "in_progress" checked={state === "in_progress"} onChange={(a) => {setState(a.target.value)}}></input>
                                <h1>In-Progress</h1>
                            </div>
                            <div className="flex flex-row gap-1">
                                <input type="radio" value= "completed" checked={state === "completed"} onChange={(a) => {setState(a.target.value)}}></input>
                                <h1>Completed</h1>
                            </div>
                            <button className="black-button  w-[100%]" onClick={() => {props.changeRepairStatus(props.id, state); setShowPopup3(false)}}>Change</button>
                        </div>
                        <button onClick={() => setShowPopup3(false)}className="absolute top-4 right-4 text-xl font-bold cursor-pointer">x</button>
                    </div>
                </div>
      )}
    </div>
  );
}

RepairBox.defaultProps = {
  title: "Fix cable",
  assigned_date: "30 April 2137",
  completed_date: "30 April 2137",
  status: -1,
  subcontractor: "John the Guy"
};

export default RepairBox;
