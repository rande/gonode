import React, { Component, PropTypes } from 'react';
import { Link }                        from 'react-router';
import _                               from 'lodash';
import Breadcrumbs                     from '../Breadcrumbs.jsx';

class Header extends Component {
    render() {
        const { region, country, routes, dispatch } = this.props;

        let searchButton;
        if (_.find(routes, { path: 'search' })) {
            searchButton = <Link to="/" className="search_button _is-active"/>;
        } else {
            searchButton = <Link to="/search" className="search_button"/>;
        }

        return (
            <header className="header">
                <Link to="/" className="header_logo"/>
                {searchButton}
                <Breadcrumbs
                    region={region}
                    country={country}
                    dispatch={dispatch}
                />
            </header>
        );
    }
}

Header.propTypes = {
    region:   PropTypes.object,
    country:  PropTypes.object,
    routes:   PropTypes.array.isRequired,
    dispatch: PropTypes.func.isRequired
};


export default Header;
