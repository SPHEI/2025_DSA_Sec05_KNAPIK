"use client"
import "./globals.css";

import { useRouter } from 'next/navigation';
import { useEffect } from "react";
//App start - Checks auth and redirects either to login or dashboard
function App() {
  const router = useRouter();
  useEffect(() => {
    router.push("/dashboard")
})
}

export default App;