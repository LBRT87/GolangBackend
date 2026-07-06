const BASE_URL = 'http://localhost/api';

export const req = async <T> (path: string, body: object): Promise<T> => {
    const res = await fetch(`${BASE_URL}${path}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json '},
        body: JSON.stringify(body)
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "Error");
    return data as T;
};