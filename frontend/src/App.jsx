import './App.css'
import AuthProvider from "./auth/AuthProvider.jsx";
import {BrowserRouter, Route, Routes} from "react-router";
import Signup from "./auth/Signup.jsx";
import React from "react";
import LoginPage from "./auth/LoginPage.jsx";
import Navigation from "./Navigation.jsx"

function App() {
    return (
        <BrowserRouter>
            <AuthProvider>
                <img src={"../resources/Central_University_full.png"}
                     id={"logo"} alt={"Логотип Центрального университета"} key={"cu_logo"}/>
                <Routes>
                    <Route path={"/login"} element={<LoginPage/>}/>
                    <Route path={"/signup"} element={<Signup/>}/>
                    <Route path={"*"} element={<Navigation/>}/>
                </Routes>
            </AuthProvider>
        </BrowserRouter>
    );
}

export default App
