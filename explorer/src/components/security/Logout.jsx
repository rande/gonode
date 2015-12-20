import React, { Component, PropTypes } from 'react';
import { connect }                     from 'react-redux';
import { Link }                        from 'react-router';
import { FormattedMessage }            from 'react-intl';


class Logout extends Component {
    render() {
        return (
            <div className="logout">
                <div className="logout_wrapper">
                    <div className="logout_header">
                        <h2 className="logout_header_brand">
                            <strong>GO</strong>NODE
                        </h2>
                    </div>
                    <FormattedMessage id="logout.message"/>
                    <Link to="/login" className="button">
                        <FormattedMessage id="login.link"/>
                    </Link>
                </div>
            </div>
        );
    }
}

export default connect(state => {
    return {};
})(Logout);