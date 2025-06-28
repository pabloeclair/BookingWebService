import "../App.css"
import {useNavigate} from "react-router";
import {useContext} from "react";
import AuthContext from "../auth/AuthContext.jsx";
import LoginPage from "../auth/LoginPage.jsx";

export function MainAdminPage() {

    const user = useContext(AuthContext).user;
    const logout = useContext(AuthContext).logout;
    const navigate = useNavigate();

    if (!user) {
        return <LoginPage/>;
    }

    return (
        <>
        <button className={"form-button"} onClick={logout} style={{right: '20px', top:'20px', position: 'absolute'}}>Выйти</button>
        <div id={'main-container'}>
            <span className={'text-path'} onClick={() => navigate('/')}>Главная</span>
            <h1>Добро пожаловать,<br/>{user.first_name}</h1><br/>
            <div className={'buttons-container'}>
                <button onClick={() => navigate('/personal-account')} className={'main-cards button-navigate'}>👤 Личный кабинет</button>
                <button onClick={() => navigate('/admin/users')} className={'main-cards button-navigate'}>👥 Пользователи</button>
                <button onClick={() => navigate('/admin/places')} className={'main-cards button-navigate'}>🏢 Места бронирования</button>
                <button onClick={() => navigate('/admin/bookings')} className={'main-cards button-navigate'}>📅 Записи бронирования</button>
            </div>
        </div>
        </>
    );
}

export function MainUserPage() {
    const user = useContext(AuthContext).user;
    const logout = useContext(AuthContext).logout;
    const navigate = useNavigate();

    if (!user) {
        return <LoginPage/>;
    }

    return (
        <>
        <button className={"form-button"} onClick={logout} style={{right: '20px', top:'20px', position: 'absolute'}}>Выйти</button>
        <div id={'main-container'}>
            <span className={'text-path'}>Главная</span>
            <h1>Добро пожаловать,<br/>{user.first_name}</h1><br/>
            <div className={'buttons-container'}>
                <button onClick={() => navigate('/personal-account')} className={'main-cards button-navigate'}>👤 Личный кабинет</button>
                <button onClick={() => navigate('/my-bookings')} className={'main-cards button-navigate'}>📅 Мои записи</button>
                <button onClick={() => navigate('/free-bookings')} className={'main-cards button-navigate'}>🏢 Свободные аудитории</button>
            </div>
        </div>
        </>
    );
}