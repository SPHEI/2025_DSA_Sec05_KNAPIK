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

  const [payments, setPayments] = useState([{id: -1, amount: -1, payment_date: '', due_date: '', status_id: -1, renting_id: -1, transaction_reference: '', name: ''}])

  const [names, setNames] = useState(['None'])

  const [role, setRole] = useState('')

  const router = useRouter();

  const [sort, setSort] = useState('Newest')
  const [sort2, setSort2] = useState('Newest')

  const [filter, setFilter] = useState('None')

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
                    const sorted = [...data].sort((a, b) => {
                      const dateA = new Date(a.due_date).getTime();
                      const dateB = new Date(b.due_date).getTime();
                      return dateB - dateA;
                    });

                    setPayments(sorted);

                    const set = new Set(sorted.map(a => a.name))
                    set.add("None")

                    setNames(Array.from(set).reverse())
                    var s = Cookies.get("tFilter");
                    if(s != null)
                    {
                        setFilter(s)
                        Cookies.remove("tFilter")
                    }
                }
              }
              else
              {
                setPayments([])
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

  function countUnique(array: string[])
  {
    return new Set(array).size
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
              <button className="cursor-pointer h-[100%] w-[100%]" onClick={() => {Cookies.set("tSort", "Overdue");router.push("/tenants")}}>
                <div className="flex flex-col items-center justify-center">
                  <img src={card.src} width={40} />
                  
                  <b>Overdue Tenants</b>
                  <h1>{
                    countUnique(payments.filter(a => a.status_id == 3).map(a => a.name))
                  }</h1>
                
                </div>
              </button>
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
          {role === "1" && 
          <div className="white-box w-[50%] py-4">
            <div className="flex flex-row gap-1">
              <h1 className="text-xl">Tenant Filter:</h1>
              <select className="input-box w-[50%]" value={filter} onChange={(a) => {setFilter(a.target.value)}}>
                  {names.map((a, index) => (<option key={index} value={a}>{a}</option>))}
              </select>
            </div>
          </div>
          }

          <div className="white-box w-[50%] py-4">
            <div className="flex flex-col items-left justify-start w-full h-full gap-2">
              <b className="text-xl">Pending Payments</b>
              <div className="flex flex-row gap-1 min-w-[300px]">
                <h1 className="text-xl">Sort:</h1>
                <select className="input-box w-[25%]" value={sort} onChange={(a) => {setSort(a.target.value)}}>
                    <option value="Newest">Newest</option>
                    <option value="Oldest">Oldest</option>
                    <option value="High">Highest Amount</option>
                    <option value="Low">Lowest Amount</option>
                    <option value="Overdue">Overdue</option>
                </select>
              </div>
              <div className="flex flex-col gap-1">
                {(sort == 'Newest' ? payments : sort == "Oldest" ? [...payments].reverse() : [...payments].sort((a, b) => {
                        if(sort == "Overdue")
                        {
                          const priorityA = a.status_id == 3 ? 1 : 0
                          const priorityB = b.status_id == 3 ? 1 : 0
                          return priorityB - priorityA;
                        }
                        else
                        {
                          const priorityA = a.amount 
                          const priorityB = b.amount
                          return sort == "High" ? priorityB - priorityA : priorityA - priorityB;
                        }
                  })).filter((a) => filter == "None" || a.name == filter).map((a,index) => a.status_id != 2 && (
                  <div key={index} className="flex flex-row">
                    <PaymentBox 
                    date={"Due: " + a.due_date.split("T")[0]} 
                    name={role === "1" ? "Tenant: " + a.name  : ""} 
                    type={isOverdue(a.due_date) ? "Overdue" : ""} 
                    amount={a.amount} /> 
                    {role === "2" && <button className="black-button" onClick={()=>{pay(a.id)}}>Pay</button>}
                  </div>))}
              </div>
              <b className="text-xl">Payment history</b>
              <div className="flex flex-row gap-1 min-w-[300px]">
                <h1 className="text-xl">Sort:</h1>
                <select className="input-box w-[25%]" value={sort2} onChange={(a) => {setSort2(a.target.value)}}>
                    <option value="Newest">Newest</option>
                    <option value="Oldest">Oldest</option>
                    <option value="High">Highest Amount</option>
                    <option value="Low">Lowest Amount</option>
                </select>
              </div>
              <div className="flex flex-col gap-1">
                {(sort2 == 'Newest' ? payments : sort2 == "Oldest" ? [...payments].reverse() : [...payments].sort((a, b) => {
                      const priorityA = a.amount 
                      const priorityB = b.amount
                      return sort2 == "High" ? priorityB - priorityA : priorityA - priorityB;
                  })).filter((a) => filter == "None" || a.name == filter).map((a,index) => a.status_id == 2 && 
                <PaymentBox key={index}  
                date={"Due: " + a.due_date.split("T")[0]}  
                name={role === "1" ? "Tenant: " + a.name : ""} 
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
