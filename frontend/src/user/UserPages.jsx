import { useContext, useEffect, useState } from "react";
import AuthContext from "../auth/AuthContext";
import LoginPage from "../auth/LoginPage";
import { useNavigate } from "react-router";
import "../App.css"
import "./UserPages.css"

// todo: do update user form
export function PersonalAccount() {
    const user = useContext(AuthContext).user;
    const logout = useContext(AuthContext).logout;
    const navigate = useNavigate();

    let patronymic;
    if (!user?.patronymic) {
        patronymic = <i>Отсутствует</i>;
    } else {
        patronymic = user.patronymic;
    }

    if (!user) {
        return <LoginPage/>;
    }

    return (
        <>
        <button className={"form-button"} onClick={logout} style={{right: '20px', top:'20px', position: 'absolute'}}>Выйти</button>
        <div id="main-container">
            <span className={'text-path'} onClick={() => navigate('..')}>Главная</span> 
            <span className={'text-path'}>/</span> 
            <span className={'text-path'} onClick={() => navigate('/personal-account')}>Личный кабинет</span>
            <div id="modal-personal-account">
                <h1>Личный кабинет</h1>
                <div className="field-personal-account">
                    <h2>Имя</h2>{user.first_name}
                </div>
                <div className="field-personal-account">
                    <h2>Фамилия</h2>{user.second_name}
                </div>
                <div className="field-personal-account">
                    <h2>Отчество</h2>{patronymic}
                </div>
                <div className="field-personal-account">
                    <h2>Роль</h2>{user.role}
                </div>
            </div>  
        </div>
        </>
    );
}

export function FreeBookings() {
    const user = useContext(AuthContext).user;
    const logout = useContext(AuthContext).logout;
    const navigate = useNavigate();
    const [isLoading, setIsLoading] = useState(true);
    const [rooms, setRooms] = useState(null);
    const [error, setError] = useState(null);

    useEffect(() => {
        getAllRooms(setIsLoading, setRooms, setError);
    }, []);

    if (!user) {
        return <LoginPage/>;
    }

    return (
        <>
        <button className={"form-button"} onClick={logout} style={{right: '20px', top:'20px', position: 'absolute'}}>Выйти</button>
        <div id="main-container">
            <span className={'text-path'} onClick={() => navigate('..')}>Главная</span> 
            <span className={'text-path'}>/</span> 
            <span className={'text-path'} onClick={() => navigate('/free-bookings')}>Свободные аудитории</span>

            {error && <div className={'modal error'} style={{right: '0', left: '1px'}}>Ошибка<br/>{error}</div>}
            {isLoading && <img src={'/3-dots-loader.svg'} alt={'Загрузка'} className="dots-loader"/>}
            <div className="main-cards">
                {rooms && rooms.map((room) => (
                    <div key={room.id} className="room-card">
                        {room.image && <img src={room.image} alt={room.name} />}
                        {!room.image && <img src="/image_not_found.jpeg" alt={room.name}/>}
                        <h2>{room.name}</h2>
                        <p>{room.description}</p>
                        <p>Вместимость – {room.size}</p>
                    </div>
                ))}
            </div>   
            {rooms && rooms.length === 0 && <div className={'modal error'} style={{right: '0', left: '1px'}}>Ошибка<br/>Нет доступных комнат</div>}
            
        </div>
        </>
    );
}

async function getAllRooms(setIsLoading, setRooms, setError) {
    try {
        const response = await fetch("http://localhost:8080/api/v1/rooms");
        setIsLoading(false);
        if (!response.ok) {
            throw new Error("К сожалению, нет свободных аудиторий на ближайшее время");
        }
        setRooms(await response.json());
    } catch (err) {
        setError("К сожалению, нет свободных аудиторий на ближайшее время");
        setIsLoading(false);
    }          
}

export function MyBookings() {

}