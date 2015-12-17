import './styles/main.styl';
import 'font-awesome/less/font-awesome.less';

import 'intl';
import 'babel-polyfill';
import './modernizr';
import './match-media';
import React, { Component, PropTypes } from 'react';
import ReactDOM                        from 'react-dom';
import { Provider }                    from 'react-redux';
import { Router }                      from 'react-router';
import { syncReduxAndRouter }          from 'redux-simple-router';
import { addLocaleData, IntlProvider } from 'react-intl';
import { history, getRoutes }          from './routing';
import configureStore                  from './store/configure-store';
import { en, fr }                      from './i18n';

const store = configureStore();

syncReduxAndRouter(history, store);

ReactDOM.render(
    <Provider store={store}>
        <IntlProvider key="intl" locale="en" messages={en}>
            <Router history={history}>
                {getRoutes(store)}
            </Router>
        </IntlProvider>
    </Provider>,
    document.getElementById('app')
);
