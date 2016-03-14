import React, { Component, PropTypes } from 'react';
import { connect }                     from 'react-redux';
import classNames                      from 'classnames';
import ReactCSSTransitionGroup         from 'react-addons-css-transition-group';
import { FormattedMessage }            from 'react-intl';
import NodeRevisionsItem               from './NodeRevisionsItem.jsx';
import { nodeRevisionsSelector }       from '../../selectors/nodes-selector';
import { fetchNodeRevisionsIfNeeded }  from '../../actions';


class NodeRevisions extends Component {
    static displayName = 'NodeRevisions';

    static propTypes = {
        uuid:             PropTypes.string.isRequired,
        node:             PropTypes.object,
        isFetching:       PropTypes.bool.isRequired,
        revisionsByMonth: PropTypes.array.isRequired,
        nextPage:         PropTypes.number.isRequired,
        fetchMore:        PropTypes.func.isRequired
    };

    constructor(props) {
        super(props);

        this.handleMoreClick = this.handleMoreClick.bind(this);
    }

    handleMoreClick() {
        const { fetchMore, uuid, nextPage } = this.props;
        fetchMore(uuid, nextPage);
    }

    render() {
        const {uuid, node, isFetching, revisionsByMonth, nextPage } = this.props;

        const revNodes = [];
        revisionsByMonth.forEach(month => {
            revNodes.push(
                <span key={`revisions.month.${month.year}.${month.month}`} className="node_revisions_range-header">
                    <span className="node_revisions_range-header_label">
                        <FormattedMessage id={`month.${month.month}.long`} />&nbsp;
                        {month.year}
                    </span>
                </span>
            );

            month.items.forEach(revision => {
                revNodes.push(
                    <NodeRevisionsItem
                        key={`revision.${revision.revision}`}
                        isCurrent={!!(node && node.revision === revision.revision)}
                        uuid={uuid}
                        revision={revision}
                    />
                );
            });
        });

        return (
            <div className="node_revisions">
                <div className="node_revisions_wrapper">
                    <ReactCSSTransitionGroup
                        transitionName="node_revisions_item"
                        transitionEnterTimeout={400}
                        transitionLeaveTimeout={400}
                    >
                        {revNodes}
                    </ReactCSSTransitionGroup>
                </div>
                <span className={classNames('node_revisions_loader', {
                    'node_revisions_loader-loading': isFetching
                })} />
                {(nextPage > 0) && (
                    <span className="node_revisions_more" onClick={this.handleMoreClick}>
                        <i className="fa fa-angle-down"/>
                    </span>
                )}
            </div>
        );
    }
}

const mapDispatchToProps = dispatch => ({
    fetchMore: (uuid, page) => dispatch(fetchNodeRevisionsIfNeeded(uuid, page))
});


export default connect(
    nodeRevisionsSelector,
    mapDispatchToProps
)(NodeRevisions);
