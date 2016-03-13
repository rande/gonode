import React, { Component, PropTypes } from 'react';
import { connect }                     from 'react-redux';
import { Link }                        from 'react-router';
import { nodeRevisionSelector }        from '../../selectors/nodes-selector';
import NodeInfo                        from './NodeInfo.jsx';
import Breadcrumb                      from '../Breadcrumb.jsx';


class NodeRevision extends Component {
    static displayName = 'NodeRevision';

    static propTypes = {
        nodeUuid:   PropTypes.string.isRequired,
        revisionId: PropTypes.number.isRequired,
        revision:   PropTypes.object,
        isFetching: PropTypes.bool.isRequired
    };

    render() {
        const { nodeUuid, revisionId, revision, isFetching } = this.props;

        if (!revision) {
            return null;
        }

        return (
            <div className="node-main">
                <header className="panel-header">
                    <Link to={`/nodes`} className="panel-header_close">
                        <i className="fa fa-close" />
                    </Link>
                    <Breadcrumb items={[
                        { label: revision.name, path:  `/nodes/${nodeUuid}` },
                        { label: `revision ${revisionId}` }
                    ]} />
                </header>
                <div className="panel-body">
                    <NodeInfo node={revision} />
                </div>
            </div>
        );
    }
}


export default connect(nodeRevisionSelector)(NodeRevision);
