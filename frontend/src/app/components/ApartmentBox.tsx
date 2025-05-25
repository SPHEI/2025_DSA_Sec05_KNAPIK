import React from "react";

interface ApartmentProps {
    name: String
  }

function ApartmentBox(props : ApartmentProps)
{
    return(
         <div className="white-box w-[50%] h-[200px]">
            <div className="flex flex-col">
                <h1>Name: {props.name}</h1>
                <h1>Street: default street</h1>
                <h1>Building Number: 000 000 000</h1>
                <h1>Building Name: 1 demo street</h1>
                <h1>Flat Number: 0$</h1>
            </div>
            <button className="black-button">Delete</button>
        </div>
    )
}

ApartmentBox.defaultProps = {
    name: "Default Name"
}
export default ApartmentBox;