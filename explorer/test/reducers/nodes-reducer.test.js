import expect       from 'expect';
import nodesReducer from '../../src/reducers/nodes-reducer';
import * as types   from '../../src/constants/ActionTypes';


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
            type:    types.REQUEST_NODES,
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

    it('should handle the RECEIVE_NODES action', () => {
        expect(nodesReducer({}, {
            type:         types.RECEIVE_NODES,
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

    it('should handle the SELECT_NODE action', () => {
        expect(nodesReducer({}, {
            type:     types.SELECT_NODE,
            nodeUuid: 'plouc'
        }))
            .toEqual({
                currentUuid: 'plouc'
            })
        ;
    });

    it('should handle the RECEIVE_NODE_CREATION action', () => {
        expect(nodesReducer({}, {
            type: types.RECEIVE_NODE_CREATION,
            node: {}
        }))
            .toEqual({
                didInvalidate: true
            })
        ;
    });

    it('should handle the RECEIVE_NODE_UPDATE action', () => {
        expect(nodesReducer({}, {
            type: types.RECEIVE_NODE_UPDATE,
            node: {}
        }))
            .toEqual({
                didInvalidate: true
            })
        ;
    });
});