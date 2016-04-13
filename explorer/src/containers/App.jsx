import React, { PropTypes } from 'react';
import { connect }          from 'react-redux';


const App = ({ content }) => content;

App.propTypes = {
    content:  PropTypes.element.isRequired
};


export default connect(() => ({}))(App);
