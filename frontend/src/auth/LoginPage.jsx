import './Auth.css'
import {Link, useNavigate} from 'react-router';
import {useContext, useState} from "react";
import AuthContext from "./AuthContext.jsx";
import generateSHA256Hash from '../GenerateHash';
import MainAdminPage from "../admin/MainAdminPage.jsx";

function LoginPage() {

    const user = useContext(AuthContext).user;
    const logout = useContext(AuthContext).logout;
    const navigate = useNavigate();
    const [error, setError] = useState(null);

    if (user && user.role !== 'USER') {
        return <MainAdminPage/>;
    }

    if (user) {
        return (
            <>
                Привет user {user.first_name}<br/>
                <button className={"form-button"} onClick={logout}>Выйти</button>
                <br/>
                <button className={"form-button"} onClick={() => {
                    navigate("/signup")
                }}>Тест
                </button>
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
        <div className={'modal error'}>{error}</div>
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
            const hashPassword = await generateSHA256Hash(form.password);
            const response = await fetch('http://localhost:8080/users?' + new URLSearchParams({
                email: form.email,
                key: hashPassword,
            }).toString());

            setForm(prevForm => ({
                ...prevForm,
                email: '',
                password: '',
            }))

            const data = await response.json();
            if (!response.ok) {
                let errorMessage;
                switch(response.status) {
                    case 404:
                        errorMessage = 'Аккаунта с указанной почтой не существует';
                        break;
                    case 401:
                        errorMessage = 'Неверный пароль';
                        break;
                    default:
                        errorMessage = 'Произошла серверная ошибка';
                }
                throw new Error(errorMessage);
            }
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