import './App.css'
import AuthProvider from "./auth/AuthProvider.jsx";
import {BrowserRouter, Route, Routes} from "react-router";
import LoginPage from "./auth/LoginPage.jsx";
import SignupPage from './auth/SignupPage.jsx';
import {MainAdminPage, MainUserPage} from './user/MainPages.jsx';
import ErrorNotFound from './ErrorNotFound.jsx';
import { AllRooms, CreateBooking, MyBookings, PersonalAccount } from './user/UserPages.jsx';

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
                    <Route path={"/"} element={<MainUserPage/>}/>
                    <Route path={"/personal-account"} element={<PersonalAccount/>}/>
                    <Route path={"/rooms"} element={<AllRooms/>}/>
                    <Route path={"/rooms/:id"} element={<CreateBooking/>}/>
                    <Route path={"/my-bookings"} element={<MyBookings/>}/>
                    <Route path={"*"} element={<ErrorNotFound/>}/>
                </Routes>
            </AuthProvider>
        </BrowserRouter>
    );
}

export default App
