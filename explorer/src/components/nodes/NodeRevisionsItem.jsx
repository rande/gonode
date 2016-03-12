import React, { PropTypes } from 'react';


const NodeRevisionsItem = ({ revision, isCurrent }) => (
    <div className="node_revisions_item">
        <span className="node_revisions_item_circle">
            {revision.revision}
        </span>
        {isCurrent && <span className="node_revisions_item_current" />}
    </div>
);

NodeRevisionsItem.displayName = 'NodeRevisionsItem';

NodeRevisionsItem.propTypes = {
    uuid:      PropTypes.string.isRequired,
    isCurrent: PropTypes.bool.isRequired,
    revision:  PropTypes.object.isRequired
};


export default NodeRevisionsItem;
