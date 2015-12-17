import React, { Component, PropTypes } from 'react';
import { connect }                     from 'react-redux';
import { Link }                        from 'react-router';
import classNames                      from 'classnames';
import { FormattedMessage }            from 'react-intl';
import Pager                           from '../components/Pager.jsx';
import NodesList                       from '../components/nodes/NodesList.jsx';
import { setNodesPagerOptions }        from '../actions';


class Nodes extends Component {
    handlePagerChange(pagerData) {
        const { dispatch } = this.props;
        const {
            perPage
        } = pagerData;

        dispatch(setNodesPagerOptions({
            itemsPerPage: perPage
        }));
    }

    render() {
        const { nodes, itemsPerPage, isFetching, content } = this.props;

        const panelClasses = classNames(
            'second-panel',
            { '_is-opened': content }
        );

        const overlayClasses = classNames(
            'content-overlay',
            { '_is-opened': content }
        );

        return (
            <div>
                <div className="page-header">
                    <h2 className="page-header_title">
                        <FormattedMessage id="nodes.title"/>
                    </h2>
                    <Link to="/nodes/create" className="page-header_button">
                        <i className="fa fa-plus"/>&nbsp;
                        <FormattedMessage id="node.create.button"/>
                    </Link>
                    {isFetching && <span className="loader"/>}
                </div>
                <Pager
                    perPageOptions={[5, 10, 16, 32]}
                    perPage={itemsPerPage}
                    onChange={this.handlePagerChange.bind(this)}
                />
                <NodesList nodes={nodes}/>
                <Link to="/nodes" className={overlayClasses}/>
                <div className={panelClasses}>
                    {content}
                </div>
            </div>
        );
    }
}

Nodes.propTypes = {
    nodes:        PropTypes.array.isRequired,
    itemsPerPage: PropTypes.number.isRequired,
    dispatch:     PropTypes.func.isRequired
};

export default connect((state) => {
    const { nodes: {
        items,
        itemsPerPage,
        isFetching
    } } = state;

    return {
        nodes: items,
        itemsPerPage,
        isFetching
    };
})(Nodes);
