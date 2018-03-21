import { Client, Query } from './api';
import {stringify} from 'query-string';

async function getClient() : Promise<Client> {
    let c = new Client('http://localhost:2508/api', 'v1.0');
    await c.signin('admin', 'admin');

    return c;
}

describe('Test api', () => {
    it('should authenticate', async () => {
        // with
        let c = new Client('http://localhost:2508/api', 'v1.0');

        // when
        let result = await c.signin('admin', 'admin');

        // then
        expect(result).toBe(true);
        expect(c.isAuthenticated()).toBe(true);
    });

    it('should not authenticate', async () => {
        // with
        let c = new Client('http://localhost:2508/api', 'v1.0');

        // when
        let result = await c.signin('admin', '');

        // then
        expect(result).toBe(false);
        expect(c.isAuthenticated()).toBe(false);
    });

    it('should logout', async () => {
        // with
        let c = await getClient();

        // when
        let result = await c.logout();

        // then
        expect(result).toBe(true);
        expect(c.isAuthenticated()).toBe(false);
    });

    it('should search node limit', async () => {
        // with
        let c = await getClient();

        let query = new Query(1, 2);

        let pager = await c.find(query);

        expect(pager).toBeDefined();
        expect(pager.perPage).toEqual(2);
    });

    it('should search node type media.image', async () => {
        // with
        let c = await getClient();

        let query = new Query(1, 2)
            .addType('media.image');

        let pager = await c.find(query);

        expect(pager).toBeDefined();
        expect(pager.perPage).toEqual(2);
        expect(pager.items).toHaveLength(2);

        pager.items.map((node) => {
            expect(node.type).toEqual('media.image');
        })
    });

    it('should search node with order by', async () => {
        // with
        let c = await getClient();

        let query = new Query(1, 2)
            .addOrderBy('weight', 'DESC')
            .addOrderBy('name');

        let pager = await c.find(query);

        expect(pager).toBeDefined();
        expect(pager.perPage).toEqual(2);
        expect(pager.items).toHaveLength(2);
    });

    it('should search node with meta and data', async () => {
        // with
        let c = await getClient();

        let query = new Query(1, 2)
            .addOrderBy('weight', 'DESC')
            .addOrderBy('name')
            .addData('username', 'user12')
        ;

        let pager = await c.find(query);

        expect(pager).toBeDefined();
        expect(pager.perPage).toEqual(2);
        expect(pager.items).toHaveLength(1);
    });

    it ('should find one', async () => {
        let c = await getClient();

        let query = new Query(1, 1);

        let pager = await c.find(query);

        query.addUuid(pager.items[0].uuid);

        let result = await c.findOne(query);

        expect(result).not.toBe(false);
        //expect(result.uuid).toEqual(pager.items[0].uuid);
    });

    it('test queryString (format: bracket)', () => {
        let r = stringify({
            'foo': 'bar',
            'type': ['one', 'two']
        }, {arrayFormat: 'bracket'});

        expect(r).toEqual("foo=bar&type[]=one&type[]=two");

    });

    it('test queryString (format: none)', () => {
        let r = stringify({
            'foo': 'bar',
            'type': ['one', 'two']
        }, {arrayFormat: 'none'});

        expect(r).toEqual("foo=bar&type=one&type=two");
    });
});
