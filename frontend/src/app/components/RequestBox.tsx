import React from "react";
import { useState, useEffect } from "react";
import Cookies from "js-cookie";
import { useRouter, usePathname } from 'next/navigation';

import RepairBox from "./RepairBox";

interface RequestProps {
  title: string;
  description: string;
  date: string;
  status: string;
}

function RequestBox(props: RequestProps) {
  const [showPopup,setShowPopup] = useState(false)
  const [showPopup2,setShowPopup2] = useState(false)

  const [title, setTitle] = useState('')
  const [subcontractor, setSubcontractor] = useState(1)

  const [subcontractors, setSubcontractors] = useState([{id: 1, email: "no"}])

  const [role, setRole] = useState('')

  const pathname = usePathname();
useEffect(() => {
        const fetchData = async () => {
          //Page setup goes here
          var a = Cookies.get("role");
          if(a != null)
          {
            setRole(a)
          }
        }
        fetchData();
    },[pathname])


  const line = "flex flex-row gap-1";
  return (
    <div className="white-box w-[50%] h-[200px]">
      <div className="flex flex-row items-center justify-between gap-8 w-full px-8 ">
        <div className="flex flex-col w-[50%]">
          <b className="text-xl">Title: {props.title}</b>
          <h1 className="text-xl">Status: {props.status}</h1>
          <button className="black-button w-[100%]" onClick={()=>{setShowPopup(true)}}>View Details</button>
        </div>
        <div className={props.status === 'Open' ? "status-box-yellow" : "status-box-green"}></div>
      </div>
      {showPopup && (
                <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/50 overflow-auto">
                    <div className="white-box w-[60%] py-4 rounded-lg relative max-h-[80vh] overflow-y-auto">
                        <div className="flex flex-col gap-2 w-[100%] relative top-50 py-4">
                            <b className="text-4xl">Request Details</b>
                            <h1>Title</h1>
                            <h1 className="input-box">{props.title}</h1>
                            <h1>Description</h1>
                            <h1 className="input-box">{props.description}</h1>
                            <h1>Date Submitted: {props.date}</h1>
                            <h1>Status: {props.status}</h1>
                            <button className="black-button">{props.status === 'Open' ? "Close Request" : "Reopen Request"}</button>
                            <div></div><div></div><div></div>
                            <div className="page-head w-[100%]">
                              <b className="text-4xl">Associated Repairs:</b>
                              {role === "1" && <button className="black-button w-[50%]" onClick={() => {setShowPopup2(true)}}>+ Add Repair</button>}
                            </div>
                            <RepairBox title="Fix cable" assigned_date="1-01-2024" completed_date="" status="in-progress" subcontractor="John" />
                            <RepairBox title="Fix cable" assigned_date="1-01-2024" completed_date="" status="pending" subcontractor="John" />
                            <RepairBox title="Fix cable" assigned_date="1-01-2024" completed_date="" status="completed" subcontractor="John" />
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
                                    <b className="w-[34%]">Subcontractor</b>
                                </div>
                                <div className={line}>
                                    <input className="input-box w-[34%]" placeholder="Name" value={title} onChange={(a) => {setTitle(a.target.value)}}/>
                                    <select className="input-box w-[34%]" value={subcontractor}onChange={(a) => {setSubcontractor(Number(a.target.value))}}>
                                        {subcontractors.map((a,index) => (<option key={index} value={a.id}>{a.email}</option>))}
                                    </select>
                                </div>
                                <button className="black-button w-[34%]">Add</button>
                        </div>
                        <button onClick={() => setShowPopup2(false)}className="absolute top-4 right-4 text-xl font-bold cursor-pointer">x</button>
                    </div>
                </div>
      )}
    </div>
  );
}

RequestBox.defaultProps = {
  title: "Electricity",
  description: "blablabla",
  date: "30 April 2137",
  status: "Pending",
};

export default RequestBox;
