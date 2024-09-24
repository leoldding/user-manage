import React from "react";
import { BrowserRouter as Router, Route, Routes, Navigate } from "react-router-dom"; 
import Login from "./pages/Login";
import UserProfile from "./pages/UserProfile";
import AdminProfile from "./pages/AdminProfile";
import NotFound from "./pages/NotFound";

const App: React.FC = () => {
    return (
        <>
            <Router>
                <Routes>
                    <Route path="/" element={<Login />} />
                    <Route path="/u/:username" element={<UserProfile />} />
                    <Route path="a/:username" element={<AdminProfile />} />
                    <Route path="/404" element={<NotFound />} />
                    <Route path="*" element={<Navigate to="/404" />} />
                </Routes>
            </Router>
        </>
    )
}

export default App
