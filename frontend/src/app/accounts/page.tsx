"use client"
import "../globals.css";
import { useState, useEffect } from "react";
import { useRouter, usePathname } from 'next/navigation';
import Cookies from "js-cookie";

//Admin only
function Dashboard() {
    const [ready,setReady] = useState(false)
    const [error, setError] = useState('none')

    const [apartaments,setApartaments] = useState([{id: -1,name: '', street: '', building_number: '', building_name: '',flat_number:'',owner_id:-1 }])
    const [specialities,setSpecialities] = useState([{id: -1, name: ''}])

    const router = useRouter();

    const pathname = usePathname();

    const [showPopup,setShowPopup] = useState(false)


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
        const res = await fetch('http://localhost:8080/apartament/list?token=' + t)
            const data = await res.json();
            if(data.message)
            {
                setError(data.message)
            }
            else
            {
                setApartaments(data);
            }
        }
        catch(err: any)
        {
            setError(err.message);
        }
        await reFetchSpecialities();

        setReady(true);
        }
        fetchData();
    },[pathname])

    async function reFetchSpecialities()
    {
        var t = Cookies.get("token");
        try{
        const res = await fetch('http://localhost:8080/subspec',{
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
                //alert(JSON.stringify(data))
                setSpecialities(data);
            }
        }
        catch(err: any)
        {
            setError(err.message);
        }
    }

    const [role, setRole] = useState('Tenant')

    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [phone, setPhone] = useState('');
    const [password, setPassword] = useState('');
    const [repassword, setRepassword] = useState('');
    const [apartment, setApartment] = useState(1);
    const [date, setDate] = useState('');
    const [address, setAddress] = useState('');
    const [nip, setNip] = useState('');
    const [speciality, setSpeciality] = useState(1);

    const [newSpeciality, setNewSpeciality] = useState('');

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


    const formatNip = (value: string) => {
        const digits = value.replace(/\D/g, '').substring(0, 10);
        return digits;
    }
    function nipChange(e: React.ChangeEvent<HTMLInputElement>){
        const a = formatNip(e.target.value);
        setNip(a);
    }

    async function addSpeciality()
    {
        var t = Cookies.get("token");
        try {
            const res = await fetch('http://localhost:8080/addsubspec',{
                method:'POST',
                body: JSON.stringify({ 
                    "token" : t,
                    "name" : newSpeciality
                })
            });
            if(res.ok)
            {
                alert("Speciality added succesfully.");
            }
            else
            {
                var data = await res.json()
                alert(data.message)
            }
            reFetchSpecialities()
            setShowPopup(false)
        } catch (err: any) {
            setError(err.message)
        } finally{
            setReady(true);
        }
    }

    function clear()
    {
        setName('')
        setEmail('')
        setPhone('')
        setPassword('')
        setRepassword('')
        setApartment(1)
        setDate('')
        setAddress('')
        setNip('')
        setSpeciality(1)
        setNewSpeciality('')
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
        // alert(JSON.stringify({ 
        //             "token": t,
        //             name, 
        //             password,
        //             email,
        //             phone,
        //             apartment,
        //             rent,
        //             address,
        //             nip,
        //             speciality,
        //             "role":a
        //         }))
        try {
            const res = await fetch('http://localhost:8080/adduser',{
                method:'POST',
                body: JSON.stringify({ 
                    "token": t,
                    "user" :{
                        "name": name, 
                        "password": password,
                        "email" : email,
                        "phone" : phone,
                        "role_id":a
                    }
                })
            });
            const data = await res.json();
            //alert(JSON.stringify(data))
            if(data.message)
            {
                alert(data.message)
            }
            else
            {
                if(role === 'Tenant')
                {
                    const res2 = await fetch('http://localhost:8080/renting/start',{
                        method:'POST',
                        body: JSON.stringify({ 
                            "token": t,
                            "renting" : {
                                "apartment_id" : apartment,
                                "user_id" : data.id,
                                "start_date" : date + "T00:00:00Z"
                            }
                        })
                    });
                    if(res2.ok)
                    {
                        alert("Tenant Added Succesfully")
                        clear()
                    }
                    else
                    {
                        var data2 = await res2.json()
                        alert(data2.message)
                    }
                }
                else if(role === 'Subcontractor')
                {

                    const res2 = await fetch('http://localhost:8080/subcon/add',{
                        method:'POST',
                        body: JSON.stringify({ 
                            "token": t,
                            "subcontractor":{
                                "user_id" : data.id,
                                "address" : address,
                                "NIP" : nip,
                                "speciality_id" : speciality
                            }
                        })
                    });
                    if(res2.ok)
                    {
                        alert("Subcontractor Added Succesfully")
                        clear()
                    }
                    else
                    {
                        var data2 = await res2.json()
                        alert(data2.message)
                    }
                }
                else
                {
                    alert("Admin added succesfully.");
                    clear()
                }
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
                    <div className="page-head w-[55%] min-w-[500px]">
                        <b className="text-4xl">Add Users</b> 
                    </div>
                    <div className="white-box w-[55%] py-8 min-w-[500px]">
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
                                <input className="input-box w-[26%]" placeholder="Name" value={name} onChange={(a) => {setName(a.target.value)}}/>
                                <input className="input-box w-[26%]" placeholder="E-mail" value={email} onChange={(a) => {setEmail(a.target.value)}} />
                                <input className="input-box w-[26%]" placeholder="Phone" value={phone} onChange={phoneChange}/>
                            </div>
                            <div className={line}>
                                <b className="w-[26%]">Password</b>
                                <b className="w-[26%]">Repeat password</b>
                            </div>
                            <div className={line}>
                                <input className="input-box w-[26%]" placeholder="Password" type="password" value={password} onChange={(a) => {setPassword(a.target.value)}}/>
                                <input className="input-box w-[26%]" placeholder="Repeat Password" type="password" value={repassword} onChange={(a) => {setRepassword(a.target.value)}}/>
                            </div>
                            {role === "Tenant" && (
                                <div>
                                    <div className={line}>
                                        <b className="w-[26%]">Apartment</b>
                                        <b className="w-[26%]">Start Date</b>
                                    </div>
                                    <div className={line}>
                                        <select className="input-box w-[26%]" value={apartment} onChange={(a) => {setApartment(Number(a.target.value))}}>
                                            {apartaments.map((a,index) => (<option key={index} value={a.id}>{a.name}</option>))}
                                        </select>
                                        <input
                                        type="date"
                                        className="input-box w-[26%]"
                                        value={date}
                                        onChange={(e) => setDate(e.target.value)}
                                        />
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
                                        <input className="input-box w-[26%]" placeholder="Address" value={address} onChange={(a) => {setAddress(a.target.value)}}/>
                                        <input className="input-box w-[26%]" placeholder="NIP" value={nip} onChange={nipChange}/>
                                        <select className="input-box w-[26%]" onChange={(a) => {setSpeciality(Number(a.target.value))}}>
                                            {specialities.map((a,index) => (<option key={index} value={a.id}>{a.name}</option>))}
                                        </select>
                                    </div>
                                </div>
                            )}
                            <div className={line}>
                                <button className="black-button w-[26%]" onClick={sendData}>Add</button>
                                {role === "Subcontractor" && (<button className="black-button w-[26%] relative left-[26.5%]" onClick={() => {setShowPopup(true)}}>+ Add Speciality</button>)}
                            </div>
                        </div>
                    </div>
                    {showPopup && (
                    <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/50">
                        <div className="white-box w-[20%] py-4 rounded-lg relative min-w-[400px]">
                            <div className="flex flex-col gap-2 w-[100%]">
                                <b className="text-4xl">Add Speciality</b>
                                <input className="input-box" placeholder="Speciality Name" onChange={(a)=>{setNewSpeciality(a.target.value)}}/>
                                <button className="black-button" onClick={addSpeciality}>Add</button>
                            </div>
                            <button onClick={() => setShowPopup(false)}className="absolute top-4 right-4 text-xl font-bold cursor-pointer">x</button>
                        </div>
                    </div>
                    )}
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
