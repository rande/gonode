import { combineReducers }  from 'redux';
import security             from './security-reducer';
import nodes                from './nodes-reducer';
import nodesByUuid          from './nodes-by-uuid-reducer';
import nodesRevisionsByUuid from './nodes-revisions-by-uuid-reducer';

export default {
    security,
    nodes,
    nodesByUuid,
    nodesRevisionsByUuid
};
