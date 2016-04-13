import React, { PropTypes } from 'react';
import classNames           from 'classnames';
import { FormattedMessage } from 'react-intl';
import { Link }             from 'react-router';


const NodeRevisionsItem = ({ uuid, revision, isCurrent }) => {
    const updatedAt = new Date(revision.updated_at);

    const day       = updatedAt.getDate();
    const hours     = updatedAt.getHours();
    const minutes   = updatedAt.getMinutes();

    return (
        <div className={classNames('node_revisions_item', { 'node_revisions_item-current': isCurrent })}>
            <Link
                to={`/nodes/${uuid}/revisions/${revision.revision}`}
                className="node_revisions_item_circle"
                activeClassName="node_revisions_item_circle-active"
            >
                <span className="node_revisions_item_day">
                    {`${day < 10 ?  '0' : ''}${day}`}
                </span>
                <span className="node_revisions_item_time">
                    {`${hours < 10 ?  '0' : ''}${hours}:${minutes < 10 ?  '0' : ''}${minutes}`}
                </span>
                {isCurrent && <span className="node_revisions_item_current"/>}
            </Link>
        </div>
    );
};

NodeRevisionsItem.displayName = 'NodeRevisionsItem';

NodeRevisionsItem.propTypes = {
    uuid:      PropTypes.string.isRequired,
    revision:  PropTypes.object.isRequired,
    isCurrent: PropTypes.bool.isRequired
};


export default NodeRevisionsItem;
