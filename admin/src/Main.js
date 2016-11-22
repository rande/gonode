// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

import React from 'react';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import AppBar from 'material-ui/AppBar';
import Drawer from 'material-ui/Drawer';
import { SMALL } from 'material-ui/utils/withWidth';
import { connect } from 'react-redux';
import { Router, Route, IndexRoute } from 'react-router';

import { toggleDrawer, hideDrawer } from './apps/guiApp';
import About from './components/About';
import Dashboard from './components/Dashboard';
import LoginForm from './components/LoginForm';
import MenuList from './components/MenuList';
import ListPanel from './components/node/ListPanel';
import ViewPanel from './components/node/ViewPanel';

const Container = props => (
    <div>{props.children}</div>
);

Container.propTypes = {
    children: React.PropTypes.any,
};

const Main = (props) => {
    let DrawerOpen = false;
    let marginLeft = 0;
    if (props.width === SMALL && props.DrawerOpen) {
        DrawerOpen = true;
    }

    if (props.width !== SMALL) {
        DrawerOpen = true;
        marginLeft = 250;
    }

    var mainScreen;

    if (props.roles.length == 0 ) {
        // load valid components
        mainScreen = <div>You don't have any roles associated to your account</div>
    } else if(props.roles.length == 1 && props.roles.includes('anonymous')) {
        mainScreen = <LoginForm />
    } else {
        mainScreen = <div>
            <AppBar
                title={props.Title}
                iconClassNameRight="muidocs-icon-navigation-expand-more"
                onLeftIconButtonTouchTap={props.toggleDrawer}
                showMenuIconButton={props.width === SMALL}
            />

            <Drawer open={DrawerOpen} docked width={marginLeft}>
                <AppBar
                    title={props.Title}
                    onLeftIconButtonTouchTap={props.toggleDrawer}
                    showMenuIconButton={props.width === SMALL}
                />

                <MenuList />
            </Drawer>

            <div className="foobar" style={{ marginLeft: `${marginLeft}px` }}>
                <Router history={props.history} >
                    <Route path="/" component={Container}>
                        <IndexRoute component={Dashboard} />
                        <Route path="nodes" component={ListPanel} />
                        <Route path="nodes/:uuid" component={ViewPanel} />
                        <Route path="about" component={About} />
                    </Route>
                </Router>
            </div>
        </div>;
    }

    return (<MuiThemeProvider muiTheme={props.Theme}>
        {mainScreen}
    </MuiThemeProvider>);
};

Main.propTypes = {
    Theme: React.PropTypes.object,
    Title: React.PropTypes.string,
    DrawerOpen: React.PropTypes.bool,
    toggleDrawer: React.PropTypes.func,
    history: React.PropTypes.object,
    width: React.PropTypes.number.isRequired,
};

const mapStateToProps = (state) => {
    return {
        ...state.guiApp,
        user: state.userApp.user,
        roles: state.userApp.roles,
        userState: state.userApp.state,
    };
};

const mapDispatchToProps = dispatch => ({
    toggleDrawer: () => {
        dispatch(toggleDrawer());
    },
});

export default connect(mapStateToProps, mapDispatchToProps)(Main);
