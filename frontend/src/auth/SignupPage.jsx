import './Auth.css'
import {Link} from 'react-router';
import {useState} from "react";

function SignupPage() {
    const [error, setError] = useState(null);
    const [user, setUser] = useState(false);

    if (error) {
        return (
            <>
            <div className={'form-container'}>
                <h1>Регистрация</h1>
                <span className={"text-gray"}>Уже есть аккаунт?</span>
                <Link to={"../login"} className={"text-link"}>Войти</Link>
                <br/><br/>
                <SignupForm setUser={setUser} setError={setError}/>
            </div>
            <div className={'modal error'}>Ошибка<br/>{error}</div>
            </>
        );
    }

    if (user) {
        return (
            <>
            <div className={'form-container'}>
                <h1>Регистрация</h1>
                <span className={"text-gray"}>Уже есть аккаунт?</span>
                <Link to={"../login"} className={"text-link"}>Войти</Link>
                <br/><br/>
                <SignupForm setUser={setUser} setError={setError}/>
            </div>
            <div className={'modal ok'}>Регистрация прошла успешна<br/>Перейдите на страницу авторизации</div>
            </>
        );
    }

    return (
        <>
        <div className={'form-container'}>
            <h1>Регистрация</h1>
            <span className={"text-gray"}>Уже есть аккаунт?</span>
            <Link to={"../login"} className={"text-link"}>Войти</Link>
            <br/><br/>
            <SignupForm setUser={setUser} setError={setError}/>
        </div>
        </>
    );
}

function SignupForm({ setUser, setError }) {

    const [form, setForm] = useState({
        first_name: '',
        second_name: '',
        patronymic: '',
        email: '',
        password: '',
    });

    const handleSubmit = async (event) => {
        event.preventDefault();

        try {
            const response = await fetch('http://localhost:7070/api/v1/signup', {
                method: 'POST',
                body: JSON.stringify(form),
                headers: {
                    'Content-Type': 'application/json; charset=utf-8',
                },
            });

            if (!response.ok) {
                let errorMessage;
                if (response.status === 409) {
                    errorMessage = 'Аккаунт с указанной почтой уже существует';
                } else {
                    errorMessage = 'Произошла серверная ошибка';
                }
                throw new Error(errorMessage);
            }

            setForm(prevForm => ({
                ...prevForm,
                first_name: '',
                second_name: '',
                patronymic: '',
                email: '',
                password: '',
            }))

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

export default SignupPage;