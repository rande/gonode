import lodash                          from 'lodash';
import React, { Component, PropTypes } from 'react';
import { connect }                     from 'react-redux';
import Header                          from '../components/layout/Header.jsx';
import Navigation                      from '../components/Navigation.jsx';
import Footer                          from '../components/layout/Footer.jsx';


class Dashboard extends Component {
    render() {
        const {
            routes,
            content
        } = this.props;

        return (
            <div>
                <Navigation/>
                <div className="content">
                    {content}
                </div>
            </div>
        );
    }
}

Dashboard.propTypes = {
    routes:   PropTypes.array.isRequired,
    dispatch: PropTypes.func.isRequired
};


export default connect(state => {
    return {};
})(Dashboard);
