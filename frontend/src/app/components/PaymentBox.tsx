import React from "react";

interface PaymentProps {
    date: String
    type: String
    name: String
    amount: number
  }

function PaymentBox(props : PaymentProps)
{
    return(
        <div className="white-box-noshadow w-[100%] h-[40px]">
            <div className="flex flex-row gap-8">
                <h1>{props.date.split("T")[0]}</h1>
                <h1>{props.name}</h1>
            </div>
            <div className="flex flex-row gap-8">
                <h1>{props.type}</h1>
                <h1>{"$" + props.amount}</h1>
            </div>
        </div>
    )
}

PaymentBox.defaultProps = {
    date: "Null",
    type: "Loss",
    name: "a",
    amount: 1234,
}
export default PaymentBox;