import expect       from 'expect';
import nodesReducer from '../../src/reducers/nodes-reducer';
import {
    REQUEST_NODES,
    RECEIVE_NODES,
    SELECT_NODE
} from '../../src/constants/ActionTypes';


const defaultState = {
    isFetching:    false,
    uuids:         [],
    uuid:          null,
    itemsPerPage:  10,
    currentPage:   1,
    previousPage:  null,
    nextPage:      null,
    didInvalidate: true
};


describe('nodes reducer', () => {
    it('should return the initial state', () => {
        expect(nodesReducer(undefined, {}))
            .toEqual(defaultState)
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

    it('should store API results on RECEIVE_NODES action', () => {
        const nodes = [
            { uuid: 1 },
            { uuid: 2 },
            { uuid: 3 }
        ];

        const expectedState = {
            isFetching:    false,
            uuids:         nodes.map(node => node.uuid),
            itemsPerPage:  20,
            currentPage:   1,
            previousPage:  null,
            nextPage:      2,
            didInvalidate: false
        };

        expect(nodesReducer({}, {
            type:         RECEIVE_NODES,
            items:        nodes,
            itemsPerPage: 20,
            page:         1,
            previous:     0,
            next:         2
        }))
            .toEqual(expectedState)
        ;
    });

    it('should store the current node uuid on SELECT_NODE action', () => {
        const expectedState = {
            uuid: 'plouc'
        };

        expect(nodesReducer({}, {
            type:     SELECT_NODE,
            nodeUuid: 'plouc'
        }))
            .toEqual(expectedState)
        ;
    });
});
