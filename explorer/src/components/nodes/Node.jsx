import React, { PropTypes } from 'react';
import { connect }          from 'react-redux';
import NodeRevisions        from './NodeRevisions.jsx';


const Node = ({ uuid, content }) => {
    return (
        <div className="node-show">
            <NodeRevisions uuid={uuid}/>
            {content}
        </div>
    );
};

Node.displayName = 'Node';

Node.propTypes = {
    uuid:    PropTypes.string.isRequired,
    content: PropTypes.element.isRequired
};


export default connect(({ nodes: { uuid } }) => ({ uuid }))(Node);
