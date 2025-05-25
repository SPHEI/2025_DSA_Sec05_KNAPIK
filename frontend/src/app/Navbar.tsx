"use client"
import { useState, useEffect } from "react";
import { useRouter, usePathname } from 'next/navigation';
import logo from './logo.png';
import Cookies from "js-cookie";

function Navbar() {
    const [ready,setReady] = useState(false)

    const [userType, setUserType] = useState("not logged in");
    const [currentPage, setCurrentPage] = useState("")

    const router = useRouter();

    //Refresh every time the page changes
    const pathname = usePathname();
    useEffect(refresh, [pathname])
    function refresh()
    {
        setCurrentPage(pathname.substring(1))
        var a = Cookies.get("role");
        if(a != null){
            switch(a){
                case "1":
                    setUserType("admin")
                    break
                case "2":
                    setUserType("tenant")
                    break
                case "3":
                    setUserType("subcontractor")
                    break
            }
        }
        else
        {
            setUserType("not logged in")
        }
        setReady(true)
    }

    async function LogOut()
    {
        //Log out Properly later
        setUserType("not logged in")
        try{
            var t = Cookies.get("token")
            var r = await fetch('http://localhost:8080/logout',{
                method: 'POST',
                body: JSON.stringify({
                    "token" : t 
                })
              });
            //alert(JSON.stringify(await r.json()))
        }catch(err: any){
            alert(err.message)
        }
        Cookies.remove("token")
        Cookies.remove("role")
        router.push("/login")
    }

    function DebugNextUser()
    {
        if      (userType == "tenant") {setUserType("admin")}
        else if (userType == "admin") {setUserType("subcontractor")}
        else if (userType == "subcontractor") {setUserType("not logged in")}
        else {setUserType("tenant")}
    }

    const button =          "h-[35px] px-2 cursor-pointer text-[var(--text)] tracking-wide transition-colors duration-200 rounded-lg hover:bg-[#c0c0c0] active:bg-[#a0a0a0] bg-transparent text-shadow-md "
    const buttonSelected = " h-[35px] px-2 cursor-pointer text-[var(--text)] tracking-wide transition-colors duration-200 rounded-lg bg-[var(--navbutton)] shadow-md"

    if(ready)
    {
        return (
            <nav className="h-[80px] px-4 py-5 bg-[var(--borders)] top-0 w-full shadow-md flex justify-between z-50">
                <div className="flex gap-1">
                    <img className="cursor-pointer" src={logo.src} width={150} onClick={() => router.push("/")}/>
                    <button onClick={DebugNextUser}><b>{userType}</b></button>
                </div>
                {
                    currentPage != "login" ?

                        userType == "tenant" ? 
                            <div className="flex gap-1">

                                <button className={currentPage == "dashboard" ? buttonSelected : button}        onClick={() => router.push("/dashboard")}>
                                    Dashboard
                                </button>

                                <button className={currentPage == "submit-issue" ? buttonSelected : button}     onClick={() => router.push("/submit-issue")}>
                                    Submit Issue
                                </button>

                                <button className={currentPage == "requests" ? buttonSelected : button}         onClick={() => router.push("/requests")}>
                                    My Requests
                                </button>

                                <button className={currentPage == "rent-and-payment" ? buttonSelected : button} onClick={() => router.push("/rent-and-payment")}>
                                    Rent & Payments
                                </button>

                                <button className={button}                                                      onClick={LogOut}>
                                    Log out
                                </button>

                            </div>

                        :userType == "admin" ? 
                            <div className="flex gap-1">

                                <button className={currentPage == "dashboard" ? buttonSelected : button}        onClick={() => router.push("/dashboard")}>
                                    Dashboard
                                </button>

                                <button className={currentPage == "tenants" ? buttonSelected : button}          onClick={() => router.push("/tenants")}>
                                    Tenants
                                </button>

                                <button className={currentPage == "requests" ? buttonSelected : button}         onClick={() => router.push("/requests")}>
                                    Requests
                                </button>

                                <button className={currentPage == "rent-and-payment" ? buttonSelected : button} onClick={() => router.push("/rent-and-payment")}>
                                    Payments
                                </button>

                                <button className={currentPage == "reports" ? buttonSelected : button}          onClick={() => router.push("/reports")}>
                                    Reports
                                </button>

                                <button className={currentPage == "accounts" ? buttonSelected : button}         onClick={() => router.push("/accounts")}>
                                    Add Users
                                </button>

                                <button className={button} onClick={LogOut}>
                                    Log out
                                </button>

                            </div>

                        :userType == "subcontractor" ?
                            <div className="flex gap-1">

                                <button className={currentPage == "dashboard" ? buttonSelected : button}        onClick={() => router.push("/dashboard")}>
                                    Dashboard
                                </button>

                                <button className={currentPage == "requests" ? buttonSelected : button}         onClick={() => router.push("/requests")}>
                                    My Requests
                                </button>

                                <button className={button}                                                      onClick={LogOut}>
                                    Log out
                                </button>
                                
                            </div>

                        :<button className={button} onClick={() => router.push("/login")}>Log In</button>
                    :<h1></h1>
                }
            </nav>
        );
    }
};

export default Navbar;
