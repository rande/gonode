import React, { PropTypes } from 'react';


const NodeRevisionsItem = ({ revision }) => (
    <div className="node_revisions_item">
        <span className="node_revisions_item_circle">
            {revision.revision}
        </span>
    </div>
);

NodeRevisionsItem.displayName = 'NodeRevisionsItem';

NodeRevisionsItem.propTypes = {
    uuid:     PropTypes.string.isRequired,
    revision: PropTypes.object.isRequired
};


export default NodeRevisionsItem;
