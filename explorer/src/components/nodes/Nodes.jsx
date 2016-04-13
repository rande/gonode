import React, { Component, PropTypes } from 'react';
import { connect }                     from 'react-redux';
import { Link }                        from 'react-router';
import classNames                      from 'classnames';
import { FormattedMessage }            from 'react-intl';
import Pager                           from '../Pager.jsx';
import NodesList                       from './NodesList.jsx';
import { fetchNodesIfNeeded }          from '../../actions';
import { history }                     from '../../routing';
import { nodesSelector }               from '../../selectors/nodes-selector';


class Nodes extends Component {
    static displayName = 'Nodes';

    static propTypes = {
        nodes:        PropTypes.array.isRequired,
        itemsPerPage: PropTypes.number.isRequired,
        currentPage:  PropTypes.number.isRequired,
        dispatch:     PropTypes.func.isRequired
    };

    handlePagerChange(pagerData) {
        history.push(`/nodes?p=1&pp=${pagerData.perPage}`);
    }

    fetchNodes() {
        const { dispatch, location, itemsPerPage, currentPage } = this.props;
        const { query } = location;

        dispatch(fetchNodesIfNeeded({
            page:    query.p ? parseInt(query.p) : currentPage,
            perPage: query.pp ? parseInt(query.pp) : itemsPerPage
        }));
    }

    componentDidMount() {
        this.fetchNodes();
    }

    componentDidUpdate() {
        this.fetchNodes();
    }

    render() {
        const { nodes, itemsPerPage, currentPage, previousPage, nextPage, isFetching, content } = this.props;

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
                <div className="nodes-wrapper">
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
                        page={currentPage}
                        previousPage={previousPage}
                        nextPage={nextPage}
                        onChange={this.handlePagerChange.bind(this)}
                    />
                    <NodesList nodes={nodes}/>
                </div>
                <Link to="/nodes" className={overlayClasses}/>
                <div className={panelClasses}>
                    {content}
                </div>
            </div>
        );
    }
}


export default connect(nodesSelector)(Nodes);
