"use client";
import "../globals.css";
import { useState, useEffect } from "react";
import card from "./card.png";
import PaymentBox from "../components/PaymentBox";
import calendar from "./calendar.png";

//Shared by admin and tenants
function RentAndPayment() {
  const [ready, setReady] = useState(false);
  const [error, setError] = useState("none");
  useEffect(() => {
    //Page setup goes here
    setReady(true);
  }, []);

  if (ready) {
    if (error == "none") {
      return (
        <main>
          <div className="page-head w-[50%]">
            <b className="text-4xl">Rent & Payments</b>
          </div>
          <div className="flex flex-row gap-4 w-[50%]">
            <div className="white-box h-[150px] w-[100%]">
              <div className="flex flex-col items-center justify-center">
                <img src={card.src} width={40} />
                <b>Current Rent</b>
                <h1>150$</h1>
              </div>
            </div>
            <div className="white-box h-[150px] w-[100%]">
              <div className="flex flex-col items-center justify-center">
                <img src={card.src} width={40} />
                <b>Rent Status</b>
                <h1>paid</h1>
              </div>
            </div>
            <div className="white-box h-[150px] w-[100%]">
              <div className="flex flex-col items-center justify-center">
                <img src={calendar.src} width={40} />
                <b>Next Rent Due</b>
                <h1>28.06.2025</h1>
              </div>
            </div>
          </div>
          <button className="black-button">Pay rent</button>
          <div className="white-box w-[50%] py-4">
            <div className="flex flex-col items-left justify-start w-full h-full gap-2">
              <b className="text-xl">Payment history</b>
              <div className="flex flex-col gap-1">
                <PaymentBox date={"April 30 2025"} type={""} amount={2500} />
                <PaymentBox date={"April 30 2025"} type={""} amount={2500} />
                <PaymentBox date={"April 30 2025"} type={""} amount={2500} />
                <PaymentBox date={"April 30 2025"} type={""} amount={2500} />
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
