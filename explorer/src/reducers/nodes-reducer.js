import _               from 'lodash';
import { UPDATE_PATH } from 'redux-simple-router';
import * as types      from '../constants/ActionTypes';


export default function nodes(state = {
    isFetching:    false,
    items:         [],
    itemsPerPage:  10,
    currentUuid:   null,
    didInvalidate: true
}, action) {
    switch (action.type) {
        case types.REQUEST_NODES:
            return Object.assign({}, state, {
                isFetching: true
            });

        case types.RECEIVE_NODES:
            return Object.assign({}, state, {
                isFetching: false,
                items:      action.items
            });

        case types.SELECT_NODE:
            return Object.assign({}, state, {
                currentUuid: action.nodeUuid
            });

        case types.RECEIVE_NODE_CREATION:
            return Object.assign({}, state, {
                didInvalidate: true
            });

        case types.SET_NODES_PAGER_OPTIONS:
            return Object.assign({}, state, {
                itemsPerPage:  action.itemsPerPage,
                didInvalidate: true
            });

        default:
            return state;
    }
}
