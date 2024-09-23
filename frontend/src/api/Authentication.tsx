export async function Login(username: string, password: string): Promise<string> {
    try {
        const response = await fetch("/api/login", {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({Username: username, Password: password}),
        });
        if (!response.ok) {
            throw new Error("Unable to login");
        }
        const role: string = await response.text();
        return role;
    } catch (error) {
        return "";
    }
}

export async function Logout(): Promise<boolean> {
    try {
        const response = await fetch("/api/logout", {
            method: "GET",
            credentials: "include",
        });
        if (!response.ok) {
            throw new Error("Unable to logout");
        }
        return true;
    } catch (error) {
        return false;
    }
}
