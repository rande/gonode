import {
    REQUEST_NODE_REVISIONS,
    RECEIVE_NODE_REVISIONS,
    INVALIDATE_NODE_REVISIONS
} from '../constants/ActionTypes';

const assign = Object.assign || require('object.assign');


function nodeRevisions(state = {
    isFetching:    false,
    items:         [],
    didInvalidate: false
}, action) {
    switch (action.type) {
        case REQUEST_NODE_REVISIONS:
            return assign({}, state, {
                isFetching: true
            });

        case RECEIVE_NODE_REVISIONS:
            return assign({}, state, {
                isFetching:    false,
                didInvalidate: false,
                items:         action.items
            });

        case INVALIDATE_NODE_REVISIONS:
            return assign({}, state, {
                didInvalidate: true
            });

        default:
            return state;
    }
}


export default function nodesRevisionsByUuid(state = {}, action) {
    switch (action.type) {
        case REQUEST_NODE_REVISIONS:
        case RECEIVE_NODE_REVISIONS:
        case INVALIDATE_NODE_REVISIONS:
            return assign({}, state, {
                [action.uuid]: nodeRevisions(state[action.uuid], action)
            });

        default:
            return state;
    }
}
