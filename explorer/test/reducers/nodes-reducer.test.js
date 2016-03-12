import expect       from 'expect';
import nodesReducer from '../../src/reducers/nodes-reducer';
import {
    REQUEST_NODES,
    RECEIVE_NODES,
    SELECT_NODE,
    RECEIVE_NODE_CREATION,
    RECEIVE_NODE_UPDATE
} from '../../src/constants/ActionTypes';


describe('nodes reducer', () => {
    it('should return the initial state', () => {
        expect(nodesReducer(undefined, {}))
            .toEqual({
                didInvalidate: true,
                isFetching:    false,
                items:         [],
                itemsPerPage:  10,
                previousPage:  null,
                currentPage:   1,
                nextPage:      null,
                currentUuid:   null
            })
        ;
    });

    it('should handle the REQUEST_NODES action', () => {
        expect(nodesReducer({}, {
            type:    REQUEST_NODES,
            perPage: 10,
            page:    1
        }))
            .toEqual({
                isFetching:   true,
                itemsPerPage: 10,
                currentPage:  1
            })
        ;
    });

    it('should store API results when receiving a RECEIVE_NODES action', () => {
        expect(nodesReducer({}, {
            type:         RECEIVE_NODES,
            items:        [],
            itemsPerPage: 20,
            page:         1,
            previous:     0,
            next:         2
        }))
            .toEqual({
                isFetching:    false,
                items:         [],
                itemsPerPage:  20,
                currentPage:   1,
                previousPage:  null,
                nextPage:      2,
                didInvalidate: false
            })
        ;
    });

    it('should store the current node uuid when receiving a SELECT_NODE action', () => {
        expect(nodesReducer({}, {
            type:     SELECT_NODE,
            nodeUuid: 'plouc'
        }))
            .toEqual({
                currentUuid: 'plouc'
            })
        ;
    });

    it('should set current state as invalid when receiving a RECEIVE_NODE_CREATION action', () => {
        expect(nodesReducer({}, {
            type: RECEIVE_NODE_CREATION,
            node: {}
        }))
            .toEqual({
                didInvalidate: true
            })
        ;
    });

    it('should replace an already fetched node with the updated one when receiving a RECEIVE_NODE_UPDATE action', () => {
        expect(nodesReducer({
            items: [{ uuid: 1, name: 'old' }]
        }, {
            type: RECEIVE_NODE_UPDATE,
            node: { uuid: 1, name: 'new' }
        }))
            .toEqual({
                items: [{ uuid: 1, name: 'new' }]
            })
        ;
    });
});
