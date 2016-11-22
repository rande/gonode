// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

import React from 'react';

import {
    Table,
    TableBody,
    TableHeader,
    TableHeaderColumn,
    TableRow,
    TableRowColumn
} from 'material-ui/Table';

import {Tabs, Tab} from 'material-ui/Tabs';


const ViewNodeDebug = (props) => {

    const {node} = props;

    return <Tabs>
        <Tab label="Node">
            <Table selectable={false}>
                <TableBody displayRowCheckbox={false}>
                    <TableRow>
                        <TableRowColumn>UUID</TableRowColumn>
                        <TableRowColumn>{node.uuid}</TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Type</TableRowColumn>
                        <TableRowColumn>{node.type}</TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Name</TableRowColumn>
                        <TableRowColumn>{node.name}</TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Slug</TableRowColumn>
                        <TableRowColumn>{node.slug}</TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Path</TableRowColumn>
                        <TableRowColumn>{node.path}</TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Status</TableRowColumn>
                        <TableRowColumn>{node.status}</TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Weight</TableRowColumn>
                        <TableRowColumn>{node.weight}</TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Revision</TableRowColumn>
                        <TableRowColumn>{node.revision}</TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Version</TableRowColumn>
                        <TableRowColumn>{node.version}</TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Created At</TableRowColumn>
                        <TableRowColumn>{node.created_at}</TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Updated At</TableRowColumn>
                        <TableRowColumn>{node.updated_at}</TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Enabled</TableRowColumn>
                        <TableRowColumn>{node.enabled ? `yes` : `no`}</TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Deleted</TableRowColumn>
                        <TableRowColumn>{node.deleted ? `yes` : `no`}</TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Parents</TableRowColumn>
                        <TableRowColumn>{node.parents ? node.parents.map((uuid, key) => {
                            return <a key={key} href={`/node/${uuid}`} onClick={(e) => props.loadNode(uuid)}>{uuid}</a>
                        }) : null}</TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Updated By</TableRowColumn>
                        <TableRowColumn><a href={`/node/${node.updated_by}`} onClick={(e) => props.loadNode(node.updated_by)}>{node.updated_by}</a></TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Created By</TableRowColumn>
                        <TableRowColumn><a href={`/node/${node.created_by}`} onClick={(e) => props.loadNode(node.created_by)}>{node.created_by}</a></TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Parent Uuid</TableRowColumn>
                        <TableRowColumn><a href={`/node/${node.parent_uuid}`} onClick={(e) => props.loadNode(node.parent_uuid)}>{node.parent_uuid}</a></TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Set Uuid</TableRowColumn>
                        <TableRowColumn><a href={`/node/${node.set_uuid}`} onClick={(e) => props.loadNode(node.set_uuid)}>{node.set_uuid}</a></TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Source</TableRowColumn>
                        <TableRowColumn><a href={`/node/${node.source}`} onClick={(e) => props.loadNode(node.source)}>{node.source}</a></TableRowColumn>
                    </TableRow>
                    <TableRow>
                        <TableRowColumn>Access</TableRowColumn>
                        <TableRowColumn>{JSON.stringify(node.access, null, 2)}</TableRowColumn>
                    </TableRow>
                </TableBody>
            </Table>
        </Tab>
        <Tab label="Data">
            <pre>{JSON.stringify(node.data, null, 2)}</pre>
        </Tab>

        <Tab label="Meta">
            <pre>{JSON.stringify(node.meta, null, 2)}</pre>
        </Tab>
    </Tabs>;
};

ViewNodeDebug.propTypes = {
    node: React.PropTypes.object,
    loadNode: React.PropTypes.func,
};

export default ViewNodeDebug;