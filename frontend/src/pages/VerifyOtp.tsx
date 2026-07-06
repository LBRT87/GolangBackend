import { useState } from "react";
import { req } from "../api/api";

interface VerifyOtpProps {
    email: string;
    onSuccess: () => void;
}

export default function VerifyOtp({ email, onSuccess }: VerifyOtpProps) {
    const [error, setError] = useState('');

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError('');
        try {
            const fields = Object.fromEntries(new FormData(e.currentTarget));
            await req<void>('/auth/register/verify', { email, code: fields.code });
            alert('Registration success');
            onSuccess();
        } catch (err) {
            setError(err instanceof Error ? err.message : "Verification failed")
        }
    }

    return (
        <form onSubmit={handleSubmit} className="w-full max-w-md bg-slate-800 p-6 rounded-2xl shadow-xl space-y-4">
        <h2 className="text-xl font-bold text-center">Verifikasi OTP</h2>
        <p className="text-xs text-slate-400 text-center">Masukkan kode verifikasi yang dikirim ke {email}</p>
        {error && <div className="bg-red-500/20 border border-red-500 text-red-300 p-2.5 rounded-lg text-xs">{error}</div>}
        
        <input type="text" name="code" placeholder="Kode OTP" maxLength={6} required className="w-full text-center bg-slate-950 border border-slate-700 p-2.5 rounded-lg text-lg tracking-widest font-bold text-white focus:outline-none" />
        <button className="w-full bg-emerald-600 hover:bg-emerald-500 py-2.5 rounded-lg text-sm font-semibold text-white">Verifikasi</button>
        </form>
    );
}