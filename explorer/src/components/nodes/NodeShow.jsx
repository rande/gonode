import React, { Component, PropTypes } from 'react';
import { Link }                        from 'react-router';
import classNames                      from 'classnames';


class NodeShow extends Component {
    render() {
        const { node } = this.props;

        return (
            <div>
                {node.name}
            </div>
        );
    }
}

NodeShow.propTypes = {
    node: PropTypes.object.isRequired
};


export default NodeShow;
