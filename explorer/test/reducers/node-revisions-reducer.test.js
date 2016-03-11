import expect               from 'expect';
import nodeRevisionsReducer from '../../src/reducers/node-revisions-reducer';
import {
    REQUEST_NODE_REVISIONS,
    RECEIVE_NODE_REVISIONS,
    INVALIDATE_NODE_REVISIONS
} from '../../src/constants/ActionTypes';


describe('node revisions reducer', () => {
    it('should return the initial state', () => {
        expect(nodeRevisionsReducer(undefined, {})).toEqual({});
    });

    const uuid = 'a1b2c3';

    it('should handle the REQUEST_NODE_REVISIONS action', () => {
        expect(nodeRevisionsReducer({}, {
            type: REQUEST_NODE_REVISIONS,
            uuid
        }))
            .toEqual({
                [uuid]: {
                    didInvalidate: false,
                    isFetching:    true,
                    items:         []
                }
            })
        ;
    });

    it('should handle the RECEIVE_NODE_REVISIONS action', () => {
        const expectedRevisions = [1, 2, 3];

        expect(nodeRevisionsReducer({}, {
            type: RECEIVE_NODE_REVISIONS,
            uuid,
            items: expectedRevisions
        }))
            .toEqual({
                [uuid]: {
                    didInvalidate: false,
                    isFetching:    false,
                    items:         expectedRevisions
                }
            })
        ;
    });

    it('should handle the INVALIDATE_NODE_REVISIONS action', () => {
        expect(nodeRevisionsReducer({}, {
            type: INVALIDATE_NODE_REVISIONS,
            uuid
        }))
            .toEqual({
                [uuid]: {
                    didInvalidate: true,
                    isFetching:    false,
                    items:         []
                }
            })
        ;
    });
});
