import {
    GET_LIST,
    GET_ONE,
    GET_MANY,
    GET_MANY_REFERENCE,
    CREATE,
    UPDATE,
    DELETE
} from 'admin-on-rest';

import {fetchJson, queryParameters} from '../utils/fetch';

/**
 * @param {String} the base url for the api
 * @param {String} type One of the constants appearing at the top if this file, e.g. 'UPDATE'
 * @param {String} resource Name of the resource to fetch, e.g. 'posts'
 * @param {Object} params The REST request params, depending on the type
 * @returns {Object} { url, options } The HTTP request parameters
 */
const convertRESTRequestToHTTP = (apiUrl, type, resource, params) => {

    console.log( localStorage.getItem('token'));

    let url = '';

    const options = {
        headers: new Headers({
            Accept: 'application/json',
            Authorization: 'Bearer ' + localStorage.getItem('token')
        })
    };

    switch (type) {
        case GET_LIST: {
            const {page, perPage} = params.pagination;
            const {field, order} = params.sort;

            const query = {
                page: page,
                per_page: perPage,
                order_by: `${field},${order}`,
                type: resource
                //filter: JSON.stringify(params.filter),
            };

            url = `${apiUrl}/nodes?${queryParameters(query)}`;

            break;
        }
        case GET_ONE:
            url = `${apiUrl}/nodes/${params.id}`;
            break;
        case GET_MANY: {
            const query = {
                filter: JSON.stringify({id: params.ids}),
            };
            url = `${apiUrl}/nodes?${queryParameters(query)}`;
            break;
        }
        case GET_MANY_REFERENCE: {
            const {page, perPage} = params.pagination;
            const {field, order} = params.sort;
            const query = {
                sort: JSON.stringify([field, order]),
                range: JSON.stringify([(page - 1) * perPage, (page * perPage) - 1]),
                filter: JSON.stringify({...params.filter, [params.target]: params.id}),
            };
            url = `${apiUrl}/${resource}?${queryParameters(query)}`;
            break;
        }
        case UPDATE:
            url = `${apiUrl}/nodes/${params.id}`;
            options.method = 'PUT';
            options.body = JSON.stringify(params.data);
            break;
        case CREATE:
            params.data['type'] = resource;
            url = `${apiUrl}/nodes`;
            options.method = 'POST';
            options.body = JSON.stringify(params.data);
            break;
        case DELETE:
            url = `${apiUrl}/nodes/${params.id}`;
            options.method = 'DELETE';
            break;
        default:
            throw new Error(`Unsupported fetch action type ${type}`);
    }

    return {url, options};
};

/**
 * @param {Object} response HTTP response from fetch()
 * @param {String} type One of the constants appearing at the top if this file, e.g. 'UPDATE'
 * @param {String} resource Name of the resource to fetch, e.g. 'posts'
 * @param {Object} params The REST request params, depending on the type
 * @returns {Object} REST response
 */
const convertHTTPResponseToREST = (response, type, resource, params) => {
    const {json} = response;

    switch (type) {
        case GET_LIST:
            return {
                data: json.elements.map( e => {
                    e['id'] = e.uuid;

                    return e;
                }),
                total: json.next > 0 ? json.per_page * json.next : json.elements.length,
            };
        case CREATE:
            return {data: {...params.data, id: json.uuid}};
        default:
            return {data: json};
    }
};

/**
 * @param {string} apiUrl
 * @returns {function(*=, *=, *=)}
 */
export default (apiUrl) => {

    /**
     * @param {string} type Request type, e.g GET_LIST
     * @param {string} resource Resource name, e.g. "posts"
     * @param {Object} payload Request parameters. Depends on the request type
     * @returns {Promise} the Promise for a REST response
     */
    return (type, resource, params) => {
        const {url, options} = convertRESTRequestToHTTP(apiUrl, type, resource, params);
        return fetchJson(url, options)
            .then(response => convertHTTPResponseToREST(response, type, resource, params));
    };
}