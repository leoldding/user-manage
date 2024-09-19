import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

const Login: React.FC = () => {
    const navigate = useNavigate();
    const [credentials, setCredentials] = useState<{ Username: string, Password: string }>({ Username: "", Password: "" });

    const handleCredentialChange = (field: string, value: string) => {
        setCredentials((prevState) => ({
            ...prevState,
            [field]: value
        }));
    }

    const enableLoginButton = () => {
        const { Username, Password } = credentials || {};
        const isUsernameValid = /^[a-zA-Z\d]+$/.test(Username);
        const isPasswordValid = /^[a-zA-Z\d]+$/.test(Password);
        return isUsernameValid && isPasswordValid;
    }

    const handleLogin = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
        const response = fetch("/ping")
        console.log(response)
        navigate("/u/" + credentials.Username);
    }

    return (
        <>
            <h1>login</h1>
            <form>
                <input
                    type="text"
                    placeholder="username"
                    value={credentials.Username}
                    onChange={(e) => handleCredentialChange("Username", e.target.value)}
                />
                <input
                    type="password"
                    placeholder="password"
                    value={credentials.Password}
                    onChange={(e) => handleCredentialChange("Password", e.target.value)}
                />
                <button type="submit" disabled={!enableLoginButton()} onClick={handleLogin}>login</button>
            </form>
        </>
    )
}

export default Login;
