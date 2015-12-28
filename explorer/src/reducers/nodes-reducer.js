import _               from 'lodash';
import { UPDATE_PATH } from 'redux-simple-router';
import * as types      from '../constants/ActionTypes';

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
        case types.REQUEST_NODES:
            return assign({}, state, {
                isFetching:   true,
                currentPage:  action.page,
                itemsPerPage: action.perPage
            });

        case types.RECEIVE_NODES:
            return assign({}, state, {
                isFetching:    false,
                items:         action.items,
                itemsPerPage:  action.itemsPerPage,
                currentPage:   action.page,
                previousPage:  (action.previous !== 0 ? action.previous : null),
                nextPage:      (action.next     !== 0 ? action.next     : null),
                didInvalidate: false
            });

        case types.SELECT_NODE:
            return assign({}, state, {
                currentUuid: action.nodeUuid
            });

        case types.RECEIVE_NODE_CREATION:
            return assign({}, state, {
                didInvalidate: true
            });

        default:
            return state;
    }
}
