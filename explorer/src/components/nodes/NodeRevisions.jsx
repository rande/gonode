import React, { Component, PropTypes } from 'react';
import NodeRevisionsItem               from './NodeRevisionsItem.jsx';
import { fetchNodeRevisionsIfNeeded }  from '../../actions';


class NodeRevisions extends Component {
    fetchRevisions() {
        const { dispatch, uuid } = this.props;

        dispatch(fetchNodeRevisionsIfNeeded(uuid));
    }

    componentDidMount() {
        this.fetchRevisions();
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

NodeRevisions.propTypes = {
    uuid:      PropTypes.string.isRequired,
    revisions: PropTypes.array.isRequired,
    dispatch:  PropTypes.func.isRequired
};

export default NodeRevisions;
