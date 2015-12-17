import React                 from 'react';
import { Route, IndexRoute } from 'react-router';
import App                   from '../containers/App.jsx';
import Home                  from '../components/Home.jsx';
import Nodes                 from '../containers/Nodes.jsx';
import Node                  from '../containers/Node.jsx';
import NodeShow              from '../components/nodes/NodeShow.jsx';
import NodeCreate            from '../components/nodes/NodeCreate.jsx';


import {
    onEnterApp,
    onEnterNodes,
    onEnterNode
} from './hooks';


export default function getRoutes(store) {
    return (
        <Route path="/" onEnter={onEnterApp(store)} component={App}>
            <IndexRoute
                components={{ content: Home }}
            />
            <Route
                path="nodes"
                onEnter={onEnterNodes(store)}
                components={{ content: Nodes }}
            >
                <Route
                    path="create"
                    components={{ content: NodeCreate }}
                />
                <Route
                    path=":node_uuid"
                    onEnter={onEnterNode(store)}
                    components={{ content: Node }}
                />
            </Route>
        </Route>
    );
}
