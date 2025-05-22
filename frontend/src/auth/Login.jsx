import './Auth.css'
import {Link} from 'react-router';
import {useContext} from "react";
import AuthContext from "./AuthContext.jsx";

function Login() {

    const {user, login, logout} = useContext(AuthContext);

    if (user) {
        return (
            <>
                Привет user {user.header}
                <button onClick={logout}>Выйти</button>
            </>

        );
    }

    const handleSubmit = (event) => {
        const email = event.target.email.value;
        login(email)
    };

    return (
        <div id={"root-container"}>
            <h1>Авторизация</h1>
            <span className={"text-gray"}>Нет аккаунта?</span>
            <Link to={"../user/signup"} className={"text-link"}>Зарегистрироваться</Link>
            <br/><br/>

            <form onSubmit={handleSubmit}>
                <div className={"form-container"}>
                    <label>
                        E-mail<br/>
                        <input className={"form-field"}
                               type={"email"}
                               name={"email"}
                               placeholder={"mail@example.ru"}
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
                            required={true}
                        />
                    </label>
                </div>
                <br/>
                <button className={"form-button"} type={"submit"}>Войти</button>
            </form>
        </div>
    )
}

export default Login