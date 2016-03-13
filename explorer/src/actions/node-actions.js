import Api         from '../Api';
import { history } from '../routing';
import {
    REQUEST_NODE,
    RECEIVE_NODE,
    REQUEST_NODE_CREATION,
    RECEIVE_NODE_CREATION,
    REQUEST_NODE_UPDATE,
    RECEIVE_NODE_UPDATE
} from '../constants/ActionTypes';
import {
    invalidateNodeRevisions,
    fetchNodeRevisionsIfNeeded
} from './node-revisions-actions';


function requestNode(uuid) {
    return {
        type: REQUEST_NODE,
        uuid
    };
}

function receiveNode(node) {
    return {
        type: RECEIVE_NODE,
        uuid: node.uuid,
        node
    };
}

function fetchNode(uuid) {
    return (dispatch, getState) => {
        dispatch(requestNode(uuid));
        Api.node(uuid, getState().security.token)
            .then(node => {
                return dispatch(receiveNode(node));
            })
        ;
    };
}

function shouldFetchNode(state, uuid) {
    const node = state.nodesByUuid[uuid];
    if (!node) {
        return true;
    }

    return node.didInvalidate;
}

export function fetchNodeIfNeeded(uuid) {
    return (dispatch, getState) => {
        if (shouldFetchNode(getState(), uuid)) {
            return dispatch(fetchNode(uuid));
        }
    };
}

function requestNodeCreation(nodeData) {
    return {
        type: REQUEST_NODE_CREATION,
        nodeData
    };
}

function receiveNodeCreation(node) {
    return {
        type: RECEIVE_NODE_CREATION,
        uuid: node.uuid,
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
        type: REQUEST_NODE_UPDATE,
        nodeData
    };
}

function receiveNodeUpdate(node) {
    return {
        type: RECEIVE_NODE_UPDATE,
        uuid: node.uuid,
        node
    };
}

export function updateNode(nodeData) {
    return (dispatch, getState) => {
        dispatch(requestNodeUpdate(nodeData));
        Api.updateNode(nodeData, getState().security.token)
            .then(node => {
                dispatch(receiveNodeUpdate(node));
                dispatch(invalidateNodeRevisions(node.uuid));
                dispatch(fetchNodeRevisionsIfNeeded(node.uuid));
                history.push(`/nodes/${node.uuid}`);
            })
        ;
    };
}
