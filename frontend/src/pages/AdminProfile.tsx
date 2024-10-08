import React, { useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { IsAuthenticated, Logout as LogoutAPI } from "../api/Authentication";

const UserProfile: React.FC = () => {
    const navigate = useNavigate();
    const { username } = useParams();

    useEffect(() => {
        const checkAuth = async () => {
            try {
                const isAuth = await IsAuthenticated();
                if (!isAuth) {
                    navigate("/");
                }
            } catch (error) {
                console.error(error);
            }
        };
        checkAuth();
    }, []);

    const handleLogout = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
    
        LogoutAPI();

        navigate("/");
    }

    return (
        <>
            <h1>admin profile</h1>
            <p>Username: { username }</p>
            <button type="button" onClick={handleLogout}>logout</button>
        </>
    )
}

export default UserProfile;
