// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

import React from 'react';
import Dialog from 'material-ui/Dialog';
import FlatButton from 'material-ui/FlatButton';
import AppBar from 'material-ui/AppBar';

import { reduxForm, Field, submit } from 'redux-form'

import {TextField} from 'redux-form-material-ui'

import { connect } from 'react-redux';
import { push } from 'react-router-redux';

import {authenticateUser} from '../apps/userApp';

let LoginForm = (props) => {
    let disabled = props.state == 'AUTHENTICATING_USER';
    let label = "Connect";

    if (disabled) {
        label = "Authenticating ...";
    }

    const actions = [
        <FlatButton
            label={label}
            primary={true}
            disabled={disabled}
            onClick={props.onClick}
        />,
    ];

    let title = <AppBar
        title="Authentification"
        showMenuIconButton={false}
      />;

    return <Dialog title={title} actions={actions} modal={false} open={true}>
        <form onSubmit={props.onSubmit}>
            <Field
                name="login"
                hintText="Enter your login"
                floatingLabelText="Login"
                fullWidth={true}
                component={TextField}
                disabled={disabled}
            />
            <br />
            <Field
                name="password"
                hintText="Password Field"
                floatingLabelText="Password"
                type="password"
                fullWidth={true}
                component={TextField}
                disabled={disabled}
            />
        </form>
    </Dialog>
};

LoginForm.propTypes = {
    onSubmit: React.PropTypes.func,
    errors: React.PropTypes.object,
    values: React.PropTypes.object,
    state: React.PropTypes.string,
};

const mapStateToProps = state => ({
    initialValues: state.userApp.login,
    state: state.userApp.state
});

const mapDispatchToProps = dispatch => ({
    onTouchStart: (mirror) => {
        dispatch(push(`/mirror/${mirror.Id}`));
    },
    homepage: () => {
        dispatch(push('/'));
    },
    onSubmit: (values) => {
        dispatch(authenticateUser(values.login, values.password));
    },
    onClick: () => {
        dispatch((submit('login')))
    }
});

LoginForm = reduxForm({
  form: 'login'
})(LoginForm);

LoginForm = connect(mapStateToProps, mapDispatchToProps)(LoginForm);

export default LoginForm;



