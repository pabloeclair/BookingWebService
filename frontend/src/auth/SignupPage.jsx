import './Auth.css'
import {Link} from 'react-router';
import {useState} from "react";

function SignupPage() {
    const [error, setError] = useState(null);
    const [user, setUser] = useState(false);

    return (
        <>
        <div className={'form-container'}>
            <h1>Регистрация</h1>
            <span className={"text-gray"}>Уже есть аккаунт?</span>
            <Link to={"../login"} className={"text-link"}>Войти</Link>
            <br/><br/>
            <SignupForm setUser={setUser} setError={setError}/>
        </div>
        <Modal user={user} error={error} />
        </>
    );
}

function SignupForm({ setUser, setError }) {

    const [form, setForm] = useState({
        first_name: '',
        second_name: '',
        patronymic: '',
        email: '',
        password: ''
    });

    const handleSubmit = async (event) => {
        event.preventDefault();
        console.log('Регистрация:', form);

        const http = `http://localhost:8080/users`

        try {
            const response = await fetch(http, {
                method: 'POST',
                body: JSON.stringify(form),
                headers: {
                    'Content-Type': 'application/json; charset=UTF-8',
                },
            });
            const data = await response.json();
            if (!response.ok) {
                    throw new Error(data.errorCode);
            }
            setError(null);
            setUser(true)
        } catch (err) {
            setError(err.message);
            setUser(false);
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
            <label>Имя<br/>
                <input className={"form-field"}
                    type={"text"}
                    name={"first_name"}
                    placeholder={"Иван"}
                    value={form.first_name}
                    onChange={handleChange}
                    required={true}
                />
            </label>
            <br/>
            <label>Фамилия<br/>
                <input className={"form-field"}
                    type={"text"}
                    name={"second_name"}
                    placeholder={"Иванов"}
                    value={form.second_name}
                    onChange={handleChange}
                    required={true}
                />
            </label>
            <br/>
            <label>Отчество<br/>
                <input className={"form-field"}
                    type={"text"}
                    name={"patronymic"}
                    placeholder={"Иванович"}
                    value={form.patronymic}
                    onChange={handleChange}
                    required={false}
                />
            </label>
            <br/>
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

function Modal({ user, error }) {
    if (user) {
        return <div className={'modal ok'}>Регистрация прошла успешна<br/>Перейдите на страницу авторизации</div>;
    }
    if (!error) return null;
    switch (error) {
        case '400 BAD_REQUEST':
            return <div className={'modal error'}>Аккаунт с указанной почтой уже существует</div>;
        default:
            return <div className={'modal error'}>Произошла серверная ошибка<br/>Пожалуйста, обновите страницу или обратитесь на ресепшен на 4 этаже</div>;
    }
}


export default SignupPage;