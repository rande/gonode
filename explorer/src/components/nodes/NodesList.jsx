import React, { PropTypes } from 'react';
import NodesListItem        from './NodesListItem.jsx';


const NodesList = ({ nodes }) => (
    <div className="nodes-list">
        {nodes.map(node => (
            <NodesListItem
                key={node.uuid}
                node={node}
            />
        ))}
    </div>
);

NodesList.displayName = 'NodesList';

NodesList.propTypes = {
    nodes: PropTypes.array.isRequired
};


export default NodesList;
