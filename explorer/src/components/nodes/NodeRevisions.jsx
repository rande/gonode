import React, { Component, PropTypes } from 'react';
import { connect }                     from 'react-redux';
import ReactCSSTransitionGroup         from 'react-addons-css-transition-group';
import NodeRevisionsItem               from './NodeRevisionsItem.jsx';
import { fetchNodeRevisionsIfNeeded }  from '../../actions';


class NodeRevisions extends Component {
    static displayName = 'NodeRevisions';

    static propTypes = {
        uuid:           PropTypes.string.isRequired,
        node:           PropTypes.object,
        isFetching:     PropTypes.bool.isRequired,
        revisions:      PropTypes.array.isRequired,
        fetchRevisions: PropTypes.func.isRequired
    };

    componentDidMount() {
        const { fetchRevisions, uuid } = this.props;
        fetchRevisions(uuid);
    }

    render() {
        const { uuid, node, revisions } = this.props;

        return (
            <div className="node_revisions">
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
        );
    }
}

const mapStateToProps = ({ nodesByUuid, nodesRevisionsByUuid }, { uuid }) => {
    let node = null;
    if (nodesByUuid[uuid]) {
        node = nodesByUuid[uuid].node;
    }

    let revisions = {
        items:      [],
        isFetching: true
    };
    if (nodesRevisionsByUuid[uuid]) {
        revisions = nodesRevisionsByUuid[uuid];
    }

    return {
        isFetching: revisions.isFetching,
        revisions:  revisions.items,
        node
    };
};

const mapDispatchToProps = dispatch => ({
    fetchRevisions: (nodeUuid) => dispatch(fetchNodeRevisionsIfNeeded(nodeUuid))
});


export default connect(
    mapStateToProps,
    mapDispatchToProps
)(NodeRevisions);
