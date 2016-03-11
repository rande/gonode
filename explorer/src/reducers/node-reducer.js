import _          from 'lodash';
import * as types from '../constants/ActionTypes';

const assign = Object.assign || require('object.assign');


function node(state = {
    isFetching:    false,
    node:          null,
    didInvalidate: false
}, action) {
    switch (action.type) {
        case types.REQUEST_NODE:
            return assign({}, state, {
                isFetching:    true,
                didInvalidate: false
            });

        case types.RECEIVE_NODE:
            return assign({}, state, {
                isFetching:    false,
                didInvalidate: false,
                node:          action.node
            });

        case types.RECEIVE_NODE_UPDATE:
            return assign({}, state, {
                didInvalidate: true
            });

        default:
            return state;
    }
}

export default function nodesByUuid(state = {}, action) {
    switch (action.type) {
        case types.REQUEST_NODE:
        case types.RECEIVE_NODE:
        case types.RECEIVE_NODE_UPDATE:
            return assign({}, state, {
                [action.nodeUuid]: node(state[action.nodeUuid], action)
            });

        default:
            return state;
    }
}
