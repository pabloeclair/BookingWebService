import './Auth.css'
import {Link} from 'react-router';
import {useState} from "react";

function SignupPage() {

    return (
        <div id={"root-container"}>
            <h1>Регистрация</h1>
            <span className={"text-gray"}>Уже есть аккаунт?</span>
            <Link to={"../login"} className={"text-link"}>Войти</Link>
            <br/><br/>
            <SignupForm />
        </div>
    );
}

function SignupForm() {

    const [error, setError] = useState(null);
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
    switch (error.message) {
        case '400 BAD_REQUEST':
            return <>Аккаунт с указанной почтой уже существует</>;
        default:
            return <>Произошла серверная ошибка<br/>Пожалуйста, обновите страницу или обратитесь на ресепшен на 4 этаже</>;
    }
}

function HtmlForm({ form, handleChange, handleSubmit }) {
    return (
        <form onSubmit={handleSubmit}>
            <div className={"form-container"}>
                <label>
                    Имя<br/>
                    <input className={"form-field"}
                            type={"text"}
                            name={"first_name"}
                            placeholder={"Иван"}
                            value={form.first_name}
                            onChange={handleChange}
                            required={true}
                    />
                </label>
            </div>
            <div className={"form-container"}>
                <label>
                    Фамилия<br/>
                    <input className={"form-field"}
                            type={"text"}
                            name={"second_name"}
                            placeholder={"Иванов"}
                            value={form.second_name}
                            onChange={handleChange}
                            required={true}
                    />
                </label>
            </div>
            <div className={"form-container"}>
                <label>
                    Отчество<br/>
                    <input className={"form-field"}
                            type={"text"}
                            name={"patronymic"}
                            placeholder={"Иванович"}
                            value={form.patronymic}
                            onChange={handleChange}
                            required={false}
                    />
                </label>
            </div>
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
                    <input className={"form-field"}
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

export default SignupPage;