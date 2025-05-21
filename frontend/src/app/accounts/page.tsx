"use client"
import { format } from "path";
import "../globals.css";
import { useState, useEffect } from "react";

//Admin only
function Dashboard() {
    const [ready,setReady] = useState(false)
    const [error, setError] = useState('none')
    useEffect(() => {
        //Page setup goes here
        setReady(true);
    },[])

    const [role, setRole] = useState('Tenant')

    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [phone, setPhone] = useState('');
    const [password, setPassword] = useState('');
    const [repassword, setRepassword] = useState('');
    const [apartment, setApartment] = useState('');
    const [rent, setRent] = useState('');
    const [address, setAddress] = useState('');
    const [nip, setNip] = useState('');
    const [speciality, setSpeciality] = useState('');

    function formatPhone(value: string) {
        let cleaned = value.replace(/[^\d+\s]/g, '');
        if (cleaned.startsWith('+')) {
            cleaned = '+' + cleaned.slice(1).replace(/[+]/g, '');
        } else {
            cleaned = cleaned.replace(/[+]/g, '');
        }
        return cleaned;
    };
    function phoneChange(e: React.ChangeEvent<HTMLInputElement>){
        const formatted = formatPhone(e.target.value);
        setPhone(formatted);
    };

    const formatRent = (value: string) => {
        const digits = value.replace(/[^\d.]/g, '');
        return digits;
    }
    function rentChange(e: React.ChangeEvent<HTMLInputElement>){
        const a = formatRent(e.target.value);
        setRent(a);
    }


    const formatNip = (value: string) => {
        const digits = value.replace(/\D/g, '').substring(0, 10);
        return digits;
    }
    function nipChange(e: React.ChangeEvent<HTMLInputElement>){
        const a = formatNip(e.target.value);
        setNip(a);
    }

    function sendData()
    {
        var a = role + "\n\nName: " + name + "\nE-mail: " + email + "\nPhone number: " + phone + "\nPassword: " + password + "\nRepeated password: " + repassword + "\n";
        if(role === "Tenant")
        {
            a += "Apartment: " + apartment + "\nRent: " + rent + "\n"; 
        } 
        if(role === "Subcontractor")
        {
            a += "Address: " + address + "\nNIP: " + nip + "\nSpeciality: " + speciality;
        }
        alert(a)
    }

    const line = "flex flex-row gap-1";
    if(ready)
    {
        if(error == 'none')
        {
            return (
                <main>
                    <div className="page-head w-[55%]">
                        <b className="text-4xl">Add Users</b> 
                    </div>
                    <div className="white-box w-[55%] py-8">
                        <div className="flex flex-col gap-2 w-[100%]">
                            <select className="input-box w-[25.5%]" defaultValue={'Tenant'} onChange={(a) => {setRole(a.target.value)}}>
                                <option value="Tenant">Tenant</option>
                                <option value="Subcontractor">Subcontractor</option>
                                <option value="Admin">Admin</option>
                            </select>
                            <div className={line}>
                                <input className="input-box" placeholder="Name" value={name} onChange={(a) => {setName(a.target.value)}}/>
                                <input className="input-box" placeholder="E-mail" value={email} onChange={(a) => {setEmail(a.target.value)}} />
                                <input type="tel" className="input-box" placeholder="Phone" pattern="[0-9]*" value={phone} onChange={phoneChange}/>
                            </div>
                            <div className={line}>
                                <input className="input-box" placeholder="Password" type="password" value={password} onChange={(a) => {setPassword(a.target.value)}}/>
                                <input className="input-box" placeholder="Repeat Password" type="password" value={repassword} onChange={(a) => {setRepassword(a.target.value)}}/>
                            </div>
                            {role === "Tenant" && (
                                <div className={line}>
                                    <input className="input-box" placeholder="Apartment" value={apartment} onChange={(a) => {setApartment(a.target.value)}}/>
                                    <input className="input-box" placeholder="Rent" value={rent} onChange={rentChange}/>
                                </div>
                            )}
                            {role === "Subcontractor" && (
                                <div className={line}>
                                    <input className="input-box" placeholder="Address" value={address} onChange={(a) => {setAddress(a.target.value)}}/>
                                    <input className="input-box" placeholder="NIP" value={nip} onChange={nipChange}/>
                                    <input className="input-box" placeholder="Speciality" value={speciality} onChange={(a) => {setSpeciality(a.target.value)}}/>
                                </div>
                            )}
                            <button className="black-button w-[25.5%]" onClick={sendData}>Add</button>
                        </div>
                    </div>
                </main>
            );
        }
        else
        {
            return (
                <main>
                    <b>An error has occured:</b>
                    <h1>{error}</h1>
                </main>
            );
        }
    }
    else
    {
        return(
            <main>
                <h1>Loading...</h1>
            </main>
        )
    }
};

export default Dashboard;
