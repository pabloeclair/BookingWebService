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
    const user = useContext(AuthContext).user;
    const navigate = useNavigate();

    const { id: roomId } = useParams();

    const [isRoomAndBookingsLoading, setIsRoomAndBookingsLoading] = useState(true);
    const [isFormSending, setIsFormSending] = useState(false);

    const [error, setError] = useState(null);
    const [room, setRoom] = useState(null);
    const [roomBookings, setRoomBookings] = useState([]);

    const today = dayjs();
    const [dateBooking, setDateBooking] = useState(today);
    const [timeStart, setTimeStart] = useState(today.add(1, 'hour'));
    const [timeEnd, setTimeEnd] = useState(today.add(2, 'hour'));
    const [errorTimeStart, setErrorTimeStart] = useState(null);
    const [errorTimeEnd, setErrorTimeEnd] = useState(null);
    const [errorDateBooking, setErrorDateBooking] = useState(null);

    useEffect(() => {
        async function getRoomAndBookingsById() {
            try {
                const responseRoom = await fetch("http://localhost:8080/api/v1/rooms/"+roomId);
                if (!responseRoom.ok) {
                    if (responseRoom.status === 404) {
                        setError('Страница не найдена');
                        return;
                    } else {
                        setError('Произошла серверная ошибка');
                        return;
                    }
                } 
      
                let data = await responseRoom.json();
                setRoom(data);

                const responseBookings = await fetch("http://localhost:8080/api/v1/bookings/room/"+roomId)
                if (!responseBookings.ok) {
                    setError('Произошла серверная ошибка');
                    return;
                }
                data = await responseBookings.json();
                setRoomBookings(data);
            } catch (err) {
                setError('Произошла серверная ошибка');
            } finally {
                setIsRoomAndBookingsLoading(false);
            }
        }
        getRoomAndBookingsById();
    }, [])
    
    const handleSubmit = async (event) => {
        event.preventDefault();
        setIsFormSending(true);
        try {
            const response = await fetch('http://localhost:8080/api/v1/bookings', {
                method: 'POST',
                body: {
                    'room_id': roomId,
                    'user_id': user.id,
                    'booking_date': dateBooking.format('DD.MM.YYYY'),
                    'booking_start': timeStart.format('HH:mm:ss'),
                    'booking_end': timeEnd.format('HH:mm:ss')
                },
                headers: {
                    'Authorization': user.key,
                    'Content-Type': 'application/json; charset=utf-8'
                }
            });

            if (!response.ok) {
                const data = await response.json();
                switch (response.status) {
                    case 400:
                        setError(data.error_message);
                        return;
                    case 403:
                        setError(data.error_message);
                        return;
                    case 401:
                        setError(
                            <>Ваша сессия истекла. Пожалуйста, 
                            <Link to={"../login"} className={"text-link"}>авторизуйтесь</Link>
                            заново.</>
                        );
                        return;
                    default:
                        setError('Произошла серверная ошибка');
                        return;
                }
            }
        } catch (err) {
            setError('Произошла серверная ошибка');
        } finally {
            setIsFormSending(false);
        }
        
    }

    if (error === 'Страница не найдена') {
        return <ErrorNotFound />;
    }

    if (!user) {
        return <LoginPage />;
    }

    if (isRoomAndBookingsLoading) {
        return (
            <>
            <button className={"form-button"} onClick={logout} style={{right: '20px', top:'20px', position: 'absolute'}}>Выйти</button>
            <div id="main-container">
                <span className={'text-path'} onClick={() => navigate('/')}>Главная</span> 
                <span className={'text-path'}>/</span> 
                <span className={'text-path'} onClick={() => navigate('/rooms')}>Свободные аудитории</span>
                <span className={'text-path'}>/</span>
                <span className={'text-path'}>...</span>  
                <img src={'/3-dots-loader.svg'} alt={'Загрузка'} className="loader-dots"/>
            </div>
            </>
        );
    }

    const handleDateBooking = (newValue) => {
        setDateBooking(newValue);
        const resultValidation = validateDateBooking(newValue, today);
        if (resultValidation.isValid) {
            setErrorDateBooking(null);
        } else {
            setErrorDateBooking(resultValidation.error);
        }
    }

    const handleTimeStart = (newValue) => {
        setTimeStart(newValue);
        const resultValidation = validateTimeStart(newValue, timeEnd, dateBooking, today, roomBookings);
        if (resultValidation.isValid) {
            setErrorTimeStart(null);
        } else {
            setErrorTimeStart(resultValidation.error);
        }
    }

    const handleTimeEnd = (newValue) => {
        setTimeEnd(newValue);
        const resultValidation = validateTimeEnd(timeStart, newValue, dateBooking, roomBookings, today);
        if (resultValidation.isValid) {
            setErrorTimeEnd(null);
        } else {
            setErrorTimeEnd(resultValidation.error);
        }
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
                        onChange={handleDateBooking}
                        disablePast 
                        slotProps={{
                                textField: {
                                    error: Boolean(errorDateBooking),
                                    helperText: errorDateBooking,
                                },
                        }}
                        sx={{ width: '60%' }} 
                    /><br/><br/><br/>
                    <MobileTimePicker 
                        label="Начало бронирования"
                        value={timeStart}
                        onChange={handleTimeStart}
                        minTime={dateBooking && dateBooking.isSame(today, 'day') ? today : dayjs().startOf('day')}
                        maxTime={timeEnd ? timeEnd.subtract(1, 'minute') : dayjs().endOf('day')}
                        slotProps={{
                                textField: {
                                    error: Boolean(errorTimeStart),
                                    helperText: errorTimeStart,
                                },
                        }}
                        sx={{ width: '60%' }} 
                    /><br/><br/><br/>
                    <MobileTimePicker 
                        label="Конец бронирования"
                        value={timeEnd}
                        onChange={handleTimeEnd}
                        minTime={timeStart ? timeStart.add(1, 'minute') : (dateBooking && dateBooking.isSame(today, 'day') ? today : dayjs().startOf('day'))}
                        slotProps={{
                                textField: {
                                    error: Boolean(errorTimeEnd),
                                    helperText: errorTimeEnd,
                                },
                        }}
                        sx={{ width: '60%' }} 
                    /><br/><br/>
                </LocalizationProvider>
                <button className="form-button" style={{marginBottom: '50px'}} onSubmit={handleSubmit}>Отправить</button>
                {isFormSending && <img src={'/ring-loader.svg'} alt={'Загрузка'} className="loader-ring"/>}
            </div>  
        </div>
        </>
    );
}

function isOverlap(start1, end1, start2, end2) {
    return start1.isBefore(end2) && start2.isBefore(end1);
}

function validateDateBooking(dateBooking, today) {
    if (!dateBooking) {
        return {isValid: false, error: 'Данное поле обязательно'};
    }

    if (dateBooking.isBefore(today)) {
        return {isValid: false, error: 'Дата должна быть не раньше текущего дня'};
    }

    return {isValid: true, error: ''};
}

function validateTimeStart(timeStart, timeEnd, dateBooking, today, roomBookings) {
    if (!timeStart) {
        return {isValid: false, error: 'Данное поле обязательно'};
    }

    if (!timeEnd || !dateBooking) {
        return {isValid: false, error: 'Пожалуйста заполните все остальные поля'};
    }

    const startDateTime = dateBooking
        .hour(timeStart.hour())
        .minute(timeStart.minute())
        .second(0)
        .millisecond(0);

    if (dateBooking.isSame(today, 'day') && startDateTime.isBefore(today)) {
        return {isValid: false, error: 'Время начала должно быть не раньше текущего времени'}
    }

    // время начала < время конца
    if (!timeStart.isBefore(timeEnd)) {
        return {isValid: false, error: 'Время начала должно быть раньше времени окончания'};
    }

    // пересечения с существующими бронированиями
    for (const booking of roomBookings) {
        if (dayjs(booking.booking_date).isSame(dateBooking, 'day')) {
            const existingStart = dayjs(booking.booking_start);
            const existingEnd = dayjs(booking.booking_end);

            if (isOverlap(timeStart, timeEnd, existingStart, existingEnd)) {
                return {isValid: false, error: 'Время начала пересекается с существующим бронированием: ' + existingStart.format('HH:mm') + ' - ' + existingEnd.format('HH:mm')};
            }
        }
    }

    return {isValid: true, error: ''};
}

function validateTimeEnd(timeStart, timeEnd, dateBooking, roomBookings, today) {
    if (!timeEnd) {
        return {isValid: false, error: 'Данное поле обязательно'};
    }

    if (!timeStart || !dateBooking) {
        return {isValid: false, error: 'Пожалуйста заполните все остальные поля'};
    }

    const endDateTime = dateBooking
        .hour(timeEnd.hour())
        .minute(timeEnd.minute())
        .second(0)
        .millisecond(0);

    if (dateBooking.isSame(today, 'day') && endDateTime.isBefore(dateBooking)) {
        return {isValid: false, error: 'Время начала должно быть не раньше текущего времени'}
    }

    // время конца > время начала
    if (!timeEnd.isAfter(timeStart)) {
        return {isValid: false, error: 'Время окончания должно быть позже времени начала'};
    }

    // пересечения с существующими бронированиями
    for (const booking of roomBookings) {
        if (dayjs(booking.booking_date).isSame(dateBooking, 'day')) {
            const existingStart = dayjs(booking.booking_start);
            const existingEnd = dayjs(booking.booking_end);

            if (isOverlap(timeStart, timeEnd, existingStart, existingEnd)) {
                return {isValid: false, error: 'Время окончания пересекается с существующим бронированием: ' + existingStart.format('HH:mm') + ' - ' + existingEnd.format('HH:mm')};
            }
        }
    }

    return {isValid: true, error: ''};
}

export function MyBookings() {

}