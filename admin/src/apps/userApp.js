// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

import base64 from 'base64-js';

// @todo: use a middleware to avoid coupling userApp and nodeApp ...
import {searchNodes, getNodeServices} from './nodeApp';

export const DEFINE_USER = 'DEFINE_USER';
export const AUTHENTICATE_USER = 'AUTHENTICATE_USER';
export const AUTHENTICATING_USER = 'AUTHENTICATING_USER';
export const AUTHENTICATED_USER = 'AUTHENTICATED_USER';
export const AUTHENTICATION_FAILED_USER = 'AUTHENTICATION_FAIL_USER';

export const defineUser = (user, roles) => ({
    type: DEFINE_USER,
    user,
    roles
});

export const logoutUser = () => {
    return {
        type: DEFINE_USER,
        user: defaultState.user,
        roles: defaultState.roles,
    }
};

export const authenticateUser = (login, password) => {
    return (dispatch) => {

        let form = new URLSearchParams();
        form.append('username', login);
        form.append('password', password);

        fetch('/api/v1.0/login', {
            method: 'POST',
            body: form
        }).then((res) => {
            return res.json();
        }).then((data) => {
            if (data.status == "OK") {
                let token = data.token.split('.')[1];

                switch (token.length % 4) {
                    case 2: {
                        token += '==';
                        break;
                    }
                    case 3: {
                        token += '=';
                        break;
                    }
                }

                let roles = JSON.parse(String.fromCharCode(...base64.toByteArray(token))).rls;

                dispatch({
                    type: AUTHENTICATED_USER,
                    user: login,
                    roles: roles,
                    auth: data.token
                });

                dispatch(searchNodes());
                dispatch(getNodeServices());
            } else {
                dispatch({
                    type: AUTHENTICATION_FAILED_USER,
                    login: login,
                    password: password
                })
            }

        }).catch((err) => {
            console.log(err) ;
        });

        dispatch({
            type: AUTHENTICATING_USER,
        });
    }
};

const defaultState = {
    user: 'anonymous',
    roles: ['anonymous'],
    login: {
        login: 'admin',
        password: 'admin'
    },
    state: DEFINE_USER,
    auth: null
};


export function userApp(state = defaultState, action) {
    switch (action.type) {
    case DEFINE_USER:
        return {
            ...state,
            user: action.user,
            roles: action.roles,
            state: DEFINE_USER
        };

    case AUTHENTICATE_USER:
        return {
            ...state,
            state: AUTHENTICATE_USER
        };

    case AUTHENTICATING_USER:
        return {
            ...state,
            state: AUTHENTICATING_USER
        };

    case AUTHENTICATION_FAILED_USER:
        return {
            ...state,
            state: AUTHENTICATION_FAILED_USER,
            login: {
                login: action.login,
                password: action.password
            }
        };

    case AUTHENTICATED_USER:
        return {
            ...state,
            user: action.user,
            roles: action.roles,
            auth: action.auth,
            state: AUTHENTICATED_USER,
            login: {
                login: null,
                password: null
            }
        };

    default:
        return state;
    }
}
