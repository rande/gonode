import React, { Component, PropTypes } from 'react';


class NodeRevisionsItem extends Component {
    render() {
        const { revision } = this.props;

        return (
            <div className="node_revisions_item">
                <span className="node_revisions_item_circle">
                    {revision.revision}
                </span>
            </div>
        );
    }
}

NodeRevisionsItem.propTypes = {
    uuid:     PropTypes.string.isRequired,
    revision: PropTypes.object.isRequired
};

export default NodeRevisionsItem;
