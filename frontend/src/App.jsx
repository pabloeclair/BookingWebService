import './App.css'
import AuthProvider from "./auth/AuthProvider.jsx";
import {BrowserRouter, Route, Routes} from "react-router";
import LoginPage from "./auth/LoginPage.jsx";
import Navigation from "./Navigation.jsx"
import SignupPage from './auth/SignupPage.jsx';
import MainAdminPage from './user/MainAdminPage.jsx';

function App() {
    return (
        <BrowserRouter>
            <AuthProvider>
                <img src={"/Central_University_full.png"}
                     id={"logo"} alt={"Логотип Центрального университета"} key={"cu_logo"}/>
                <Routes>
                    <Route path={"/login"} element={<LoginPage/>}/>
                    <Route path={"/signup"} element={<SignupPage/>}/>
                    <Route path={"/admin"} element={<MainAdminPage/>}/>
                    <Route path={"*"} element={<Navigation/>}/>
                </Routes>
            </AuthProvider>
        </BrowserRouter>
    );
}

export default App
