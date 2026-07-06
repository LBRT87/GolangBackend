export interface User {
    id: number,
    username: string,
    email: string,
    role: 'student' | 'lecturer',
    dob: string
}

export interface LoginResponse {
    access_token: string,
    refresh_token: string,
    user: User
}

export interface RegisterResponse {
    message: string
}