import React, { Component } from 'react';


class Footer extends Component {
    shouldComponentUpdate() {
        return false;
    }

    render() {
        return (
            <footer className="footer">
                <span className="footer_item">
                    Copyright Raphaël Benitte 2015 &copy;
                </span>
            </footer>
        );
    }
}


export default Footer;
