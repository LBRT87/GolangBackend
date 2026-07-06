import type { User } from "../dto/types";

interface DashboardProps {
    user: User,
    onLogout: () => void,
}

export default function Dashboard({ user, onLogout }: DashboardProps) {
    return (
        <div className="w-full max-w-md bg-slate-800 p-6 rounded-2xl shadow-xl space-y-4">
            <div className="flex justify-between items-center">
                <h1 className="text-lg font-bold text-white">Halo, {user.username}</h1>
                <button onClick={onLogout} className="bg-red-600 hover:bg-red-500 px-3 py-1.5 rounded-lg text-xs font-semibold text-white transition-all">
                Keluar
                </button>
            </div>
            
            <div className="text-xs space-y-2 border-t border-slate-700 pt-3 text-slate-400">
                <p>Email: <span className="text-slate-250 font-medium">{user.email}</span></p>
                <p>Role: <span className="text-emerald-450 font-bold uppercase">{user.role}</span></p>
            </div>
        </div>
    );
}