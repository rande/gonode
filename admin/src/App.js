import React from 'react';
import {Admin, Resource} from 'admin-on-rest';

import {NodeList, NodeShow, NodeCreate} from './components/nodes';
import {UserList, UserShow, UserCreate} from './components/users';

import apiClient from './services/api';
import authClient from './services/authClient';

import PostIcon from 'material-ui/svg-icons/action/book';
import UserIcon from 'material-ui/svg-icons/social/group';
import ImageIcon from 'material-ui/svg-icons/image/image';
import ListIcon from 'material-ui/svg-icons/action/list';
import VideoIcon from 'material-ui/svg-icons/action/theaters';


import {CRUD_CREATE_FAILURE} from "admin-on-rest";
import {stopSubmit} from 'redux-form';
import {put, takeEvery} from "redux-saga/effects";

export function* errorSagas() {
    yield takeEvery(CRUD_CREATE_FAILURE, crudCreateFailure);
}

function makeViolations(flatErrors) {
    let errors = {}
    for (let field in flatErrors) {
        field.split('.').reduce((o, v, i, s) => {
            if (s.length - 1 === i) { // last element, set the value
                o[v] = flatErrors[field].reduce((o, v) => `${o} / ${v}`);
            }

            if (!o[v]) {
                o[v] = {};
            }

            return o[v];
        }, errors)
    }

    return errors;
}

export function* crudCreateFailure(action) {
    yield put(stopSubmit('record-form', makeViolations(action.payload)));
}


const App = () => (
    <Admin title="GoNode" authClient={authClient} restClient={apiClient('/api/v1.0')} customSagas={[errorSagas]}>
        <Resource name="core.index" list={NodeList} create={NodeCreate} show={NodeShow} options={{label: "Index"}}
                  icon={ListIcon}/>
        <Resource name="blog.post" list={NodeList} create={NodeCreate} show={NodeShow} options={{label: "Posts"}}
                  icon={PostIcon}/>
        <Resource name="core.user" list={UserList} create={UserCreate} show={UserShow} options={{label: "Users"}}
                  icon={UserIcon}/>
        <Resource name="media.youtube" list={NodeList} create={NodeCreate} show={NodeShow} options={{label: "Youtube"}}
                  icon={VideoIcon}/>
        <Resource name="media.image" list={NodeList} create={NodeCreate} show={NodeShow} options={{label: "Image"}}
                  icon={ImageIcon}/>
        <Resource name="feed.index" list={NodeList} create={NodeCreate} show={NodeShow} options={{label: "Feed Index"}}
                  icon={ListIcon}/>
        <Resource name="core.raw" list={NodeList} create={NodeCreate} show={NodeShow} options={{label: "Raw"}}/>
    </Admin>
);

export default App;