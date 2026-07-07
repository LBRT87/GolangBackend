import { useState } from "react";
import Login from "./pages/Login";
import Register from "./pages/Register";
import type { User } from "./dto/types";
import Dashboard from "./pages/Dashboard";
import VerifyOtp from "./pages/VerifyOtp";
import { GoogleOAuthProvider } from "@react-oauth/google";

export default function App() {
  const [page, setPage] = useState<'login' | 'register' | 'otp' | 'dashboard'>('login');
  const [email, setEmail] = useState('');
  const [user, setUser] = useState<User | null>(null);
  
  const handleLogout = () => {
    localStorage.clear();
    setUser(null);
    setPage('login');
  }
  
  return (
    <GoogleOAuthProvider clientId="176129897203-q9gkojlulcvl9ido78qlra5d69tffr4o.apps.googleusercontent.com">
      <div className="min-h-screen bg-slate-900 text-white flex items-center justify-center p-4">
        { page === 'login' && (<Login onLoginSuccess={(u) => { setUser(u); setPage('dashboard'); }} onNavigate={setPage}/>) }
        { page === 'register' && (<Register onNavigate={setPage} onOtpSent={(em) => { setEmail(em); setPage('otp'); }}/>) }
        { page === 'otp' && (<VerifyOtp email={email} onSuccess={ () => setPage('login') }/>) }
        { page === 'dashboard' && user && (<Dashboard user={user} onLogout={handleLogout}/>)}
      </div>
    </GoogleOAuthProvider>
  );
}
