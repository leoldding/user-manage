import React from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Logout as LogoutAPI } from "../api/Authentication";

const UserProfile: React.FC = () => {
    const navigate = useNavigate();
    const { username } = useParams();

    const handleLogout = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
    
        LogoutAPI();

        navigate("/");
    }

    return (
        <>
            <h1>user profile</h1>
            <p>Username: { username }</p>
            <button type="button" onClick={handleLogout}>logout</button>
        </>
    )
}

export default UserProfile;
