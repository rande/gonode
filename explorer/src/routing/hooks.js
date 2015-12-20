import {
    fetchNodesIfNeeded,
    selectNode,
    fetchNodeIfNeeded,
    logout
} from '../actions';
import history from './history';


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

export function onEnterNodes(store) {
    return ensureAuthenticated(store, function () {
        store.dispatch(fetchNodesIfNeeded());
    });
}

export function onEnterNode(store) {
    return ensureAuthenticated(store, function (nextState) {
        const { node_uuid } = nextState.params;

        store.dispatch(selectNode(node_uuid));
        store.dispatch(fetchNodeIfNeeded(node_uuid));
    });
}
