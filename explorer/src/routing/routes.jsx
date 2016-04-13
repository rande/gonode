import React                 from 'react';
import { Route, IndexRoute } from 'react-router';
import App                   from '../containers/App.jsx';
import Dashboard             from '../containers/Dashboard.jsx';
import Home                  from '../components/Home.jsx';
import Login                 from '../components/security/Login.jsx';
import Logout                from '../components/security/Logout.jsx';
import Nodes                 from '../components/nodes/Nodes.jsx';
import Node                  from '../components/nodes/Node.jsx';
import NodeCreate            from '../components/nodes/NodeCreate.jsx';
import NodeShow              from '../components/nodes/NodeShow.jsx';
import NodeEdit              from '../components/nodes/NodeEdit.jsx';
import NodeRevision          from '../components/nodes/NodeRevision.jsx';


import {
    ensureAuthenticated,
    onEnterApp,
    onEnterNode,
    onEnterNodeRevision,
    onEnterLogout
} from './hooks';


export default function getRoutes(store) {
    return (
        <Route path="/" onEnter={onEnterApp(store)} component={App}>
            <Route path="login" components={{ content: Login }}/>
            <Route path="logout" onEnter={onEnterLogout(store)} components={{ content: Logout }}/>
            <Route components={{ content: Dashboard }} >
                <IndexRoute onEnter={ensureAuthenticated(store)} components={{ content: Home }}/>
                <Route path="nodes" onEnter={ensureAuthenticated(store)} components={{ content: Nodes }}>
                    <Route path="create" components={{ content: NodeCreate }}/>
                    <Route path=":node_uuid" components={{ content: Node }}>
                        <IndexRoute onEnter={onEnterNode(store)} components={{ content: NodeShow }}/>
                        <Route path="edit" onEnter={onEnterNode(store)} components={{ content: NodeEdit }}/>
                        <Route path="revisions/:revision_id" onEnter={onEnterNodeRevision(store)} components={{ content: NodeRevision }}/>
                    </Route>
                </Route>
            </Route>
        </Route>
    );
}
