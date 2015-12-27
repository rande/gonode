import React, { Component, PropTypes } from 'react';
import { reduxForm }                   from 'redux-form';
import classNames                      from 'classnames';


class LoginForm extends Component {
    render() {
        const {
            isFetching,
            fields: {
                login,
                password
            },
            handleSubmit,
            resetForm,
            submitting
        } = this.props;

        const submitClasses = classNames(
            'button',
            { '_is-disabled': isFetching }
        );

        return (
            <div className="login_form">
                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <input type="text" placeholder="login" {...login}/>
                    </div>
                    <div className="form-group">
                        <input type="password" placeholder="password" {...password}/>
                    </div>
                    <button className={submitClasses} onClick={handleSubmit}>Login</button>
                </form>
            </div>
        );
    }
}

LoginForm.propTypes = {
    isFetching: PropTypes.bool.isRequired
};


export default reduxForm({
    form:   'node',
    fields: ['login', 'password']
})(LoginForm);
