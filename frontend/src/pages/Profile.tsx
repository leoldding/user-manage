import React from "react";
import { useNavigate, useParams } from "react-router-dom";

const Profile: React.FC = () => {
    const navigate = useNavigate();
    const { username } = useParams();

    const handleLogout = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
        navigate("/");
    }

    return (
        <>
            <h1>profile</h1>
            <p>Username: { username }</p>
            <p>Name:</p>
            <p>Role:</p>
            <button type="button" onClick={handleLogout}>logout</button>
        </>
    )
}

export default Profile;
