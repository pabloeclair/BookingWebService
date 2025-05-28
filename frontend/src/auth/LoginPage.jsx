import './Auth.css'
import {Link, useNavigate} from 'react-router';
import {useContext, useState} from "react";
import AuthContext from "./AuthContext.jsx";

function LoginPage() {

    const user = useContext(AuthContext).user;
    const logout = useContext(AuthContext).logout;
    const navigate = useNavigate();
    const [error, setError] = useState(null);

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
        <>
        <div className={'form-container'}>
            <h1>Авторизация</h1>
            <span className={"text-gray"}>Нет аккаунта?</span>
            <Link to={"../signup"} className={"text-link"}>Зарегистрироваться</Link>
            <br/><br/>
            <LoginForm setError={setError}/>
        </div>
        <Error error={error}/>
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
        console.log('Регистрация:', form);

        const email = form.email;
        const password = form.password;
        const http = `http://localhost:8080/users?email=${email}&password=${password}`

        try {
            const response = await fetch(http);
            const data = await response.json();
            if (!response.ok) {
                    throw new Error(data.errorCode);
            }
            login(data);
            setError(null);
        } catch (err) {
            setError(err);
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

function Error({ error }) {
    if (error == null) {
        return;
    } else if (error.message === "404 NOT_FOUND") {
        return <div className={'modal error'}>Аккаунта с указанной почтой не существует<br/>Пожалуйста, зарегистрируйтесь</div>;
    } 
    return <div className={'modal error'}>Произошла серверная ошибка<br/>Пожалуйста, обновите страницу или<br/>обратитесь на ресепшен на 4 этаже</div>;
}

export default LoginPage;