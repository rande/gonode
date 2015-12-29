import React, { Component, PropTypes }         from 'react';
import { Link }                                from 'react-router';
import { FormattedMessage, FormattedRelative } from 'react-intl';
import classNames                              from 'classnames';
import { history }                             from '../../routing';


class NodesListItem extends Component {
    constructor(props) {
        super(props);
        this.handleClick = this.handleClick.bind(this);
    }

    handleClick() {
        const { node } = this.props;
        history.push(`/nodes/${node.uuid}`);
    }

    render() {
        const { node } = this.props;

        return (
            <div onClick={this.handleClick} className="nodes-list_item">
                <h3 className="nodes-list_item_title">{node.name}</h3>
                <div className="nodes-list_item_meta">
                    <div className="nodes-list_item_type">
                        <i className="fa fa-hashtag"/>
                        {node.type}
                    </div>
                    <div className="nodes-list_item_creation">
                        <i className="fa fa-calendar-o"/>
                        <FormattedMessage
                            id="node.created_ago"
                            values={{
                                createdAt: (
                                    <FormattedRelative
                                        value={new Date(node.created_at)}
                                        style="numeric"
                                    />
                                ),
                                updatedAt: (
                                    <FormattedRelative
                                        value={new Date(node.updated_at)}
                                        style="numeric"
                                    />
                                )
                            }}
                        />
                    </div>
                    <div className="nodes-list_item_actions">
                        <Link to={`/nodes/${node.uuid}/edit`} className="button" onClick={e => e.stopPropagation()}>
                            <FormattedMessage id="node.edit.link"/>
                        </Link>
                    </div>
                </div>
            </div>
        );
    }
}

NodesListItem.propTypes = {
    node: PropTypes.object.isRequired
};


export default NodesListItem;
