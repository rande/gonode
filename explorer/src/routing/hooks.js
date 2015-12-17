import {
    fetchNodesIfNeeded,
    selectNode,
    fetchNodeIfNeeded
} from '../actions';


export function onEnterApp() {
}

export function onEnterNodes(store) {
    return () => {
        store.dispatch(fetchNodesIfNeeded());
    };
}

export function onEnterNode(store) {
    return (routing) => {
        const { node_uuid } = routing.params;

        store.dispatch(selectNode(node_uuid));
        store.dispatch(fetchNodeIfNeeded(node_uuid));
    };
}
