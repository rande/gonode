import React                from 'react';
import { Link }             from 'react-router';
import { FormattedMessage } from 'react-intl';


const Navigation = () => (
    <div className="navigation">
        <Link to="/" className="navigation_brand">
            <strong>GO</strong>
            NODE
        </Link>
        <Link to="/" className="navigation_item">
            <i className="fa fa-home"/>
            <FormattedMessage id="nav.home"/>
        </Link>
        <Link to="/nodes" className="navigation_item" activeClassName="navigation_item-active">
            <i className="fa fa-folder-o"/>
            <FormattedMessage id="nav.nodes"/>
        </Link>
        <Link to="/logout" className="navigation_item" activeClassName="navigation_item-active">
            <i className="fa fa-sign-out"/>
            <FormattedMessage id="logout.link"/>
        </Link>
    </div>
);


export default Navigation;
