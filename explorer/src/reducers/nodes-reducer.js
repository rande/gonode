import { UPDATE_PATH } from 'redux-simple-router';
import {
    REQUEST_NODES,
    RECEIVE_NODES,
    SELECT_NODE,
    RECEIVE_NODE_UPDATE,
    RECEIVE_NODE_CREATION
} from '../constants/ActionTypes';

const assign = Object.assign || require('object.assign');


export default function nodes(state = {
    isFetching:    false,
    items:         [],
    itemsPerPage:  10,
    currentPage:   1,
    previousPage:  null,
    nextPage:      null,
    currentUuid:   null,
    didInvalidate: true
}, action) {
    switch (action.type) {
        case REQUEST_NODES:
            return assign({}, state, {
                isFetching:   true,
                currentPage:  action.page,
                itemsPerPage: action.perPage
            });

        case RECEIVE_NODES:
            return assign({}, state, {
                isFetching:    false,
                items:         action.items,
                itemsPerPage:  action.itemsPerPage,
                currentPage:   action.page,
                previousPage:  (action.previous !== 0 ? action.previous : null),
                nextPage:      (action.next     !== 0 ? action.next     : null),
                didInvalidate: false
            });

        case SELECT_NODE:
            return assign({}, state, {
                currentUuid: action.nodeUuid
            });

        case RECEIVE_NODE_UPDATE:
            const { items } = state;
            const updatedItems = items.map(item => {
                if (item.uuid === action.node.uuid) {
                    return action.node;
                }

                return item;
            });

            return assign({}, state, {
                items: updatedItems
            });

        case RECEIVE_NODE_CREATION:
            return assign({}, state, {
                didInvalidate: true
            });

        default:
            return state;
    }
}
