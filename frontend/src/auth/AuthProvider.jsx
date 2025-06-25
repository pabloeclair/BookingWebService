import React, {useState} from 'react';
import AuthContext from './AuthContext';

function AuthProvider({children}) {
    const [user, setUser] = useState(null);

    const login = (response) => setUser({
        email: response.email,
        first_name: response.first_name,
        second_name: response.second_name,
        patronymic: response.patronymic,
        role: response.role,
        key: response.key
    });
    const logout = () => setUser(null);

    return (
        <AuthContext.Provider value={{user, login, logout}}>
            {children}
        </AuthContext.Provider>
    );
}

export default AuthProvider
