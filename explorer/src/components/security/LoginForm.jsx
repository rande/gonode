import React, { PropTypes } from 'react';
import { reduxForm }        from 'redux-form';
import classNames           from 'classnames';

const LoginForm = ({
    isFetching,
    fields: { login, password },
    handleSubmit
}) => (
    <div className="login_form">
        <form onSubmit={handleSubmit}>
            <div className="form-group">
                <input type="text" placeholder="login" {...login}/>
            </div>
            <div className="form-group">
                <input type="password" placeholder="password" {...password}/>
            </div>
            <button
                className={classNames('button', { '_is-disabled': isFetching })}
                onClick={handleSubmit}
            >
                Login
            </button>
        </form>
    </div>
);

LoginForm.propTypes = {
    isFetching:   PropTypes.bool.isRequired,
    fields:       PropTypes.object.isRequired,
    handleSubmit: PropTypes.func.isRequired
};

export default reduxForm({
    form:   'node',
    fields: ['login', 'password']
})(LoginForm);
