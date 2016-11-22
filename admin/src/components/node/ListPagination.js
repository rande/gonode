// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

import React from 'react';

import RaisedButton from 'material-ui/RaisedButton';

const style = {
  margin: 6,
};

const ListPagination = (props) => {
    return <div>
      <RaisedButton label="Previous" style={style} disabled={props.previous < 1} onTouchTap={() => props.searchNodes(null, props.previous, null)}/>
      <RaisedButton label="Home" primary={true} style={style} onTouchTap={() => props.searchNodes(null, 1, null)} />
      <RaisedButton label="Next" secondary={true} style={style} disabled={props.next < 1} onTouchTap={() => props.searchNodes(null, props.next, null)}/>
    </div>
};

ListPagination.propTypes = {
    next:     React.PropTypes.number,
    previous: React.PropTypes.number,
    state:    React.PropTypes.string,
    searchNodes: React.PropTypes.func,
};

export default ListPagination;