import { useState } from "react";
import { req } from "../api/api";
import type { RegisterResponse } from "../dto/types";

interface RegisterProps {
    onOtpSent: (email: string) => void;
    onNavigate: (page: 'login' | 'register' | 'otp') => void;
}

export default function Register({ onNavigate, onOtpSent }: RegisterProps){
    const [error, setError] = useState('');

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError('');
        const fields = Object.fromEntries(new FormData(e.currentTarget));
        console.log(fields)
        try {
            await req<RegisterResponse>('auth/register', fields);
            onOtpSent(fields.email as string);
            onNavigate('otp');
        } catch (err) {
            setError(err instanceof Error ? err.message : "Register failed");
        }
    }

    return (
        <div className="w-full max-w-md bg-slate-800 p-6 rounded-2xl shadow-xl space-y-4">
            <h2 className="text-xl font-bold text-center">Daftar Akun</h2>
            {error && <div className="bg-red-500/20 border border-red-500 text-red-300 p-3 rounded-lg text-xs">{error}</div>}
            
            <form onSubmit={handleSubmit} className="space-y-4">
                <input type="text" name="username" placeholder="Username" minLength={3} required className="w-full bg-slate-950 border border-slate-700 p-2.5 rounded-lg text-sm text-white focus:outline-none focus:border-emerald-500" />
                <input type="email" name="email" placeholder="Email" required className="w-full bg-slate-950 border border-slate-700 p-2.5 rounded-lg text-sm text-white focus:outline-none focus:border-emerald-500" />
                <input type="password" name="password" placeholder="Password" minLength={6} required className="w-full bg-slate-950 border border-slate-700 p-2.5 rounded-lg text-sm text-white focus:outline-none focus:border-emerald-500" />
                <input type="date" name="dob" required className="w-full bg-slate-950 border border-slate-700 p-2.5 rounded-lg text-sm text-white focus:outline-none focus:border-emerald-500" />
                <button className="w-full bg-emerald-600 hover:bg-emerald-500 py-2.5 rounded-lg text-sm font-semibold text-white">Daftar</button>
            </form>
            
            <p className="text-center text-xs text-slate-400">Sudah punya akun? <button onClick={() => onNavigate('login')} className="text-emerald-400 hover:underline">Login</button></p>
        </div>
    );
}