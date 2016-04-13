import {
    REQUEST_NODE_REVISIONS,
    RECEIVE_NODE_REVISIONS,
    INVALIDATE_NODE_REVISIONS,
    SELECT_NODE_REVISION
} from '../constants/ActionTypes';

const assign = Object.assign || require('object.assign');

function nodeRevision(state = {
    isFetching:    false,
    revision:      null,
    didInvalidate: false
}, action) {
    switch (action.type) {
        case RECEIVE_NODE_REVISIONS:
            return assign({}, state, {
                isFetching:    false,
                revision:      action.revision,
                didInvalidate: false
            });

        default:
            return state;
    }
}

function nodeRevisions(state = {
    isFetching:    false,
    ids:           [],
    id:            null,
    byRevisionId:  {},
    page:          1,
    nextPage:      0,
    didInvalidate: false
}, action) {
    switch (action.type) {
        case REQUEST_NODE_REVISIONS:
            return assign({}, state, {
                isFetching: true
            });

        case RECEIVE_NODE_REVISIONS:
            const { ids, byRevisionId } = state;

            const receivedIds = action.items.map(item => item.revision);
            const newObjects  = assign({}, byRevisionId, action.items.reduce((newRevisions, revision) => {
                newRevisions[revision.revision] = nodeRevision(undefined, assign({}, action, {
                    revision
                }));

                return newRevisions;
            }, {}));

            let resultingIds = [];
            if (action.page === 1) {
                resultingIds = receivedIds;
            } else {
                resultingIds = ids.concat(receivedIds);
            }

            return assign({}, state, {
                byRevisionId:  newObjects,
                ids:           resultingIds,
                page:          action.page,
                nextPage:      action.nextPage,
                isFetching:    false,
                didInvalidate: false
            });

        case INVALIDATE_NODE_REVISIONS:
            return assign({}, state, {
                didInvalidate: true
            });

        case SELECT_NODE_REVISION:
            return assign({}, state, {
                id: action.id
            });

        default:
            return state;
    }
}


export default function nodesRevisionsByUuid(state = {}, action) {
    switch (action.type) {
        case REQUEST_NODE_REVISIONS:
        case RECEIVE_NODE_REVISIONS:
        case INVALIDATE_NODE_REVISIONS:
        case SELECT_NODE_REVISION:
            return assign({}, state, {
                [action.uuid]: nodeRevisions(state[action.uuid], action)
            });

        default:
            return state;
    }
}
