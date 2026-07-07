import { useState } from "react";
import { req } from "../api/api";
import type { LoginResponse, User } from "../dto/types";
import { useGoogleLogin } from "@react-oauth/google";

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

    const loginGoogle = useGoogleLogin({
        onSuccess: async (tokenResponse) => {
            setError('');
            try{
                const accessToken = tokenResponse.access_token;
                const data = await req<LoginResponse>('/auth/google/login', { token: accessToken });

                localStorage.setItem('access_token', data.access_token);
                localStorage.setItem('refresh_token', data.refresh_token);
                localStorage.setItem('user_role', data.user.role)
            } catch (err){
                setError(err instanceof Error ? err.message : 'Google login failed')
            }
        },
        onError: () => setError('Google login failed 2')
    });

    return (
        <div className="w-full max-w-md bg-slate-800 p-6 rounded-2xl shadow-xl space-y-4">
        <h2 className="text-xl font-bold text-center">Login ESdemy</h2>
        {error && <div className="bg-red-500/20 border border-red-500 text-red-300 p-3 rounded-lg text-xs">{error}</div>}
        
        <form onSubmit={handleSubmit} className="space-y-4">
            <input type="text" name="email" placeholder="Email / Username" required className="w-full bg-slate-950 border border-slate-700 p-2.5 rounded-lg text-sm text-white focus:outline-none focus:border-emerald-500" />
            <input type="password" name="password" placeholder="Password" required className="w-full bg-slate-950 border border-slate-700 p-2.5 rounded-lg text-sm text-white focus:outline-none focus:border-emerald-500" />
            <button className="w-full bg-emerald-600 hover:bg-emerald-500 py-2.5 rounded-lg text-sm font-semibold text-white">Masuk</button>
        </form>
        <button className="relative inline-flex w-full items-center justify-center rounded-md border border-gray-400 bg-white px-3.5 py-2.5 font-semibold text-gray-700 transition-all duration-200 hover:bg-gray-100 hover:text-black focus:bg-gray-100 focus:text-black focus:outline-none" type="button" onSubmit={() => loginGoogle}>
            <span className="mr-2 inline-block">
                <svg fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 text-emerald-600">
                <path d="M20.283 10.356h-8.327v3.451h4.792c-.446 2.193-2.313 3.453-4.792 3.453a5.27 5.27 0 0 1-5.279-5.28 5.27 5.27 0 0 1 5.279-5.279c1.259 0 2.397.447 3.29 1.178l2.6-2.599c-1.584-1.381-3.615-2.233-5.89-2.233a8.908 8.908 0 0 0-8.934 8.934 8.907 8.907 0 0 0 8.934 8.934c4.467 0 8.529-3.249 8.529-8.934 0-.528-.081-1.097-.202-1.625z"/>
                </svg>
            </span>
            Sign in with Google
        </button>
        <p className="text-center text-xs text-slate-400">Belum punya akun? <button onClick={() => onNavigate('register')} className="text-emerald-400 hover:underline">Daftar</button></p>
        </div>
    );
}