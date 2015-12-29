import React, { Component, PropTypes }   from 'react';
import { connect }                       from 'react-redux';
import { Link }                          from 'react-router';
import classNames                        from 'classnames';
import { FormattedMessage }              from 'react-intl';
import NodeForm                          from './NodeForm.jsx';
import { fetchNodeIfNeeded, updateNode } from '../../actions';

const assign = Object.assign || require('object.assign');


class NodeEdit extends Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    fetchNode() {
        const { dispatch, routeParams } = this.props;

        dispatch(fetchNodeIfNeeded(routeParams.node_uuid));
    }

    componentDidMount() {
        this.fetchNode();
    }

    handleSubmit(data) {
        const { dispatch, node } = this.props;
        const edited = assign({}, node, data);

        dispatch(updateNode(edited));
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

NodeEdit.propTypes = {
    node:       PropTypes.object,
    isFetching: PropTypes.bool.isRequired,
    dispatch:   PropTypes.func.isRequired
};


export default connect((state) => {
    const { nodes, nodesByUuid } = state;

    let node = {
        isFetching: true,
        node:       null
    };

    if (nodesByUuid[nodes.currentUuid]) {
        node = nodesByUuid[nodes.currentUuid];
    }

    return node;
})(NodeEdit);
