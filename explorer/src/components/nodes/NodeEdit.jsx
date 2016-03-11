import React, { Component, PropTypes } from 'react';
import { connect }                     from 'react-redux';
import { Link }                        from 'react-router';
import { FormattedMessage }            from 'react-intl';
import NodeForm                        from './NodeForm.jsx';
import { updateNode }                  from '../../actions';

const assign = Object.assign || require('object.assign');


class NodeEdit extends Component {
    static displayName = 'NodeEdit';

    static propTypes = {
        node:               PropTypes.object,
        isFetching:         PropTypes.bool.isRequired,
        dispatchNodeUpdate: PropTypes.func.isRequired
    };

    constructor(props) {
        super(props);

        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleSubmit(data) {
        const { dispatchNodeUpdate, node } = this.props;
        const edited = assign({}, node, data);

        dispatchNodeUpdate(edited);
    }

    render() {
        const { isFetching, node } = this.props;

        if (isFetching) {
            return null;
        }

        return (
            <div className="node-main">
                <h1 className="panel-title">
                    <FormattedMessage id="node.edit.title" values={{ name: node.name }}/>
                </h1>
                <div className="panel-body">
                    <NodeForm onSubmit={this.handleSubmit} initialValues={node}/>
                </div>
            </div>
        );
    }
}

const mapStateToProps = ({ nodes, nodesByUuid }) => {
    let node = {
        isFetching: true,
        node:       null
    };

    if (nodesByUuid[nodes.currentUuid]) {
        node = nodesByUuid[nodes.currentUuid];
    }

    return node;
};

const mapDispatchToProps = dispatch => ({
    dispatchNodeUpdate: (node) => dispatch(updateNode(node))
});


export default connect(
    mapStateToProps,
    mapDispatchToProps
)(NodeEdit);
