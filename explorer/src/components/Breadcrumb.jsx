import React, { Component, PropTypes } from 'react';
import { Link }                        from 'react-router';


const Breadcrumb = ({ items }) => (
    <div className="breadcrumb">
        {items.map((item, i) => {
            if (item.path) {
                return (
                    <Link key={i} to={item.path} className="breadcrumb_item breadcrumb_item-link">
                        {item.label}
                    </Link>
                );
            } else {
                return (
                    <span key={i} className="breadcrumb_item">
                        {item.label}
                    </span>
                );
            }
        })}
    </div>
);

Breadcrumb.displayName = 'Breadcrumb';

Breadcrumb.propTypes = {
    items: PropTypes.arrayOf(PropTypes.shape({
        link:  PropTypes.string,
        label: PropTypes.any.isRequired
    }))
};


export default Breadcrumb;
