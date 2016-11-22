// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

import React from 'react';

import {Card, CardActions, CardHeader, CardMedia, CardTitle, CardText} from 'material-ui/Card';
import FlatButton from 'material-ui/FlatButton';
import Avatar from 'material-ui/Avatar';
import FileFolder from 'material-ui/svg-icons/file/folder';
import {connect} from 'react-redux';

import {loadNode} from '../../apps/nodeApp';
import ViewNodeDebug from './ViewNodeDebug';


let ViewPanel = (props) => {
    let {node} = props;

    return <Card>
        <CardHeader
              title={node.name}
              subtitle={node.uuid}
              avatar={<Avatar icon={<FileFolder />}/>}
            />

        <div>
            <ViewNodeDebug node={node} loadNode={props.loadNode} />
        </div>
    </Card>;
};

ViewPanel.propTypes = {
    node: React.PropTypes.object,
};

const mapStateToProps = state => ({
    node: state.nodeApp.node,
});

const mapDispatchToProps = dispatch => ({
    loadNode: (uuid) => {
        dispatch(loadNode(uuid))
    }
});

export default connect(mapStateToProps, mapDispatchToProps)(ViewPanel);