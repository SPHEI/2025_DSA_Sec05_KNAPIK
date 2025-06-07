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
              <RepairBox title="Fix cable" assigned_date="1-01-2024" completed_date="" status="pending" subcontractor="John" />
              <RepairBox title="Fix cable" assigned_date="1-01-2024" completed_date="" status="completed" subcontractor="John" />
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
            <RequestBox
              title="Water leak"
              description="The water is truly leaking my lordThe water is truly leaking my lordThe water is truly leaking my lordThe water is truly leaking my lordThe water is truly leaking my lordThe water is truly leaking my lordThe water is truly leaking my lordThe water is truly leaking my lordThe water is truly leaking my lordThe water is truly leaking my lordThe water is truly leaking my lordThe water is truly leaking my lordThe water is truly leaking my lordThe water is truly leaking my lordThe water is truly leaking my lordThe water is truly leaking my lord"
              date="10 May 2025"
              status="Open"
            />
            <RequestBox
              title="Water leak"
              description="The water is truly leaking my lord"
              date="10 May 2025"
              status="Closed"
            />
            <RequestBox
              title="Water leak"
              description="The water is truly leaking my lord"
              date="10 May 2025"
              status="Open"
            />
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
