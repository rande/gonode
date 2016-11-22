// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

import React from 'react';
import { render } from 'react-dom';
import { createStore, combineReducers, applyMiddleware, compose } from 'redux';
import { Provider } from 'react-redux';
import { syncHistoryWithStore, routerReducer, routerMiddleware } from 'react-router-redux';
import thunk from 'redux-thunk';

import injectTapEventPlugin from 'react-tap-event-plugin';

import { hashHistory } from 'react-router';

import Main from './Main'; // Our custom react component
import { defineUser, userApp } from './apps/userApp';
import { nodeApp } from './apps/nodeApp';
import { guiApp, resizeApp } from './apps/guiApp';
import { searchNodes } from './apps/nodeApp';
import { reducer as formReducer } from 'redux-form'

// Needed for onTouchTap
// http://stackoverflow.com/a/34015469/988941
injectTapEventPlugin();

const middleware = routerMiddleware(hashHistory);
const reducers = combineReducers({
    userApp,
    guiApp,
    nodeApp,
    routing: routerReducer,
    form: formReducer
});

const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

let searchDispatch = store => next => action => {

    let result = next(action);

    let state = store.getState();

    if (action.type === 'redux-form/CHANGE' && action.meta.form === 'node.search') {
        const values = state.form.node.search.values;

        let params = {
            type: values.type,
            enabled: values.enabled,
            name: values.name,
        };

        store.dispatch(searchNodes(params));
    }

    return result;
};

const store = createStore(reducers, composeEnhancers(applyMiddleware(middleware, thunk, searchDispatch)));

if (window) {
    let deferTimer;
    window.addEventListener('resize', () => {
        clearTimeout(deferTimer);
        deferTimer = setTimeout(() => {
            store.dispatch(resizeApp(window.innerWidth));
        }, 200);
    });

    // init the state
    store.dispatch(resizeApp(window.innerWidth));
}

const history = syncHistoryWithStore(hashHistory, store);

render(<Provider store={store}><Main history={history} /></Provider>, document.getElementById('app'));
