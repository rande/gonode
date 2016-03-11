import { combineReducers }  from 'redux';
import security             from './security-reducer';
import nodes                from './nodes-reducer';
import nodesByUuid          from './node-reducer';
import nodesRevisionsByUuid from './node-revisions-reducer';

export default {
    security,
    nodes,
    nodesByUuid,
    nodesRevisionsByUuid
};
