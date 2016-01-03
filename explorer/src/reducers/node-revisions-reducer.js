import _          from 'lodash';
import * as types from '../constants/ActionTypes';

const assign = Object.assign || require('object.assign');


function nodeRevisions(state = {
    isFetching:    false,
    items:         [],
    didInvalidate: false
}, action) {
    switch (action.type) {
        case types.REQUEST_NODE_REVISIONS:
            return assign({}, state, {
                isFetching: true
            });

        case types.RECEIVE_NODE_REVISIONS:
            return assign({}, state, {
                isFetching:    false,
                didInvalidate: false,
                items:         action.items
            });

        case types.RECEIVE_NODE_UPDATE:
            return assign({}, state, {
                didInvalidate: true
            });

        default:
            return state;
    }
}

export default function nodesRevisionsByUuid(state = {}, action) {
    switch (action.type) {
        case types.REQUEST_NODE_REVISIONS:
        case types.RECEIVE_NODE_REVISIONS:
        case types.RECEIVE_NODE_UPDATE:
            return assign({}, state, {
                [action.uuid]: nodeRevisions(state[action.uuid], action)
            });

        default:
            return state;
    }
}
