import {
    RECEIVE_NODES,
    REQUEST_NODE,
    RECEIVE_NODE,
    RECEIVE_NODE_CREATION,
    RECEIVE_NODE_UPDATE
} from '../constants/ActionTypes';

const assign = Object.assign || require('object.assign');


function node(state = {
    isFetching:    false,
    node:          null,
    didInvalidate: false
}, action) {
    switch (action.type) {
        case REQUEST_NODE:
            return assign({}, state, {
                isFetching:    true,
                didInvalidate: false
            });

        case RECEIVE_NODE:
        case RECEIVE_NODE_CREATION:
        case RECEIVE_NODE_UPDATE:
            return assign({}, state, {
                isFetching:    false,
                didInvalidate: false,
                node:          action.node
            });

        default:
            return state;
    }
}

export default function nodesByUuid(state = {}, action) {
    switch (action.type) {
        case RECEIVE_NODES:
            return assign({}, state, action.items.reduce((newNodes, item) => {
                newNodes[item.uuid] = node(null, {
                    type: RECEIVE_NODE,
                    node: item
                });

                return newNodes;
            }, {}));

        case REQUEST_NODE:
        case RECEIVE_NODE:
        case RECEIVE_NODE_CREATION:
        case RECEIVE_NODE_UPDATE:
            return assign({}, state, {
                [action.uuid]: node(state[action.uuid], action)
            });

        default:
            return state;
    }
}
