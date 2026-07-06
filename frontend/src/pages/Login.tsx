import { useState } from "react";
import { req } from "../api/api";
import type { LoginResponse, User } from "../dto/types";

interface LoginProps {
    onLoginSuccess: (user: User) => void;
    onNavigate: (page: 'login' | 'register' | 'otp') => void;
}

export default function Login({ onLoginSuccess, onNavigate }: LoginProps) {
    const [error, setError] = useState('');

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError('');
        const fields = Object.fromEntries(new FormData(e.currentTarget));
        try {
            const data = await req<LoginResponse>('/auth/login', fields);
            
            localStorage.setItem('access_token', data.access_token);
            localStorage.setItem('refresh_token', data.refresh_token);
            localStorage.setItem('user_role', data.user.role);
            
            onLoginSuccess(data.user);
        } catch (err){
            setError(err instanceof Error ? err.message : 'Login failed');
        }
    };

    return (
        <div className="w-full max-w-md bg-slate-800 p-6 rounded-2xl shadow-xl space-y-4">
        <h2 className="text-xl font-bold text-center">Login ESdemy</h2>
        {error && <div className="bg-red-500/20 border border-red-500 text-red-300 p-3 rounded-lg text-xs">{error}</div>}
        
        <form onSubmit={handleSubmit} className="space-y-4">
            <input type="text" name="email" placeholder="Email / Username" required className="w-full bg-slate-950 border border-slate-700 p-2.5 rounded-lg text-sm text-white focus:outline-none focus:border-emerald-500" />
            <input type="password" name="password" placeholder="Password" required className="w-full bg-slate-950 border border-slate-700 p-2.5 rounded-lg text-sm text-white focus:outline-none focus:border-emerald-500" />
            <button className="w-full bg-emerald-600 hover:bg-emerald-500 py-2.5 rounded-lg text-sm font-semibold text-white">Masuk</button>
        </form>
        
        <p className="text-center text-xs text-slate-400">Belum punya akun? <button onClick={() => onNavigate('register')} className="text-emerald-400 hover:underline">Daftar</button></p>
        </div>
    );
}