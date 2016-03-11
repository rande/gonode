import React, { PropTypes } from 'react';
import { connect }          from 'react-redux';
import { FormattedMessage } from 'react-intl';
import LoginForm            from './LoginForm.jsx';
import { login }            from '../../actions';


const Login = ({ isFetching, failed, onLogin }) => (
    <div className="login">
        <div className="login_wrapper">
            <div className="login_header">
                <h2 className="login_header_brand">
                    <strong>GO</strong>NODE
                </h2>
            </div>
            {failed && (
                <div className="login_error">
                    <FormattedMessage id="login.failed"/>
                </div>
            )}
            <LoginForm
                isFetching={isFetching}
                onSubmit={onLogin}
            />
        </div>
    </div>
);

Login.propTypes = {
    isFetching: PropTypes.bool.isRequired,
    failed:     PropTypes.bool.isRequired,
    onLogin:    PropTypes.func.isRequired
};

const mapStateToProps = ({ security: { isFetching, failed } }) => {
    return { isFetching, failed };
};

const mapDispatchToProps = dispatch => ({
    onLogin: (data) => dispatch(login(data))
});


export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Login);
