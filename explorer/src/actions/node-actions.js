import * as types             from '../constants/ActionTypes';
import Api                    from '../Api';
import { fetchNodesIfNeeded } from './nodes-actions';
import { history }            from '../routing';


function receiveNode(node) {
    return {
        type: types.RECEIVE_NODE,
        nodeUuid: node.uuid,
        node
    };
}

function requestNode(nodeUuid) {
    return {
        type: types.REQUEST_NODE,
        nodeUuid
    };
}

function fetchNode(nodeUuid) {
    return dispatch => {
        dispatch(requestNode(nodeUuid));
        Api.node(nodeUuid)
            .then(node => {
                return dispatch(receiveNode(node));
            })
        ;
    };
}

function shouldFetchNode(state, nodeUuid) {
    const node = state.nodesByUuid[nodeUuid];
    if (!node || node.isFetching) {
        return true;
    }

    return false;
}

export function fetchNodeIfNeeded(nodeUuid) {
    return (dispatch, getState) => {
        if (shouldFetchNode(getState(), nodeUuid)) {
            return dispatch(fetchNode(nodeUuid));
        }
    };
}

function requestNodeCreation(nodeData) {
    return {
        type: types.REQUEST_NODE_CREATION,
        nodeData
    };
}

function receiveNodeCreation(node) {
    return {
        type: types.RECEIVE_NODE_CREATION,
        node
    };
}

export function createNode(nodeData) {
    return (dispatch, getState) => {
        dispatch(requestNodeCreation(nodeData));
        Api.createNode(nodeData)
            .then(node => {
                dispatch(receiveNodeCreation(node));
                fetchNodesIfNeeded()(dispatch, getState);
                history.push('/nodes');
            })
        ;
    };
}
