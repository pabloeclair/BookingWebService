import React, {useState} from 'react';
import AuthContext from './AuthContext';

function AuthProvider({children}) {
    const [user, setUser] = useState(null);

    const login = (response) => setUser({
        email: response.email,
        role: response.role
    });
    const logout = () => setUser(null);

    return (
        <AuthContext.Provider value={{user, login, logout}}>
            {children}
        </AuthContext.Provider>
    );
}

export default AuthProvider
