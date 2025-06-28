import { useContext, useEffect, useState } from "react";
import AuthContext from "../auth/AuthContext";
import LoginPage from "../auth/LoginPage";
import { useNavigate, useParams } from "react-router";
import "../App.css"
import "../styles/UserPages.css"
import ErrorNotFound from "../ErrorNotFound";
import { MobileTimePicker } from '@mui/x-date-pickers/MobileTimePicker';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import dayjs from "dayjs";
import 'dayjs/locale/de';

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
            <div id="big-modal">
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

export function AllRooms() {
    const user = useContext(AuthContext).user;
    const logout = useContext(AuthContext).logout;
    const navigate = useNavigate();
    const [isLoading, setIsLoading] = useState(true);
    const [rooms, setRooms] = useState(null);
    const [error, setError] = useState(null);

    useEffect(() => {
        async function getAllRooms() {
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
            <span className={'text-path'} onClick={() => navigate('/rooms')}>Свободные аудитории</span>

            {error && <div className={'modal error'} style={{right: '0', left: '1px'}}>Ошибка<br/>{error}</div>}
            {isLoading && <img src={'/3-dots-loader.svg'} alt={'Загрузка'} className="loader-dots"/>}
            <div id="main-cards">
                {rooms && rooms.map((room) => (
                    <div key={room.id} className="room-card" onClick={() => navigate('/rooms/'+room.id)}>
                        {room.image && <img src={room.image} alt={room.name} />}
                        {!room.image && <img src="/image_not_found.jpeg" alt={room.name}/>}
                        <h2>{room.name}</h2>
                        <p>Вместимость – {room.size}</p>
                    </div>
                ))}
            </div>   
            {rooms && rooms.length === 0 && <div className={'modal error'} style={{right: '0', left: '1px'}}>Ошибка<br/>Нет доступных комнат</div>}
        </div>
        </>
    );
}

export function CreateBooking() {
    const logout = useContext(AuthContext).logout;
    const { id } = useParams();
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    const [room, setRoom] = useState(null);
    const navigate = useNavigate();
    const user = useContext(AuthContext).user;
    const today = dayjs();
    const [dateBooking, setDateBooking] = useState(today);
    const [timeStart, setTimeStart] = useState(today.add(1, 'hour'));
    const [timeEnd, setTimeEnd] = useState(today.add(2, 'hour'));

    useEffect(() => {
        async function getRoomById() {
            try {
                const response = await fetch("http://localhost:8080/api/v1/rooms/"+id);
                if (!response.ok) {
                    if (response.status === 404) {
                        setError('Страница не найдена');
                        return;
                    } else {
                        setError('Произошла серверная ошибка');
                        return;
                    }
                } 
      
                const data = await response.json();
                setRoom(data);
            } catch (err) {
                setError('Произошла серверная ошибка');
            } finally {
                setIsLoading(false);
            }
        }
        getRoomById();
    }, [])

    if (error === 'Страница не найдена') {
        return <ErrorNotFound />;
    }

    if (!user) {
        return <LoginPage />;
    }

    return ( 
        <>
        <button className={"form-button"} onClick={logout} style={{right: '20px', top:'20px', position: 'absolute'}}>Выйти</button>
        <div id="main-container">
            <span className={'text-path'} onClick={() => navigate('/')}>Главная</span> 
            <span className={'text-path'}>/</span> 
            <span className={'text-path'} onClick={() => navigate('/rooms')}>Свободные аудитории</span>
            <span className={'text-path'}>/</span> 
            {room && <span className={'text-path'} onClick={() => navigate('/rooms/'+room.id)}>{room.name}</span>}

            {error && <div className={'modal error'} style={{right: '0', left: '1px'}}>Ошибка<br/>{error}</div>}
            {isLoading && <img src={'/3-dots-loader.svg'} alt={'Загрузка'} className="loader-dots"/>}
            <div id={'big-modal'}>
                <div className={'image-header-container'}>
                    {room && room.image && <img src={room.image} alt={room.name} className={'image-header'} />}
                    {room && !room.image && <img src={'/black-and-white-stripes.jpg'} className={'image-header'} />} 
                </div>
                {room && <h1 style={{fontSize: '40px', marginTop: '40px'}}>Аудитория {room.name}</h1>}
                {room && <p>{room.description}</p>}
                <h1 style={{marginTop: '40px'}}>Форма бронирования</h1>
                <br/>
                <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale="de">
                    <DatePicker 
                        label="Дата бронирования" 
                        value={dateBooking}
                        onChange={(newValue) => setDateBooking(newValue)}
                        disablePast 
                    /><br/><br/>
                    <MobileTimePicker 
                        label="Начало бронирования"
                        value={timeStart}
                        onChange={(newValue) => setTimeStart(newValue)}
                        disablePast
                    /><br/><br/>
                    <MobileTimePicker 
                        label="Конец бронирования"
                        value={timeEnd}
                        onChange={(newValue) => setTimeEnd(newValue)}
                        disablePast
                    /><br/><br/>
                </LocalizationProvider>
                <button className="form-button" style={{marginBottom: '50px'}}>Отправить</button>
            </div>  
        </div>
        </>
    );
}

export function MyBookings() {

}