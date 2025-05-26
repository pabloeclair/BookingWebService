import './Auth.css'
import {Link} from "react-router";
import {useContext} from "react";
import AuthContext from "./AuthContext.jsx";

function Signup() {

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
            <h1>Регистрация</h1>
            <span className={"text-gray"}>Уже есть аккаунт?</span>
            <Link to={"../user/login"} className={"text-link"}>Войти</Link>
            <br/><br/>

            <form onSubmit={handleSubmit}>
                <div className={"form-container"}>
                    <label>
                        Имя<br/>
                        <input className={"form-field"}
                               type={"text"}
                               name={"first_name"}
                               placeholder={"Иван"}
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
    );
}

export default Signup