import React, { Component, PropTypes } from 'react';
import { connect }                     from 'react-redux';
import { FormattedMessage }            from 'react-intl';
import LoginForm                       from './LoginForm.jsx';
import { login }                       from '../../actions';


class Login extends Component {
    handleSubmit(data) {
        const { dispatch } = this.props;
        dispatch(login(data));
    }

    render() {
        const { isFetching, loginFailed } = this.props;

        let loginError = null;
        if (loginFailed) {
            loginError = (
                <div className="login_error">
                    <FormattedMessage id="login.failed"/>
                </div>
            );
        }

        return (
            <div className="login">
                <div className="login_wrapper">
                    <div className="login_header">
                        <h2 className="login_header_brand">
                            <strong>GO</strong>NODE
                        </h2>
                    </div>
                    {loginError}
                    <LoginForm
                        isFetching={isFetching}
                        onSubmit={this.handleSubmit.bind(this)}
                    />
                </div>
            </div>
        )
    }
}

Login.propTypes = {
    isFetching:  PropTypes.bool.isRequired,
    loginFailed: PropTypes.bool.isRequired,
    dispatch:    PropTypes.func.isRequired
};


export default connect(state => {
    const { security: {
        isFetching,
        failed
    } } = state;

    return {
        isFetching,
        loginFailed: failed
    };
})(Login);
