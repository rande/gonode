import React, { Component, PropTypes }         from 'react';
import { Link }                                from 'react-router';
import { FormattedMessage, FormattedRelative } from 'react-intl';
import classNames                              from 'classnames';


class NodesListItem extends Component {
    render() {
        const { node } = this.props;

        return (
            <Link to={`/nodes/${node.uuid}`} className="nodes-list_item">
                <h3 className="nodes-list_item_title">{node.name}</h3>
                <div className="nodes-list_item_meta">
                    <span className="nodes-list_item_type">
                        <i className="fa fa-hashtag"/>
                        {node.type}
                    </span>
                    <span className="nodes-list_item_creation">
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
                    </span>
                </div>
            </Link>
        );
    }
}

NodesListItem.propTypes = {
    node: PropTypes.object.isRequired
};


export default NodesListItem;
