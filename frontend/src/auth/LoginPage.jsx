import './Auth.css'
import {Link, useNavigate} from 'react-router';
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

function HtmlForm({ form, handleChange, handleSubmit }) {
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

    const [error, setError] = useState(null);
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

    if (error) {
        return (
            <>
                <HtmlForm form={form} handleChange={handleChange} handleSubmit={handleSubmit} />
                <br/>
                <div className={"text-error"}>
                    <HandleError error={error}/>
                </div>
            </>
        );
    }
    return <HtmlForm form={form} handleChange={handleChange} handleSubmit={handleSubmit} />;
}

function HandleError({ error }) {
    if (error.message === "404 NOT_FOUND") {
        return <>Аккаунта с указанной почтой не существует<br/>Пожалуйста, зарегистрируйтесь</>;
    } 
    return <>Произошла серверная ошибка<br/>Пожалуйста, обновите страницу или обратитесь на ресепшен на 4 этаже</>;
}

export default LoginPage;