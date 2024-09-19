export async function Login(username: string, password: string): Promise<boolean> {
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
        return true;
    } catch (error) {
        return false;
    }
}

export async function logout(): Promise<boolean> {
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
