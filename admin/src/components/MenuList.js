// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

import React from 'react';
import Avatar from 'material-ui/Avatar';
import {List, ListItem} from 'material-ui/List';
import Dashboard from 'material-ui/svg-icons/action/dashboard';
import Info from 'material-ui/svg-icons/action/info';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';

import { hideDrawer } from '../apps/guiApp';
import { logoutUser } from '../apps/userApp';

const MenuList = (props) => {
    const items = [

        <ListItem
            key="logout"
            primaryText="Logout"
            leftIcon={<Dashboard viewBox="0 5 24 24" style={{width: 40, height: 40}}/>}
            onTouchTap={() => {
                props.logout();
            }}
        />,

        <ListItem
            key="dashboard"
            primaryText="Dashboard"
            leftIcon={<Dashboard viewBox="0 5 24 24" style={{width: 40, height: 40}}/>}
            onTouchTap={() => {
                props.homepage();
            }}
        />,

        <ListItem
            key="nodes"
            primaryText="Nodes"
            leftIcon={<Dashboard viewBox="0 5 24 24" style={{width: 40, height: 40}}/>}
            onTouchTap={() => {
                props.nodes();
            }}
        />,

        <ListItem
            key="about"
            primaryText="About"
            leftIcon={<Info viewBox="0 4 24 24" style={{width: 30, height: 30}}/>}
            onTouchTap={() => {
                props.about();
            }}
        />,
    ];

    return (<List>{items}</List>);
};

MenuList.propTypes = {
    onTouchStart: React.PropTypes.func,
    homepage: React.PropTypes.func,
    about: React.PropTypes.func,
    logout: React.PropTypes.func,
};

const mapStateToProps = state => ({

});

const mapDispatchToProps = dispatch => ({
    homepage: () => {
        dispatch(push('/'));
        dispatch(hideDrawer());
    },
    nodes: () => {
        dispatch(push('/nodes'));
        dispatch(hideDrawer());
    },
    about: () => {
        dispatch(push('/about'));
        dispatch(hideDrawer());
    },
    onTouchStart: (mirror) => {
        dispatch(push(`/mirror/${mirror.Id}`));
        dispatch(hideDrawer());
    },
    logout: () => {
        dispatch(logoutUser());
    }
});

export default connect(mapStateToProps, mapDispatchToProps)(MenuList);
