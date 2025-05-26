"use client"
import "../globals.css";
import { useState, useEffect } from "react";
import { useRouter, usePathname } from 'next/navigation';
import Cookies from "js-cookie";

//Admin only
function Dashboard() {
    const [ready,setReady] = useState(false)
    const [error, setError] = useState('none')

    const [apartaments,setApartaments] = useState([''])
    const [specialities,setSpecialities] = useState(['a','b'])

    const router = useRouter();

    const pathname = usePathname();

    useEffect(() => {
        const fetchData = async () => {
        //Page setup goes here
        var a = Cookies.get("role");
        if(a == null )
        {
            alert("You need to be logged in to see this page.");
            router.push("/login");
            return
        }
        if(a != "1")
        {
            alert("You need to be an admin to see this page.");
            router.push("/dashboard");
            return
        }

        var t = Cookies.get("token");
        try{
        const res = await fetch('http://localhost:8080/apartaments',{
                method:'POST',
                body: JSON.stringify({ 
                    "token": t,
                })
            });
            const data = await res.json();
            if(data.message)
            {
                setError(data.message)
            }
            else
            {
                setApartaments(data.apartaments);
            }
        }
        catch(err: any)
        {
            setError(err.message);
        }

        setReady(true);
        }
        fetchData();
    },[pathname])

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

    async function sendData()
    {
        var a = 0;
        var t = Cookies.get("token");
        if(password != repassword)
        {
            alert("Passwords don't match")
            return
        }
        if(t == null)
        {
            alert("Token not found.");
            return;
        }
        if (role === 'Admin'){a = 1}
        if (role === 'Tenant'){a = 2}
        if (role === 'Subcontractor'){a = 3}
        try {
            const res = await fetch('http://localhost:8080/adduser',{
                method:'POST',
                body: JSON.stringify({ 
                    "token": t,
                    name, 
                    password,
                    email,
                    phone,
                    apartment,
                    rent,
                    address,
                    nip,
                    speciality,
                    "role":a
                })
            });
            const data = await res.json();
            if(data.message)
            {
                alert(data.message)
            }
            else
            {
                alert("User added succesfully.");
            }
        } catch (err: any) {
            setError(err.message)
        } finally{
            setReady(true);
        }
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
                            <select className="input-box w-[26%]" defaultValue={'Tenant'} onChange={(a) => {setRole(a.target.value)}}>
                                <option value="Admin">Admin</option>
                                <option value="Tenant">Tenant</option>
                                <option value="Subcontractor">Subcontractor</option>
                            </select>
                            <div className={line}>
                                <b className="w-[26%]">Name</b>
                                <b className="w-[26%]">Email</b>
                                <b className="w-[26%]">Phone</b>
                            </div>
                            <div className={line}>
                                <input className="input-box" placeholder="Name" value={name} onChange={(a) => {setName(a.target.value)}}/>
                                <input className="input-box" placeholder="E-mail" value={email} onChange={(a) => {setEmail(a.target.value)}} />
                                <input type="tel" className="input-box" placeholder="Phone" pattern="[0-9]*" value={phone} onChange={phoneChange}/>
                            </div>
                            <div className={line}>
                                <b className="w-[26%]">Password</b>
                                <b className="w-[26%]">Repeat password</b>
                            </div>
                            <div className={line}>
                                <input className="input-box" placeholder="Password" type="password" value={password} onChange={(a) => {setPassword(a.target.value)}}/>
                                <input className="input-box" placeholder="Repeat Password" type="password" value={repassword} onChange={(a) => {setRepassword(a.target.value)}}/>
                            </div>
                            {role === "Tenant" && (
                                <div>
                                    <div className={line}>
                                        <b className="w-[26%]">Apartment</b>
                                        <b className="w-[26%]">Rent</b>
                                    </div>
                                    <div className={line}>
                                        <select className="input-box w-[26%]" onChange={(a) => {setApartment(a.target.value)}}>
                                            {apartaments.map((a,index) => (<option key={index} value={a}>{a}</option>))}
                                        </select>
                                        <input className="input-box" placeholder="Rent" value={rent} onChange={rentChange}/>
                                    </div>
                                </div>
                            )}
                            {role === "Subcontractor" && (
                                <div>
                                    <div className={line}>
                                        <b className="w-[26%]">Address</b>
                                        <b className="w-[26%]">NIP</b>
                                        <b className="w-[26%]">Speciality</b>
                                    </div>
                                    <div className={line}>
                                        <input className="input-box" placeholder="Address" value={address} onChange={(a) => {setAddress(a.target.value)}}/>
                                        <input className="input-box" placeholder="NIP" value={nip} onChange={nipChange}/>
                                        <select className="input-box w-[26%]" onChange={(a) => {setSpeciality(a.target.value)}}>
                                            {specialities.map((a,index) => (<option key={index} value={a}>{a}</option>))}
                                        </select>
                                    </div>
                                </div>
                            )}
                            <button className="black-button w-[26%]" onClick={sendData}>Add</button>
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
