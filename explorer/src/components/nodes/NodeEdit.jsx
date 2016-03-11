import React, { Component, PropTypes }   from 'react';
import { connect }                       from 'react-redux';
import { Link }                          from 'react-router';
import { FormattedMessage }              from 'react-intl';
import NodeForm                          from './NodeForm.jsx';
import { fetchNodeIfNeeded, updateNode } from '../../actions';

const assign = Object.assign || require('object.assign');


class NodeEdit extends Component {
    static displayName = 'NodeEdit';

    static propTypes = {
        node:                      PropTypes.object,
        isFetching:                PropTypes.bool.isRequired,
        dispatchFetchNodeIfNeeded: PropTypes.func.isRequired,
        dispatchNodeUpdate:        PropTypes.func.isRequired,
        routeParams:               PropTypes.object.isRequired
    };

    constructor(props) {
        super(props);

        this.handleSubmit = this.handleSubmit.bind(this);
    }

    fetchNode() {
        const { dispatchFetchNodeIfNeeded, routeParams } = this.props;
        dispatchFetchNodeIfNeeded(routeParams.node_uuid);
    }

    componentDidMount() {
        this.fetchNode();
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
            <div>
                <h1 className="panel-title">
                    <FormattedMessage id="node.edit.title" values={{ name: node.name }}/>
                </h1>
                <NodeForm onSubmit={this.handleSubmit} initialValues={node}/>
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
    dispatchFetchNodeIfNeeded: (nodeUuid) => dispatch(fetchNodeIfNeeded(nodeUuid)),
    dispatchNodeUpdate:        (node)     => dispatch(updateNode(node))
});


export default connect(
    mapStateToProps,
    mapDispatchToProps
)(NodeEdit);
