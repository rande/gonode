import React, { PropTypes }                from 'react';
import { connect }                         from 'react-redux';
import { FormattedMessage, FormattedDate } from 'react-intl';
import NodeRevisions                       from './NodeRevisions.jsx';


const Node = ({ nodeObject, revisions, dispatch }) => {
    if (nodeObject.isFetching) {
        return <span>loading node…</span>;
    }

    const { node } = nodeObject;

    return (
        <div className="node-show">
            <NodeRevisions
                uuid={node.uuid}
                dispatch={dispatch}
                revisions={revisions.items}
            />
            <div className="node-main">
                <h1 className="panel-title">{node.name}</h1>
                <div className="panel-body">
                    <ul className="node-properties">
                        <li>
                            <span className="node-properties_item_key">
                                <FormattedMessage id="node.uuid"/>
                            </span>&nbsp;
                            <span className="node-properties_item_value">{node.uuid}</span>
                        </li>
                        <li>
                            <span className="node-properties_item_key">
                                <FormattedMessage id="node.type"/>
                            </span>&nbsp;
                            <span className="node-properties_item_value">{node.type}</span>
                        </li>
                        <li>
                            <span className="node-properties_item_key">
                                <FormattedMessage id="node.status"/>
                            </span>&nbsp;
                            <span className="node-properties_item_value">
                                <FormattedMessage id={`node.status.${node.status}`}/>
                            </span>
                        </li>
                        <li>
                            <span className="node-properties_item_key">
                                <FormattedMessage id="node.revision"/>
                            </span>&nbsp;
                            <span className="node-properties_item_value">{node.revision}</span>
                        </li>
                        <li>
                            <span className="node-properties_item_key">
                                <FormattedMessage id="node.weight"/>
                            </span>&nbsp;
                            <span className="node-properties_item_value">{node.weight}</span>
                        </li>
                        <li>
                            <span className="node-properties_item_key">
                                <FormattedMessage id="node.enabled"/>
                            </span>&nbsp;
                            <span className="node-properties_item_value">{node.enabled ? 'yes' : 'no'}</span>
                        </li>
                        <li>
                            <span className="node-properties_item_key">
                                <FormattedMessage id="node.deleted"/>
                            </span>&nbsp;
                            <span className="node-properties_item_value">{node.deleted ? 'yes' : 'no'}</span>
                        </li>
                        <li>
                            <span className="node-properties_item_key">
                                <FormattedMessage id="node.created_at"/>
                            </span>&nbsp;
                            <span className="node-properties_item_value">
                                <FormattedDate
                                    value={new Date(node.created_at)}
                                    day="numeric"
                                    month="long"
                                    year="numeric"
                                />
                            </span>
                        </li>
                        <li>
                            <span className="node-properties_item_key">
                                <FormattedMessage id="node.updated_at"/>
                            </span>&nbsp;
                            <span className="node-properties_item_value">
                                <FormattedDate
                                    value={new Date(node.updated_at)}
                                    day="numeric"
                                    month="long"
                                    year="numeric"
                                />
                            </span>
                        </li>
                    </ul>
                    <pre className="node-raw">
                        {JSON.stringify(node, null, '  ')}
                    </pre>
                </div>
            </div>
        </div>
    );
};

Node.displayName = 'Node';

Node.propTypes = {
    nodeObject: PropTypes.object.isRequired,
    revisions:  PropTypes.object.isRequired,
    dispatch:   PropTypes.func.isRequired
};


export default connect((state) => {
    const { nodes, nodesByUuid, nodesRevisionsByUuid } = state;

    let nodeObject = {
        isFetching: true,
        node:       null
    };
    if (nodesByUuid[nodes.currentUuid]) {
        nodeObject = nodesByUuid[nodes.currentUuid];
    }

    let revisions = {
        items:      [],
        isFetching: true
    };
    if (nodesRevisionsByUuid[nodes.currentUuid]) {
        revisions = nodesRevisionsByUuid[nodes.currentUuid];
    }

    return {
        nodeObject,
        revisions
    };
})(Node);
