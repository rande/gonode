import expect             from 'expect';
import nodesByUuidReducer from '../../src/reducers/nodes-by-uuid-reducer';
import {
    RECEIVE_NODES,
    REQUEST_NODE,
    RECEIVE_NODE,
    RECEIVE_NODE_CREATION,
    RECEIVE_NODE_UPDATE
} from '../../src/constants/ActionTypes';


describe('nodes by uuid reducer', () => {
    it('should return the initial state on unrelated action', () => {
        expect(nodesByUuidReducer(undefined, {})).toEqual({});
    });

    it(`should store an 'augmented version' of each nodes on RECEIVE_NODES action`, () => {
        const sampleNodes = [
            { uuid: 1 },
            { uuid: 2 },
            { uuid: 3 }
        ];

        const expectedState = {};
        sampleNodes.forEach(sampleNode => {
            expectedState[sampleNode.uuid] = {
                isFetching:    false,
                node:          sampleNode,
                didInvalidate: false
            }
        });

        expect(nodesByUuidReducer({}, {
            type:  RECEIVE_NODES,
            items: sampleNodes
        })).toEqual(expectedState);
    });

    it('should set the node as fetching on REQUEST_NODE action', () => {
        const uuid          = 'test_uuid';
        const expectedState = {
            [uuid]: {
                didInvalidate: false,
                isFetching:    true,
                node:          null
            }
        };

        expect(nodesByUuidReducer({}, {
            type: REQUEST_NODE,
            uuid
        })).toEqual(expectedState);
    });

    it('should store node and switch fetching/invalid flags accordingly on RECEIVE_NODE action', () => {
        const sampleNode = {
            uuid: 'test_uuid'
        };
        const expectedState = {
            [sampleNode.uuid]: {
                didInvalidate: false,
                isFetching:    false,
                node:          sampleNode
            }
        };

        expect(nodesByUuidReducer({
            [sampleNode.uuid]: {
                didInvalidate: true,
                isFetching:    true,
                node:          null
            }
        }, {
            type: RECEIVE_NODE,
            uuid: sampleNode.uuid,
            node: sampleNode
        })).toEqual(expectedState);
    });

    it('should store node and switch fetching/invalid flags accordingly on RECEIVE_NODE_CREATION action', () => {
        const sampleNode = {
            uuid: 'test_uuid'
        };
        const expectedState = {
            [sampleNode.uuid]: {
                didInvalidate: false,
                isFetching:    false,
                node:          sampleNode
            }
        };

        expect(nodesByUuidReducer({
            [sampleNode.uuid]: {
                didInvalidate: true,
                isFetching:    true,
                node:          null
            }
        }, {
            type: RECEIVE_NODE_CREATION,
            uuid: sampleNode.uuid,
            node: sampleNode
        })).toEqual(expectedState);
    });

    it('should store node and switch fetching/invalid flags accordingly on RECEIVE_NODE_UPDATE action', () => {
        const sampleNode = {
            uuid: 'test_uuid'
        };
        const expectedState = {
            [sampleNode.uuid]: {
                didInvalidate: false,
                isFetching:    false,
                node:          sampleNode
            }
        };

        expect(nodesByUuidReducer({
            [sampleNode.uuid]: {
                didInvalidate: true,
                isFetching:    true,
                node:          null
            }
        }, {
            type: RECEIVE_NODE_UPDATE,
            uuid: sampleNode.uuid,
            node: sampleNode
        })).toEqual(expectedState);
    });
});
