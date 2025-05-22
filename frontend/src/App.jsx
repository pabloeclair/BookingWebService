import './App.css'
import AuthProvider from "./auth/AuthProvider.jsx";
import {BrowserRouter, Route, Routes} from "react-router";
import Signup from "./auth/Signup.jsx";
import React from "react";
import Login from "./auth/Login.jsx";

function App() {
    return (
        <AuthProvider>
            <BrowserRouter>
                <img src={"../resources/Central_University_full.png"}
                     id={"logo"} alt={"Логотип Центрального университета"} key={"cu_logo"}/>
                <Routes>
                    <Route path={"/user/login"} element={<Login/>}/>
                    <Route path={"user/signup"} element={<Signup/>}/>
                </Routes>
            </BrowserRouter>
        </AuthProvider>
    );
}

export default App
