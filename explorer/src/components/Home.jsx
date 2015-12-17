import React, { Component, PropTypes } from 'react';
import { FormattedMessage }            from 'react-intl';


class Home extends Component {
    render() {
        return (
            <div className="home">
                <div className="page-header">
                    <h2 className="page-header_title">
                        <FormattedMessage id="welcome"/>
                    </h2>
                </div>
            </div>
        )
    }
}


export default Home;
