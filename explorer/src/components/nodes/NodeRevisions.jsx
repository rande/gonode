import React, { Component, PropTypes } from 'react';
import { connect }                     from 'react-redux';
import ReactCSSTransitionGroup         from 'react-addons-css-transition-group';
import { FormattedMessage }            from 'react-intl';
import NodeRevisionsItem               from './NodeRevisionsItem.jsx';
import { nodeRevisionsSelector }       from '../../selectors/nodes-selector';


const NodeRevisions = ({ uuid, node, isFetching, revisions }) => (
    <div className="node_revisions">
        <div className="node_revisions_wrapper">
            <ReactCSSTransitionGroup
                transitionName="node_revisions_item"
                transitionEnterTimeout={400}
                transitionLeaveTimeout={400}
            >
                {revisions.map(revision => (
                    <NodeRevisionsItem
                        key={`revision.${revision.revision}`}
                        isCurrent={!!(node && node.revision === revision.revision)}
                        uuid={uuid}
                        revision={revision}
                    />
                ))}
            </ReactCSSTransitionGroup>
        </div>
        <span className="node_revisions_more">
            <i className="fa fa-plus" />
        </span>
    </div>
);

NodeRevisions.displayName = 'NodeRevisions';

NodeRevisions.propTypes = {
    uuid:       PropTypes.string.isRequired,
    node:       PropTypes.object,
    isFetching: PropTypes.bool.isRequired,
    revisions:  PropTypes.array.isRequired
};


export default connect(nodeRevisionsSelector)(NodeRevisions);
