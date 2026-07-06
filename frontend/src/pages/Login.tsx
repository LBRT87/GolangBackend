interface LoginProps {
    onLoginSuccess: (user: any) => void;
    onNavigate: (page: 'login' | 'register' | 'otp') => void;
}

export default function Login({ onLoginSuccess, onNavigate }: LoginProps) {
    const [error, setError] = useState('');
}