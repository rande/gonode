import expect               from 'expect';
import nodesRevisionsByUuid from '../../src/reducers/nodes-revisions-by-uuid-reducer';
import {
    REQUEST_NODE_REVISIONS,
    RECEIVE_NODE_REVISIONS,
    INVALIDATE_NODE_REVISIONS,
    SELECT_NODE_REVISION
} from '../../src/constants/ActionTypes';


describe('nodes revisions by uuid reducer', () => {
    it('should return the initial state on unrelated action', () => {
        expect(nodesRevisionsByUuid(undefined, {})).toEqual({});
    });

    it('should set revisions as fetching on REQUEST_NODE_REVISIONS action', () => {
        const uuid          = 'test_uuid';
        const expectedState = {
            [uuid]: { isFetching: true }
        };

        expect(nodesRevisionsByUuid({
            [uuid]: { isFetching: false }
        }, {
            type: REQUEST_NODE_REVISIONS,
            uuid
        })).toEqual(expectedState);
    });

    it(`should store an 'augmented version' of each revisions and switch fetching/invalid flags accordingly on RECEIVE_NODE_REVISIONS action`, () => {
        const uuid            = 'test_uuid';
        const sampleRevisions = [
            { revision: 1 },
            { revision: 2 },
            { revision: 3 }
        ];

        const byRevisionId    = {};
        sampleRevisions.forEach(sampleRevision => {
            byRevisionId[sampleRevision.revision] = {
                isFetching:    false,
                didInvalidate: false,
                revision:      sampleRevision
            };
        });

        const expectedState   = {
            [uuid]: {
                isFetching:    false,
                didInvalidate: false,
                ids:           sampleRevisions.map(sampleRevision => sampleRevision.revision),
                page:          1,
                nextPage:      0,
                byRevisionId
            }
        };

        expect(nodesRevisionsByUuid({
            [uuid]: {
                isFetching:    true,
                didInvalidate: true,
                page:          1
            }
        }, {
            type: RECEIVE_NODE_REVISIONS,
            uuid,
            page:     1,
            nextPage: 0,
            items:    sampleRevisions
        })).toEqual(expectedState);
    });

    it('should append revisions instead of replacing it if page is greater than current one on RECEIVE_NODE_REVISIONS action', () => {
        const uuid            = 'test_uuid';
        const sampleRevisions = [
            { revision: 3 },
            { revision: 4 },
            { revision: 5 }
        ];

        const byRevisionId    = {};
        sampleRevisions.forEach(sampleRevision => {
            byRevisionId[sampleRevision.revision] = {
                isFetching:    false,
                didInvalidate: false,
                revision:      sampleRevision
            };
        });

        const expectedState   = {
            [uuid]: {
                isFetching:    false,
                didInvalidate: false,
                ids:           [1, 2, 3, 4, 5],
                page:          2,
                nextPage:      0,
                byRevisionId
            }
        };

        expect(nodesRevisionsByUuid({
            [uuid]: {
                ids:           [1, 2],
                isFetching:    true,
                didInvalidate: true,
                page:          1
            }
        }, {
            type: RECEIVE_NODE_REVISIONS,
            uuid,
            page:     2,
            nextPage: 0,
            items:    sampleRevisions
        })).toEqual(expectedState);
    });

    it('should mark the revisions as invalids on INVALIDATE_NODE_REVISIONS action', () => {
        const uuid = 'test_uuid';
        const expectedState = {
            [uuid]: {
                didInvalidate: true
            }
        };

        expect(nodesRevisionsByUuid({
            [uuid]: {
                didInvalidate: false
            }
        }, {
            type: INVALIDATE_NODE_REVISIONS,
            uuid
        })).toEqual(expectedState);
    });
});
