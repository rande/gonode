import expect         from 'expect';
import { selectNode } from '../../src/actions/nodes-actions';
import * as types     from '../../src/constants/ActionTypes';


describe('nodes actions', () => {
    describe('select node', () => {
        it('should change the current node uuid', () => {
            const nodeUuid = 'plouc';
            const expectedResult = {
                type: types.SELECT_NODE,
                nodeUuid
            };

            expect(selectNode(nodeUuid)).toEqual(expectedResult);
        });
    });
});