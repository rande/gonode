// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

import React from 'react';
import {List, ListItem} from 'material-ui/List';
import Avatar from 'material-ui/Avatar';
import FileFolder from 'material-ui/svg-icons/file/folder';
import Done from 'material-ui/svg-icons/action/done';

const ListGrid = (props) => {
    let items = props.elements.map((elm, pos) => {
        let secondaryText = `${elm.type}`;
        let rightIcon = null;

        if (elm.enabled) {
            rightIcon = <Done />;
        }

        return <ListItem
            key={pos}
            leftAvatar={<Avatar icon={<FileFolder />}/>}
            rightIcon={rightIcon}
            primaryText={elm.name}
            secondaryText={secondaryText}
            onTouchTap={() => props.loadNode(elm.uuid)}
        />
    });

    return <List>{items}</List>;
};

ListGrid.propTypes = {
    elements: React.PropTypes.array,
    page:     React.PropTypes.number,
    per_page: React.PropTypes.number,
    next:     React.PropTypes.number,
    previous: React.PropTypes.number,
    state:    React.PropTypes.string,
    loadNode: React.PropTypes.func,
};

export default ListGrid;