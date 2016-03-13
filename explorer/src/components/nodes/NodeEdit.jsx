import React, { Component, PropTypes } from 'react';
import { connect }                     from 'react-redux';
import { Link }                        from 'react-router';
import { FormattedMessage }            from 'react-intl';
import NodeForm                        from './NodeForm.jsx';
import { updateNode }                  from '../../actions';
import { nodeSelector }                from '../../selectors/nodes-selector';

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
        const edited = assign({}, node.node, data);

        dispatchNodeUpdate(edited);
    }

    render() {
        const { node } = this.props;

        if (node.isFetching) {
            return null;
        }

        return (
            <div className="node-main">
                <div className="panel-header">
                    <Link to={`/nodes`} className="panel-header_close">
                        <i className="fa fa-close" />
                    </Link>
                    <h1 className="panel-title">
                        <Link to={`/nodes/${node.node.uuid}`}>
                            {node.node.name}
                        </Link>
                        &nbsp;&nbsp;|&nbsp;&nbsp;
                        <FormattedMessage id="node.edit.title" values={{ name: node.node.name }}/>
                    </h1>
                </div>
                <div className="panel-body">
                    <NodeForm onSubmit={this.handleSubmit} initialValues={node.node}/>
                </div>
            </div>
        );
    }
}

const mapDispatchToProps = dispatch => ({
    dispatchNodeUpdate: (node) => dispatch(updateNode(node))
});


export default connect(
    nodeSelector,
    mapDispatchToProps
)(NodeEdit);
