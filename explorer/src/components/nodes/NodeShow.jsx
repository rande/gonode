import React, { PropTypes }                from 'react';
import { connect }                         from 'react-redux';
import { FormattedMessage, FormattedDate } from 'react-intl';
import { Link }                            from 'react-router';
import NodeInfo                            from './NodeInfo.jsx';
import NodeDeleteButton                    from './NodeDeleteButton.jsx';


const NodeShow = ({ nodeObject }) => {
    if (nodeObject.isFetching) {
        return <div className="node-main"/>;
    }

    const { node } = nodeObject;

    return (
        <div className="node-main">
            <header className="panel-header">
                <Link to={`/nodes`} className="panel-header_close">
                    <i className="fa fa-angle-left" />
                </Link>
                <h1 className="panel-title">{node.name}</h1>
                <Link to={`/nodes/${node.uuid}/edit`} className="button button-large">
                    <i className="fa fa-pencil" />
                    <FormattedMessage id="node.edit.link"/>
                </Link>
                <NodeDeleteButton uuid={node.uuid} size="large" />
            </header>
            <div className="panel-body">
                <NodeInfo node={node} />
            </div>
        </div>
    );
};

Node.displayName = 'NodeShow';

Node.propTypes = {
    nodeObject: PropTypes.object.isRequired
};


export default connect((state) => {
    const { nodes, nodesByUuid } = state;

    let nodeObject = {
        isFetching: true,
        node:       null
    };
    if (nodesByUuid[nodes.currentUuid]) {
        nodeObject = nodesByUuid[nodes.currentUuid];
    }

    return { nodeObject };
})(NodeShow);
