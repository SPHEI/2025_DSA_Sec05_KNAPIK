import React from "react";

interface RequestProps {
  title: string;
  date: string;
  status: string;
  colorClass: string; //do statusu z tych boxow green itd
}

function RequestBox(props: RequestProps) {
  return (
    <div className="white-box w-[50%] h-[200px]">
      <div className="flex flex-row items-center justify-between gap-8 w-full px-8">
        <div className="flex flex-col">
          <h1 className="text-xl">Title: {props.title}</h1>
          <h1 className="text-xl">Date submitted: {props.date}</h1>
          <h1 className="text-xl">Status: {props.status}</h1>
          <button className="black-button">View Details</button>
        </div>
        <div className={props.colorClass}></div>
      </div>
    </div>
  );
}

RequestBox.defaultProps = {
  title: "Electricity",
  date: "30 April 2137",
  status: "Pending",
  colorClass: "status-box-yellow",
};

export default RequestBox;
