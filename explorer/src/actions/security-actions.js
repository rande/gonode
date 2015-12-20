import * as types  from '../constants/ActionTypes';
import Api         from '../Api';
import { history } from '../routing';
import cookies     from 'browser-cookies';


function attemptLogin(credentials) {
    return {
        type:  types.LOGIN_ATTEMPT,
        login: credentials.login
    };
}

function loginSucceeded(token) {
    return {
        type: types.LOGIN_SUCCEEDED,
        token
    };
}

function loginFailed() {
    return {
        type: types.LOGIN_FAILED
    };
}

export function login(credentials) {
    return (dispatch, getState) => {
        dispatch(attemptLogin(credentials));
        Api.login(credentials)
            .then(response => {
                const { status } = response;
                if (status === 'OK') {
                    const { token } = response;
                    cookies.set('token', token);
                    dispatch(loginSucceeded(token));
                    history.replaceState(null, '/nodes');
                } else {
                    dispatch(loginFailed());
                }
            })
            .catch(err => {
                dispatch(loginFailed());
            })
        ;
    };
}

export function logout() {
    cookies.erase('token');

    return {
        type: types.LOGOUT
    };
}
