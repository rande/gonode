import _          from 'lodash';
import * as types from '../constants/ActionTypes';


function node(state = {
    isFetching: false,
    node:       null
}, action) {
    switch (action.type) {
        case types.REQUEST_NODE:
            return Object.assign({}, state, {
                isFetching: true
            });

        case types.RECEIVE_NODE:
            return Object.assign({}, state, {
                isFetching: false,
                node:       action.node
            });

        default:
            return state;
    }
}

export default function nodesByUuid(state = {}, action) {
    switch (action.type) {
        case types.REQUEST_NODE:
        case types.RECEIVE_NODE:
            return Object.assign({}, state, {
                [action.nodeUuid]: node(state[action.nodeUuid], action)
            });

        default:
            return state;
    }
}
