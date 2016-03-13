import React, { PropTypes } from 'react';
import { Link }             from 'react-router';


const NodeRevisionsItem = ({ uuid, revision, isCurrent }) => (
    <div className="node_revisions_item">
        <Link
            to={`/nodes/${uuid}/revisions/${revision.revision}`}
            className="node_revisions_item_circle"
            activeClassName="node_revisions_item_circle-active"
        >
            {revision.revision}
            {isCurrent && <span className="node_revisions_item_current" />}
        </Link>
    </div>
);

NodeRevisionsItem.displayName = 'NodeRevisionsItem';

NodeRevisionsItem.propTypes = {
    uuid:      PropTypes.string.isRequired,
    revision:  PropTypes.object.isRequired,
    isCurrent: PropTypes.bool.isRequired
};


export default NodeRevisionsItem;
