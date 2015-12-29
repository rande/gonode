import _                 from 'lodash';
import superagent        from 'superagent';
import superagentPromise from 'superagent-promise';

const request = superagentPromise(superagent, Promise);

const JSON_HEADERS = {
    'Accept':       'application/json',
    'Content-Type': 'application/json'
};

const SERVER_BASE_URL = 'http://localhost:2508';
const API_BASE_URL = SERVER_BASE_URL + '/api/v1.0';

const Api = {
    /**
     * Authenticate user.
     *
     * @param {Object} credentials
     * @returns {Promise}
     */
    login(credentials) {
        const url = `${SERVER_BASE_URL}/login`;

        return request.post(url)
            .type('form')
            .send({
                'username': credentials.login,
                'password': credentials.password
            })
            .then(response => response.body)
            .then(data => {
                return data;
            })
            .catch((err) => {
                if (err.response) {
                    const { response } = err;
                    if (response.statusCode === 403) {
                        return response.body;
                    }
                }

                throw err;
            })
        ;
    },

    /**
     * Fetch nodes list
     *
     * @returns {Promise}
     */
    nodes(options, token = null) {
        const searchParams = [];
        if (options.perPage) {
            searchParams.push(`per_page=${options.perPage}`);
        }
        if (options.page) {
            searchParams.push(`page=${options.page}`);
        }

        const url = `${API_BASE_URL}/nodes?${searchParams.join('&')}`;

        const req = request.get(url);
        if (token !== null) {
            req.set('Authorization', `Bearer ${token}`);
        }

        return req
            .then(response => response.body)
        ;
    },

    /**
     * Fetch node by uuid
     *
     * @returns {Promise}
     */
    node(uuid, token = null) {
        const url = `${API_BASE_URL}/nodes/${uuid}`;

        const req = request.get(url);
        if (token !== null) {
            req.set('Authorization', `Bearer ${token}`);
        }

        return req
            .then(response => response.body)
        ;
    },

    createNode(nodeData, token = null) {
        const url = `${API_BASE_URL}/nodes`;

        const req = request.post(url);
        if (token !== null) {
            req.set('Authorization', `Bearer ${token}`);
        }

        return req
            .send(nodeData)
            .then(response => response.json())
            .then(node => {
                return node;
            })
        ;
    },

    updateNode(nodeData, token = null) {
        const url = `${API_BASE_URL}/nodes/${nodeData.uuid}`;

        const req = request.put(url);
        if (token !== null) {
            req.set('Authorization', `Bearer ${token}`);
        }

        return req
            .send(nodeData)
            .then(response => response.body)
            .then(node => {
                return node;
            })
        ;
    }
};


export default Api;