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
  subcontractors: {id: number, user_id: number, address: string, nip: string, speciality_id: number, name: string}[];
  changeRepairSubcon: (sub: number, id: number) => void;
  changeRepairStatus: (id: number, s: string) => void;
  changeFaultStatus: (i: number, id: number) => void;
  addRepair: (id: number, tit: string) => void;

}

function RequestBox(props: RequestProps) {
  const [showPopup,setShowPopup] = useState(false)
  const [showPopup2,setShowPopup2] = useState(false)

  const [title, setTitle] = useState('')

  const [role, setRole] = useState('')

  const pathname = usePathname();
useEffect(() => {
        //Page setup goes here
        var a = Cookies.get("role");
        if(a != null)
        {
          setRole(a)
        }
    },[pathname])

    const [sort, setSort] = useState('None')


    const priorityInProgress: Record<number, number> = {
        3 : 2,
        2:  3,
        1 : 1,
    };
    

  const line = "flex flex-row gap-1";
  return (
    <div className="white-box w-[50%] h-[200px] min-w-[500px]">
      <div className="flex flex-row items-center justify-between gap-8 w-full px-8 ">
        <div className="flex flex-col w-[50%]">
          <b className="text-xl">Title: {props.title}</b>
          <h1 className="text-xl">Status: {props.status == 1 ? "Open" : "Closed"}</h1>
          <button className="black-button w-[100%]" onClick={()=>{setShowPopup(true)}}>View Details</button>
        </div>
        <div className={props.status == 1 ? "status-box-yellow" : "status-box-green"}></div>
      </div>
      {showPopup && (
                <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/50">
                    <div className="white-box w-[60%] min-w-[600px] max-h-[80vh] rounded-lg bg-white flex flex-col relative">
                      <button onClick={() => setShowPopup(false)} className="absolute top-4 right-4 text-xl font-bold cursor-pointer">x</button>
                      <div className="overflow-y-auto w-[100%] px-6 py-4 flex-1">
                        <div className="flex flex-col gap-4 w-full">
                          <b className="text-4xl">Request Details</b>
                          <h1>Title</h1>
                          <h1 className="input-box">{props.title}</h1>
                          <h1>Description</h1>
                          <h1 className="input-box">{props.description}</h1>
                          <h1>Related Apartment: {props.name}</h1>
                          <h1>Date Submitted: {props.date.split("T")[0]}</h1>
                          <h1>Status: {props.status == 1 ? "Open" : "Closed"}</h1>
                          <button className="black-button"onClick={() => props.changeFaultStatus(props.status, props.id)}>
                            {props.status == 1 ? "Close Request" : "Reopen Request"}
                          </button>
                          <div className="page-head w-full flex items-center gap-2 justify-between mt-6">
                            <b className="text-4xl">Associated Repairs:</b>
                            <h1 className="text-2xl">Sort:</h1>
                            <select className="input-box w-[25%]" value={sort} onChange={(a) => {setSort(a.target.value)}}>
                                <option value="None">None</option>
                                <option value="Completed">Completed</option>
                                <option value="In-Progress">In-Progress</option>
                                <option value="Pending">Pending</option>
                            </select>
                            {role === "1" && (
                              <button className="black-button w-[30%]"onClick={() => setShowPopup2(true)}>+ Add Repair</button>
                            )}
                          </div>
                          {props.repairs.length > 0 ? (sort == 'None' ? props.repairs : [...props.repairs].sort((a, b) => {
                                const priorityA = sort == "Completed" ? a.StatusID : sort == "Pending" ? 4 - a.StatusID : priorityInProgress[a.StatusID];
                                const priorityB = sort == "Completed" ? b.StatusID : sort == "Pending" ? 4 - b.StatusID : priorityInProgress[b.StatusID];
                                return priorityB - priorityA;
                            })).map((a, index) => (
                              <RepairBox key={index} id={a.ID} title={a.Title}
                                assigned_date={a.DateAssigned} completed_date={a.DateCompleted}
                                status={a.StatusID}
                                subcontractor={a.Subcontractor} subcontractors={props.subcontractors}
                                changeRepairSubcon={props.changeRepairSubcon} changeRepairStatus={props.changeRepairStatus}
                              />
                            )) : (<h1>No Repairs</h1>)}
                        </div>
                      </div>
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
                                <button className="black-button w-[34%]" onClick={() => {props.addRepair(props.id, title); setShowPopup2(false)}}>Add</button>
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
