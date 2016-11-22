// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

import { push } from 'react-router-redux'

const SEARCHING_NODES = 'SEARCHING_NODES';
const LIST_NODES = 'LIST_NODES';
const LIST_NODE_SERVICES = 'LIST_NODE_SERVICES';

const LOAD_NODE = 'LOAD_NODE';
const LOADING_NODE = 'LOADING_NODE';


function getHeader(state) {
    const headers = new Headers();
    headers.append('Authorization', `Bearer ${state.userApp.auth}`);

    return headers;
}

export function getNodeServices() {

    return (dispatch, getState) => {
        const init = {
            method: 'GET',
            headers: getHeader(getState())
        };

        fetch(`/api/v1.0/handlers/node`, init)
            .then((res) => {
               return res.json()
            })
            .then((data => {
                dispatch({
                    services: data,
                    type: LIST_NODE_SERVICES,
                })
            }));
    }
}

export function searchNodes(params = {}, page = 1, per_page = defaultState.search.per_page) {

    return (dispatch, getState) => {
        const state = getState();

        if (params == null) { // load current params from the state
            params = state.nodeApp.search.params;
        }

        if (page == null) {
            page = state.nodeApp.search.page;
        }

        if (per_page == null) {
            per_page = state.nodeApp.search.per_page;
        }

        const init = {
            method: 'GET',
            headers: getHeader(state)
        };

        const query = new URLSearchParams();

        for (var k in params) {
            if (params[k] == null) {
                continue;
            }

            query.append(k, params[k]);
        }

        query.append('per_page', per_page);
        query.append('page', page);

        fetch(`/api/v1.0/nodes?${query.toString()}`, init)
            .then((res) => {
               return res.json()
            })
            .then((data => {
                dispatch({
                    ...data,
                    type: LIST_NODES,
                })
            }));

        dispatch({
            type: SEARCHING_NODES,
            params: params,
            page: page,
            per_page: per_page
        });
    }
}

export function loadNode(uuid) {
    return (dispatch, getState) => {

        const init = {
            method: 'GET',
            headers: getHeader(getState())
        };

        fetch(`/api/v1.0/nodes/${uuid}`, init)
            .then((res) => {
               return res.json()
            })
            .then((data => {
                dispatch({
                    node: data,
                    type: LOAD_NODE,
                });

                dispatch(push(`/nodes/${uuid}`));
            }));

        dispatch({
            type: LOADING_NODE,
            uuid: uuid,
        });
    }
}

const defaultState = {
    state: LIST_NODES,
    result: {
        elements: [],
        per_page: 16,
        page: 1,
        next: 0,
        previous: 0,
    },
    search: {
        params: {
            type: null,
            enabled: null
        },
        per_page: 24,
        page: 1
    },
    services: [],
    node: {}
};

export function nodeApp(state = defaultState, action) {
    switch (action.type) {

        case SEARCHING_NODES:
            return {
                ...state,
                state: SEARCHING_NODES,
                search: {
                    params: action.params,
                    per_page: action.per_page,
                    page: action.page,
                }
            };

        case LIST_NODES: {
            return {
                ...state,
                state: LIST_NODES,
                result: {
                    elements: action.elements,
                    per_page: action.per_page,
                    page: action.page,
                    next: action.next,
                    previous: action.previous
                }
            };
        }

        case LIST_NODE_SERVICES: {
            return {
                ...state,
                services: action.services,
            }
        }

        case LOADING_NODE: {
            return {
                ...state,
                state: LOADING_NODE,
            }
        }

        case LOAD_NODE: {
            return {
                ...state,
                node: action.node,
                state: LOAD_NODE,
            }
        }

        default:
            return state
    }
}