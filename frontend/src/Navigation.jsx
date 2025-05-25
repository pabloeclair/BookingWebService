import {Navigate} from "react-router";
import {useContext} from "react";
import AuthContext from "./auth/AuthContext.jsx";

function Navigation() {
    const {user, login, logout} = useContext(AuthContext);

    if (!user) {
        return <Navigate to={"/user/login"} replace />;
    }
    return "HelloWorld";
}

export default Navigation;