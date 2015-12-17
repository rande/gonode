import React, { Component, PropTypes }    from 'react';
import { connect }                        from 'react-redux';
import { Link }                           from 'react-router';
import classNames                         from 'classnames';
import { FormattedMessage }               from 'react-intl';
import NodeForm                           from './NodeForm.jsx';
import { createNode, fetchNodesIfNeeded } from '../../actions';


class NodeCreate extends Component {
    handleSubmit(data) {
        const { dispatch } = this.props;
        dispatch(createNode(data));
    }

    render() {
        return (
            <div>
                <h1 className="panel-title">
                    <FormattedMessage id="node.create.title"/>
                </h1>
                <NodeForm onSubmit={this.handleSubmit.bind(this)}/>
            </div>
        );
    }
}

NodeCreate.propTypes = {
    dispatch: PropTypes.func.isRequired
};


export default connect(state => {
    return {};
})(NodeCreate);
