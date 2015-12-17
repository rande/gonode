import _ from 'lodash';
import 'whatwg-fetch';

const JSON_HEADERS = {
    'Accept':       'application/json',
    'Content-Type': 'application/json'
};

const API_BASE_URL = 'http://localhost:2405/';

const Api = {
    /**
     * Fetch nodes list
     *
     * @returns {Promise}
     */
    nodes(options) {
        const searchParams = [];
        if (options.perPage) {
            searchParams.push(`per_page=${options.perPage}`)
        }

        const url = `${API_BASE_URL}nodes?${searchParams.join('&')}`;

        return fetch(url)
            .then(response => response.json())
            .then(json => {
                return json.elements;
            })
        ;
    },

    /**
     * Fetch node by uuid
     *
     * @returns {Promise}
     */
    node(uuid) {
        return fetch(`${API_BASE_URL}nodes/${uuid}`)
            .then(response => response.json())
            .then(node => {
                return node;
            })
        ;
    },

    createNode(nodeData) {
        return fetch(`${API_BASE_URL}nodes`, {
            method: 'post',
            body:   JSON.stringify(nodeData)
        })
            .then(response => response.json())
            .then(node => {
                console.log('API.CREATE_NODE', node);
                return node;
            })
        ;
    }
};


export default Api;