import * as types from '../constants/ActionTypes';
import Api        from '../Api';


function receiveNodes({
    elements,
    per_page,
    page,
    previous,
    next
}) {
    return {
        type:         types.RECEIVE_NODES,
        items:        elements,
        itemsPerPage: per_page,
        page,
        previous,
        next
    };
}

function requestNodes({ perPage = 10, page = 1 }) {
    return {
        type: types.REQUEST_NODES,
        perPage,
        page
    };
}

export function selectNode(nodeUuid) {
    return {
        type: types.SELECT_NODE,
        nodeUuid
    };
}

function fetchNodes(params) {
    return (dispatch, getState) => {
        dispatch(requestNodes(params));

        Api.nodes({
            perPage: params.perPage,
            page:    params.page
        }, getState().security.token)
            .then(nodes => {
                dispatch(receiveNodes(nodes));
            })
        ;
    };
}

function shouldFetchNodes(params, state) {
    const { nodes } = state;
    if (nodes.isFetching) {
        return false;
    }

    if (nodes.currentPage !== params.page) {
        return true;
    }

    if (nodes.itemsPerPage !== params.perPage) {
        return true;
    }

    return nodes.didInvalidate;
}

export function fetchNodesIfNeeded(params) {
    return (dispatch, getState) => {
        if (shouldFetchNodes(params, getState())) {
            return dispatch(fetchNodes(params));
        }
    };
}
