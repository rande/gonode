import React from 'react';
import { List, Edit, Create, Datagrid, ReferenceField, TextField, EditButton, DisabledInput, LongTextInput, ReferenceInput, SelectInput, SimpleForm, TextInput, SimpleShowLayout, Show, DateField } from 'admin-on-rest';

export const NodeList = (props) => {
    return (<List {...props} perPage={32}>
        <Datagrid>
            <TextField source="id"/>
            <TextField source="name"/>
            <TextField source="type"/>
        </Datagrid>
    </List>);
};

const NodeTitle = ({ record }) => {
    return <span>Node {record ? `"${record.name}"` : ''}</span>;
};

export const NodeShow = (props) => (
    <Show title={<NodeTitle />} {...props}>
        <SimpleShowLayout>
            <TextField source="id" />
            <TextField source="name" />
            <TextField source="slug" />
            <TextField source="path" />
            <TextField source="status" />
            <TextField source="weight" />
            <TextField source="revision" />
            <TextField source="version" />
            <TextField source="enabled" />
            <TextField source="deleted" />
            <DateField label="Creation date" source="created_at" />
            <DateField label="Updated date" source="updated_at" />
            <ReferenceField label="User" source="updated_by" reference="core.user">
                <TextField source="name" />
            </ReferenceField>
        </SimpleShowLayout>
    </Show>
);

const CreateNodeTitle = ({ record }) => {
    return <span>Add node</span>;
};

export const NodeCreate = (props) => (
    <Create title={<CreateNodeTitle />} {...props}>
        <SimpleForm>
            <TextInput source="name" />
        </SimpleForm>
    </Create>
);