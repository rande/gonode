import { UPDATE_PATH } from 'redux-simple-router';
import {
    REQUEST_NODES,
    RECEIVE_NODES,
    SELECT_NODE
} from '../constants/ActionTypes';

const assign = Object.assign || require('object.assign');


export default function nodes(state = {
    isFetching:    false,
    uuids:         [],
    uuid:          null,
    itemsPerPage:  10,
    currentPage:   1,
    previousPage:  null,
    nextPage:      null,
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
                uuids:         action.items.map(item => item.uuid),
                isFetching:    false,
                itemsPerPage:  action.itemsPerPage,
                currentPage:   action.page,
                previousPage:  (action.previous !== 0 ? action.previous : null),
                nextPage:      (action.next     !== 0 ? action.next     : null),
                didInvalidate: false
            });

        case SELECT_NODE:
            return assign({}, state, {
                uuid: action.nodeUuid
            });

        default:
            return state;
    }
}
