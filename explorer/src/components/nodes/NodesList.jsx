import React, { Component, PropTypes } from 'react';
import NodesListItem                   from './NodesListItem.jsx';


class NodesList extends Component {
    render() {
        const { nodes } = this.props;

        return (
            <div className="nodes-list">
                {nodes.map(node => (
                    <NodesListItem
                        key={node.uuid}
                        node={node}
                    />
                ))}
            </div>
        );
    }
}

NodesList.propTypes = {
    nodes: PropTypes.array.isRequired
};


export default NodesList;
