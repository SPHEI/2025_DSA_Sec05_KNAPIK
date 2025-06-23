"use client";
import "../globals.css";
import { useState, useEffect } from "react";
import card from "./card.png";
import PaymentBox from "../components/PaymentBox";
import calendar from "./calendar.png";
import { useRouter, usePathname } from 'next/navigation';
import Cookies from "js-cookie";

//Shared by admin and tenants
function RentAndPayment() {
  const [ready, setReady] = useState(false);
  const [error, setError] = useState("none");
  const pathname = usePathname();

  const [payments, setPayments] = useState([{id: -1, amount: -1, payment_date: '', due_date: '', status_id: -1, renting_id: -1, transaction_reference: ''}])
  const [names,setNames] = useState([{id: -1, name: '', email: '', phone: '', role_id: -1}])

  const [role, setRole] = useState('')

  useEffect(() => {
        refresh()
    },[pathname])
    
    async function refresh()
    {
        var a = Cookies.get("role");
        if(a != null)
        {
          setRole(a)
          var t = Cookies.get("token");
          try{
              await fetch('http://localhost:8080/test')
              const res = await fetch('http://localhost:8080/payments/list?token=' + t)
              const data = await res.json();
              //alert(JSON.stringify(data))
              if(data != null)
              {
                if(data.message)
                {
                    setError(data.message)
                }
                else
                {
                    setPayments(data)
                }
              }
              else
              {
                setPayments([])
              }
              if(a === "1")
              {
                const res2 = await fetch('http://localhost:8080/tenant/list?token=' + t)
                const data2 = await res2.json();
                //alert(JSON.stringify(data2));
                if(data2.message)
                {
                  setError(data2.message)
                }
                else
                {
                  setNames(data2)
                }
              }
          }catch(err: any){
              setError(err.message);
          }finally{
              setReady(true);
          }
        }
    }

  async function pay(id: number)
  {
    var t = Cookies.get("token");
    var t = Cookies.get("token");
    var d = new Date()
    var dd = String(d.getDate()).padStart(2,'0')
    var m = String(d.getMonth() + 1).padStart(2,'0')
    var y = String(d.getFullYear())

    var date = y + "-" + m + "-" + dd
    const res = await fetch('http://localhost:8080/payments/pay',{
              method:'POST',
              body: JSON.stringify({ 
                  "token" : t,
                  "payment" : {
                    "payment_date" : date + "T15:04:05Z",
                    "id" : id
                  }
        })
    });
    if(res.ok)
    {
        //alert("Paid succesfully.");
    }
    else
    {
        var data = await res.json()
        alert(data.message)
    }
    refresh()
  }

  function isOverdue(d: string)
  {
    const date = new Date(d.split("T")[0]);
    const now = new Date();
    return date < now
  }

  if (ready) {
    if (error == "none") {
      return (
        <main>
          <div className="page-head w-[50%]">
            <b className="text-4xl">{role === "2" ? "Rent & Payments" : "Payments"}</b>
          </div>
          <div className="flex flex-row gap-4 w-[50%]">
            <div className="white-box h-[150px] w-[100%]">
              {role === "2" &&
              <div className="flex flex-col items-center justify-center">
                <img src={card.src} width={40} />
                
                <b>Rent Status</b>
                <h1>{
                  payments.some(a => a.status_id == 3) ? "Overdue" :
                  payments.some(a => a.status_id == 1) ? "Unpaid" : "Paid"
                }</h1>
              
              </div>
              }
              {role === "1" &&
              <div className="flex flex-col items-center justify-center">
                <img src={card.src} width={40} />
                
                <b>Overdue Payments</b>
                <h1>{
                  payments.filter(a => a.status_id == 3).length
                }</h1>
              
              </div>
              }
            </div>
            {role === "2" && payments.filter(a => a.status_id != 2).length > 0 &&
            <div className="white-box h-[150px] w-[100%]">
              <div className="flex flex-col items-center justify-center">
                <img src={calendar.src} width={40} />
                <b>Next Rent Due</b>
                <h1>{
                  payments.filter(a => a.status_id != 2)
                  .map(p => p.due_date)
                  .reduce((earliest, current) => {return new Date(current.split("T")[0]) < new Date(earliest.split("T")[0]) ? current : earliest;},"3000-01-01")
                  .split("T")[0]
              }</h1>
              </div>
            </div>
            }
          </div>
          <div className="white-box w-[50%] py-4">
            <div className="flex flex-col items-left justify-start w-full h-full gap-2">
              <b className="text-xl">Pending Payments</b>
              <div className="flex flex-col gap-1">
                {payments.map((a,index) => a.status_id != 2 && (
                  <div key={index} className="flex flex-row">
                    <PaymentBox 
                    date={"Due: " + a.due_date.split("T")[0]} 
                    name={role === "1" ? "Renting Id: " + (a.renting_id)  : ""} 
                    type={isOverdue(a.due_date) ? "Overdue" : ""} 
                    amount={a.amount} /> 
                    {role === "2" && <button className="black-button" onClick={()=>{pay(a.id)}}>Pay</button>}
                  </div>))}
              </div>
              <b className="text-xl">Payment history</b>
              <div className="flex flex-col gap-1">
                {payments.map((a,index) => a.status_id == 2 && 
                <PaymentBox key={index}  
                date={"Due: " + a.due_date.split("T")[0]}  
                name={role === "1" ? "Renting Id: " + (a.renting_id) : ""} 
                type={"Paid on: " + a.payment_date.split("T")[0]} 
                amount={a.amount} />)}
              </div>
            </div>
          </div>
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

export default RentAndPayment;
