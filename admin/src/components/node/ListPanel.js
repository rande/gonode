// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

import React from 'react';

import {connect} from 'react-redux';

import ListSearchBar from './ListSearchBar';
import ListGrid from './ListGrid';
import ListPagination from './ListPagination'
import {searchNodes, loadNode} from '../../apps/nodeApp';

let ListPanel = (props) => {
    let searching = props.state == 'SEARCHING_NODES';

    return <div>
        <ListSearchBar {...props.search} searchNodes={props.searchNodes} services={props.services}/>
        <ListGrid {...props.result} loadNode={props.loadNode}/>
        <ListPagination {...props.result} searchNodes={props.searchNodes} />
    </div>;
};

ListPanel.propTypes = {
    result: React.PropTypes.object,
    search: React.PropTypes.object,
    state:  React.PropTypes.string,
    searchNodes: React.PropTypes.func,
    services: React.PropTypes.array,
};

const mapStateToProps = state => ({
    result: state.nodeApp.result,
    search: state.nodeApp.search,
    state:  state.nodeApp.state,
    services: state.nodeApp.services,
});

const mapDispatchToProps = dispatch => ({
    searchNodes: (params, page, per_page) => {
        dispatch(searchNodes(params, page, per_page))
    },
    loadNode: (uuid) => {
        dispatch(loadNode(uuid))
    }
});

export default connect(mapStateToProps, mapDispatchToProps)(ListPanel);