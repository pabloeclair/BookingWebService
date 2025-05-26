import {Navigate} from "react-router";
import {useContext} from "react";
import AuthContext from "./auth/AuthContext.jsx";

function Navigation() {
    const user = useContext(AuthContext).user;

    if (!user) {
        return <Navigate to={"/login"} replace />;
    }
    return "HelloWorld";
}

export default Navigation;