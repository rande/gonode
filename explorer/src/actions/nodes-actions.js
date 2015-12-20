import * as types from '../constants/ActionTypes';
import Api        from '../Api';


function receiveNodes(nodes) {
    return {
        type:  types.RECEIVE_NODES,
        items: nodes
    };
}

function requestNodes() {
    return {
        type: types.REQUEST_NODES
    };
}

export function selectNode(nodeUuid) {
    return {
        type: types.SELECT_NODE,
        nodeUuid
    };
}

function fetchNodes() {
    return (dispatch, getState) => {
        dispatch(requestNodes());

        const { nodes: {
            itemsPerPage
        } } = getState();

        Api.nodes({
            perPage: itemsPerPage
        }, getState().security.token)
            .then(nodes => {
                dispatch(receiveNodes(nodes));
            })
        ;
    };
}

function shouldFetchNodes(state) {
    const { nodes } = state;
    if (nodes.isFetching) {
        return false;
    }

    return nodes.didInvalidate;
}

export function setNodesPagerOptions({ itemsPerPage }) {
    return (dispatch, getState) => {
        dispatch({
            type: types.SET_NODES_PAGER_OPTIONS,
            itemsPerPage
        });
        return fetchNodesIfNeeded()(dispatch, getState);
    };
}

export function fetchNodesIfNeeded() {
    return (dispatch, getState) => {
        if (shouldFetchNodes(getState())) {
            return dispatch(fetchNodes());
        }
    };
}
