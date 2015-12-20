import _          from 'lodash';
import cookies    from 'browser-cookies';
import * as types from '../constants/ActionTypes';


export default function nodes(state = {
    token:      cookies.get('token'), // try to load token from cookie
    isFetching: false,
    failed:     false
}, action) {
    switch (action.type) {
        case types.LOGIN_ATTEMPT:
            return Object.assign({}, state, {
                isFetching: true,
                token:      null
            });

        case types.LOGIN_SUCCEEDED:
            return Object.assign({}, state, {
                isFetching: false,
                token:      action.token,
                failed:     false
            });

        case types.LOGIN_FAILED:
            return Object.assign({}, state, {
                isFetching: false,
                failed:     true
            });

        case types.LOGOUT:
            return Object.assign({}, state, {
                token:      null,
                isFetching: false,
                failed:     false
            });

        default:
            return state;
    }
}
