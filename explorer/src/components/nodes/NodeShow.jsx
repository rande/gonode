import React, { PropTypes }                from 'react';
import { connect }                         from 'react-redux';
import { FormattedMessage, FormattedDate } from 'react-intl';
import { Link }                            from 'react-router';
import NodeInfo                            from './NodeInfo.jsx';
import NodeDeleteButton                    from './NodeDeleteButton.jsx';
import { nodeSelector }                    from '../../selectors/nodes-selector';
import Breadcrumb                          from '../Breadcrumb.jsx';


const NodeShow = ({ node }) => {
    if (node.isFetching) {
        return <div className="node-main"/>;
    }

    return (
        <div className="node-main">
            <header className="panel-header">
                <Link to={`/nodes`} className="panel-header_close">
                    <i className="fa fa-close" />
                </Link>
                <Breadcrumb items={[
                    { label: node.node.name }
                ]} />
                <Link to={`/nodes/${node.node.uuid}/edit`} className="button button-large">
                    <i className="fa fa-pencil" />
                    <FormattedMessage id="node.edit.link"/>
                </Link>
                <NodeDeleteButton uuid={node.node.uuid} size="large" />
            </header>
            <div className="panel-body">
                <NodeInfo node={node.node} />
            </div>
        </div>
    );
};

Node.displayName = 'NodeShow';

Node.propTypes = {
    node: PropTypes.object.isRequired
};


export default connect(nodeSelector)(NodeShow);
