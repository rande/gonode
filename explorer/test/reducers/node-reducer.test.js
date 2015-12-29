import expect      from 'expect';
import nodeReducer from '../../src/reducers/node-reducer';
import * as types  from '../../src/constants/ActionTypes';


describe('node reducer', () => {
    it('should return the initial state', () => {
        expect(nodeReducer(undefined, {})).toEqual({});
    });

    const nodeUuid = 'a1b2c3';

    it('should handle the REQUEST_NODE action', () => {
        expect(nodeReducer({}, {
            type:     types.REQUEST_NODE,
            nodeUuid
        }))
            .toEqual({
                [nodeUuid]: {
                    didInvalidate: false,
                    isFetching:    true,
                    node:          null
                }
            })
        ;
    });

    it('should handle the RECEIVE_NODE action', () => {
        expect(nodeReducer({}, {
            type: types.RECEIVE_NODE,
            nodeUuid,
            node: {}
        }))
            .toEqual({
                [nodeUuid]: {
                    didInvalidate: false,
                    isFetching:    false,
                    node:          {}
                }
            })
        ;
    });

    it('should handle the RECEIVE_NODE_UPDATE action', () => {
        expect(nodeReducer({}, {
            type:     types.RECEIVE_NODE_UPDATE,
            nodeUuid,
            node:     {}
        }))
            .toEqual({
                [nodeUuid]: {
                    didInvalidate: true,
                    isFetching:    false,
                    node:          null
                }
            })
        ;
    });
});