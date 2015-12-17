import React, { Component, PropTypes } from 'react';
import { Link }                        from 'react-router';


class Breadcrumbs extends Component {
    render() {
        return (
            <div className="breadcrumbs">
                <Link to="/" className="breadcrumbs_item">Home</Link>
            </div>
        );
    }
}

Breadcrumbs.propTypes = {
};


export default Breadcrumbs;
