import React from "react";

interface PaymentProps {
    date: String
    type: String
    amount: number
  }

function PaymentBox(props : PaymentProps)
{
    return(
        <div className="white-box-noshadow w-[100%] h-[40px]">
            <h1>{props.date}</h1>
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
    amount: 1234
}
export default PaymentBox;