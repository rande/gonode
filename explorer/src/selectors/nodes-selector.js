import { createSelector } from 'reselect';

const assign = Object.assign || require('object.assign');


/**
 * Input selectors
 */
const nodesStateSelector           = state => state.nodes;
const uuidSelector                 = state => state.nodes.uuid;
const nodesByUuidSelector          = state => state.nodesByUuid;
const nodesRevisionsByUuidSelector = state => state.nodesRevisionsByUuid;


/**
 * Combined selectors.
 */
export const nodeSelector = createSelector(
    uuidSelector,
    nodesByUuidSelector,
    (uuid, nodesByUuid) => (nodesByUuid[uuid] ? { node: nodesByUuid[uuid] } : null)
);

export const nodesSelector = createSelector(
    nodesStateSelector,
    nodesByUuidSelector,
    (nodes, nodesByUuid) => {
        let items = [];
        if (nodes.uuids.length > 0) {
            items = [];
            nodes.uuids.forEach(uuid => {
                if (nodesByUuid[uuid] && nodesByUuid[uuid].node) {
                    items.push(nodesByUuid[uuid].node);
                }
            });
        }

        return assign({}, nodes, { nodes: items });
    }
);

export const nodeRevisionsSelector = createSelector(
    uuidSelector,
    nodeSelector,
    nodesRevisionsByUuidSelector,
    (uuid, { node }, nodesRevisionsByUuid) => {
        const output = {
            uuid,
            isFetching: true,
            revisions:  [],
            hasMore:    false,
            nextPage:   0,
            node:       node.node ? node.node : null
        };

        if (nodesRevisionsByUuid[uuid]) {
            const revisions = nodesRevisionsByUuid[uuid];

            output.isFetching = revisions.isFetching;
            output.revisions  = [];
            output.nextPage   = revisions.nextPage;

            if (revisions.ids.length > 0) {
                revisions.ids.forEach(id => {
                    if (revisions.byRevisionId[id] && revisions.byRevisionId[id].revision) {
                        output.revisions.push(revisions.byRevisionId[id].revision);
                    }
                });
            }
        }

        return output;
    }
);

export const nodeRevisionSelector = createSelector(
    uuidSelector,
    nodesRevisionsByUuidSelector,
    (uuid, nodesRevisionsByUuid) => {
        const output = {
            nodeUuid:   uuid,
            revisionId: null,
            isFetching: false,
            revision:   null
        };

        if (nodesRevisionsByUuid[uuid]) {
            const revisions = nodesRevisionsByUuid[uuid];

            if (revisions.id !== null) {
                output.revisionId = revisions.id;

                if (revisions.byRevisionId[revisions.id]) {
                    output.isFetching = revisions.byRevisionId[revisions.id].isFetching;
                    output.revision   = revisions.byRevisionId[revisions.id].revision;
                }
            }
        }

        return output;
    }
);
