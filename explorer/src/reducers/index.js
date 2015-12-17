import { combineReducers } from 'redux';
import nodes               from './nodes-reducer';
import nodesByUuid         from './node-reducer';

export default {
    nodes,
    nodesByUuid
};