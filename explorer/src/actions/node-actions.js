import * as types  from '../constants/ActionTypes';
import Api         from '../Api';
import { history } from '../routing';


function receiveNode(node) {
    return {
        type:     types.RECEIVE_NODE,
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
    return (dispatch, getState) => {
        dispatch(requestNode(nodeUuid));
        Api.node(nodeUuid, getState().security.token)
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

    return node.didInvalidate;
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
        Api.createNode(nodeData, getState().security.token)
            .then(node => {
                dispatch(receiveNodeCreation(node));
                history.push('/nodes');
            })
        ;
    };
}

function requestNodeUpdate(nodeData) {
    return {
        type: types.REQUEST_NODE_UPDATE,
        nodeData
    };
}

function receiveNodeUpdate(node) {
    return {
        type:     types.RECEIVE_NODE_UPDATE,
        nodeUuid: node.uuid,
        node
    };
}

export function updateNode(nodeData) {
    return (dispatch, getState) => {
        dispatch(requestNodeUpdate(nodeData));
        Api.updateNode(nodeData, getState().security.token)
            .then(node => {
                dispatch(receiveNodeUpdate(node));
                history.push(`/nodes/${node.uuid}`);
            })
        ;
    };
}
