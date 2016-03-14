import expect               from 'expect';
import nodeRevisionsByMonth from '../../src/selectors/node-revisions-by-month';


describe('nodeRevisionsByMonth()', () => {
    it('should return an empty array when given an empty array of revisions', () => {
        expect(nodeRevisionsByMonth([])).toEqual([]);
    });

    it('should return buckets with all revisions grouped by month', () => {
        const sampleRevisions = [
            { id: 5, updated_at: '03-13-2016' },
            { id: 4, updated_at: '11-14-2015' },
            { id: 3, updated_at: '11-10-2015' },
            { id: 2, updated_at: '09-02-2015' },
            { id: 1, updated_at: '09-01-2015' }
        ];

        const expectedResult = [
            { year: 2016, month: 2,  items: [sampleRevisions[0]] },
            { year: 2015, month: 10, items: [sampleRevisions[1], sampleRevisions[2]] },
            { year: 2015, month: 8,  items: [sampleRevisions[3], sampleRevisions[4]] }
        ];

        expect(nodeRevisionsByMonth(sampleRevisions)).toEqual(expectedResult);
    });

    it('should build multiple buckets for identical months with different year', () => {
        const sampleRevisions = [
            { id: 2, updated_at: '09-02-2016' },
            { id: 1, updated_at: '09-02-2015' }
        ];

        const expectedResult = [
            { year: 2016, month: 8, items: [sampleRevisions[0]] },
            { year: 2015, month: 8, items: [sampleRevisions[1]] }
        ];

        expect(nodeRevisionsByMonth(sampleRevisions)).toEqual(expectedResult);
    });
});
