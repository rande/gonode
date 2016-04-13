import React, { PropTypes } from 'react';
import { connect }          from 'react-redux';
import Navigation           from '../components/Navigation.jsx';


const Dashboard = ({ content }) => (
    <div>
        <Navigation/>
        <div className="content">
            {content}
        </div>
    </div>
);

Dashboard.propTypes = {
    content:  PropTypes.element.isRequired
};


export default connect(() => ({}))(Dashboard);
