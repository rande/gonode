import React, { PropTypes } from 'react';
import { connect }          from 'react-redux';
import { FormattedMessage } from 'react-intl';
import NodeForm             from './NodeForm.jsx';
import { createNode }       from '../../actions';


const NodeCreate = ({ onCreateNode }) => (
    <div>
        <h1 className="panel-title">
            <FormattedMessage id="node.create.title"/>
        </h1>
        <NodeForm onSubmit={onCreateNode}/>
    </div>
);

NodeCreate.displayName = 'NodeCreate';

NodeCreate.propTypes = {
    onCreateNode: PropTypes.func.isRequired
};


export default connect(state => ({}), dispatch => {
    return {
        onCreateNode: (data) => dispatch(createNode(data))
    };
})(NodeCreate);
