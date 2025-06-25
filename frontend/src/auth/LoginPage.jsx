import './Auth.css'
import {Link, useNavigate} from 'react-router';
import {useContext, useState} from "react";
import AuthContext from "./AuthContext.jsx";
import MainAdminPage from "../user/MainAdminPage.jsx";

function LoginPage() {

    const user = useContext(AuthContext).user;
    const logout = useContext(AuthContext).logout;
    const navigate = useNavigate();
    const [error, setError] = useState(null);

    if (user && (user.role === 'ADMIN' || user.role === 'MAIN_ADMIN')) {
        return <MainAdminPage/>;
    }

    if (user) {
        console.log(user)
        return (
            <>
                Привет user {user.first_name}<br/>
                <button className={"form-button"} onClick={logout}>Выйти</button>
                <br/>
                <button className={"form-button"} onClick={() => { navigate("/signup") }}>Тест</button>
            </>
        );
    }

    if (!error) {
        return (
        <>
        <div className={'form-container'}>
            <h1>Авторизация</h1>
            <span className={"text-gray"}>Нет аккаунта?</span>
            <Link to={"../signup"} className={"text-link"}>Зарегистрироваться</Link>
            <br/><br/>
            <LoginForm setError={setError}/>
        </div>
        </>
    );
    }

    return (
        <>
        <div className={'form-container'}>
            <h1>Авторизация</h1>
            <span className={"text-gray"}>Нет аккаунта?</span>
            <Link to={"../signup"} className={"text-link"}>Зарегистрироваться</Link>
            <br/><br/>
            <LoginForm setError={setError}/>
        </div>
        <div className={'modal error'}>Ошибка<br/>{error}</div>
        </>
    );
    
}

function LoginForm({ setError }) {

    const [form, setForm] = useState({
        email: '',
        password: ''
    });
    const login = useContext(AuthContext).login;

    const handleSubmit = async (event) => {
        event.preventDefault();

        try {
            const responseLogin = await fetch('http://localhost:7070/api/v1/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json; charset=utf-8' },
                body: JSON.stringify(form)
            });

            if (!responseLogin.ok) {
                let errorMessage;
                switch(responseLogin.status) {
                    case 404:
                        errorMessage = "Почта не существует";
                        break;
                    case 401:
                        errorMessage = "Пароль неверный";
                        break;
                    default:
                        errorMessage = 'Произошла серверная ошибка';
                }
                throw new Error(errorMessage);
            }

            setForm(prevForm => ({
                ...prevForm,
                email: '',
                password: '',
            }))

            const token = responseLogin.headers.get('authorization');
            const responseParseJWT = await fetch('http://localhost:7070/api/v1/user', {
                method: 'GET',
                headers: { 'Authorization': token }
            });

            if (!responseParseJWT.ok) {
                throw new Error('Произошла серверная ошибка');
            }

            let data = await responseParseJWT.json();
            data.key = token;

            login(data);
            setError(null);
        } catch (err) {
            setError(err.message);
        }
    };

    const handleChange = (event) => {
        const {name, value} = event.target;
        setForm(prevForm => ({
            ...prevForm,
            [name]: value
        }));
    };

    return (
        <form onSubmit={handleSubmit}>
            <label>E-mail<br/>
                <input className={"form-field"}
                    type={"email"}
                    name={"email"}
                    placeholder={"mail@example.ru"}
                    value={form.email}
                    onChange={handleChange}
                    required={true}
                />
            </label>
            <br/>
            <label>Пароль<br/>
                <input className={"form-field"}
                    type={"password"}
                    name={"password"}
                    value={form.password}
                    onChange={handleChange}
                    required={true}
                />
            </label>
            <br/>
            <button className={"form-button"} type={"submit"}>Войти</button>
        </form>
    );
}

export default LoginPage;