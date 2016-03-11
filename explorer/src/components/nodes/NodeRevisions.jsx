import React, { Component, PropTypes } from 'react';
import { connect }                     from 'react-redux';
import NodeRevisionsItem               from './NodeRevisionsItem.jsx';
import { fetchNodeRevisionsIfNeeded }  from '../../actions';


class NodeRevisions extends Component {
    static displayName = 'NodeRevisions';

    static propTypes = {
        uuid:           PropTypes.string.isRequired,
        isFetching:     PropTypes.bool.isRequired,
        revisions:      PropTypes.array.isRequired,
        fetchRevisions: PropTypes.func.isRequired
    };

    componentDidMount() {
        const { fetchRevisions, uuid } = this.props;
        fetchRevisions(uuid);
    }

    render() {
        const { uuid, revisions } = this.props;

        return (
            <div className="node_revisions">
                {revisions.map(revision => (
                    <NodeRevisionsItem
                        key={`revision.${revision.revision}`}
                        uuid={uuid}
                        revision={revision}
                    />
                ))}
            </div>
        );
    }
}

const mapStateToProps = ({ nodes, nodesRevisionsByUuid }) => {
    let revisions = {
        items:      [],
        isFetching: true
    };
    if (nodesRevisionsByUuid[nodes.currentUuid]) {
        revisions = nodesRevisionsByUuid[nodes.currentUuid];
    }

    return {
        isFetching: revisions.isFetching,
        revisions:  revisions.items
    };
};

const mapDispatchToProps = dispatch => ({
    fetchRevisions: (nodeUuid) => dispatch(fetchNodeRevisionsIfNeeded(nodeUuid))
});


export default connect(
    mapStateToProps,
    mapDispatchToProps
)(NodeRevisions);
