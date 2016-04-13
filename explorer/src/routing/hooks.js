import {
    selectNode,
    fetchNodeIfNeeded,
    fetchNodeRevisionsIfNeeded,
    selectNodeRevision,
    fetchNodeRevisionIfNeeded,
    logout
} from '../actions';


export function ensureAuthenticated(store, execIfAuthenticated = null) {
    return function (nextState, replaceState) {
        const { security: { token }} = store.getState();
        if (!token) {
            replaceState(null, '/login');
        } else {
            execIfAuthenticated && execIfAuthenticated(nextState, replaceState);
        }
    };
}

export function onEnterLogout(store) {
    return function () {
        store.dispatch(logout());
    };
}

export function onEnterApp(store) {
    return function () {};
}

export function onEnterNode(store) {
    return ensureAuthenticated(store, nextState => {
        const { node_uuid } = nextState.params;

        store.dispatch(selectNode(node_uuid));
        store.dispatch(fetchNodeIfNeeded(node_uuid));
        store.dispatch(fetchNodeRevisionsIfNeeded(node_uuid));
    });
}

export function onEnterNodeRevision(store) {
    return ensureAuthenticated(store, nextState => {
        const { node_uuid, revision_id } = nextState.params;

        store.dispatch(selectNode(node_uuid));
        store.dispatch(fetchNodeIfNeeded(node_uuid));

        store.dispatch(fetchNodeRevisionsIfNeeded(node_uuid));

        const revisionId = parseInt(revision_id);
        store.dispatch(selectNodeRevision(node_uuid, revisionId));
        store.dispatch(fetchNodeRevisionIfNeeded(node_uuid, revisionId));
    });
}
