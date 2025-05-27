import './Auth.css'
import {Link, Navigate, useNavigate} from 'react-router';
import {useContext, useState} from "react";
import AuthContext from "./AuthContext.jsx";

function LoginPage() {

    const user = useContext(AuthContext).user;
    const logout = useContext(AuthContext).logout;
    const navigate = useNavigate();

    if (user) {
        return (
            <>
                Привет user {user.header}<br/>
                <button className={"form-button"} onClick={logout}>Выйти</button>
                <br/>
                <button className={"form-button"} onClick={() => {navigate("/signup")}}>Тест</button>
            </>
        );
    }

    return (
        <div id={"root-container"}>
            <h1>Авторизация</h1>
            <span className={"text-gray"}>Нет аккаунта?</span>
            <Link to={"../signup"} className={"text-link"}>Зарегистрироваться</Link>
            <br/><br/>
            <LoginForm />
        </div>
    );
}

function htmlForm(form, handleChange, handleSubmit) {
    return (
        <form onSubmit={handleSubmit}>
            <div className={"form-container"}>
                <label>
                    E-mail<br/>
                    <input className={"form-field"}
                        type={"email"}
                        name={"email"}
                        placeholder={"mail@example.ru"}
                        value={form.email}
                        onChange={handleChange}
                        required={true}
                    />
                </label>
            </div>
            <div className={"form-container"}>
                <label>
                    Пароль<br/>
                    <input
                        className={"form-field"}
                        type={"password"}
                        name={"password"}
                        value={form.password}
                        onChange={handleChange}
                        required={true}
                    />
                </label>
            </div>
            <br/>
            <button className={"form-button"} type={"submit"}>Войти</button>
        </form>
    );
}

function LoginForm() {

    const [form, setForm] = useState({
        email: '',
        password: ''
    });
    const login = useContext(AuthContext).login;

    const handleSubmit = (event) => {
        event.preventDefault();
        console.log('Регистрация:', form);

        const email = form.email;
        const password = form.password;

        fetch("http://localhost:8080/users?email="+email+"&password="+password)
            .then((response) => {
                if (!response.ok) {
                    throw Error(response.json().errorCode)
                }
                const res = response.json();
                return res;
            }).then((data) => {
                login(data)
            })
            .catch((error) => error);
    };

    const handleChange = (event) => {
        const {name, value} = event.target;
        setForm(prevForm => ({
            ...prevForm,
            [name]: value
        }));
    };

    return htmlForm(form, handleChange, handleSubmit);
}


export default LoginPage;