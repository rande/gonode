import React from 'react';

import { List, Edit, Datagrid } from 'admin-on-rest';
import { ReferenceField, TextField, DateField, BooleanField } from 'admin-on-rest';
import { SimpleShowLayout, Show, ShowButton, EditButton } from 'admin-on-rest';
import { Create, TabbedForm, FormTab, SaveButton, Toolbar, DateInput, DisabledInput, LongTextInput, ReferenceInput, SelectInput, SimpleForm, TextInput, NumberInput, BooleanInput, RadioButtonGroupInput, CheckboxGroupInput } from 'admin-on-rest';

export const UserList = (props) => {
    return (<List {...props} perPage={32}>
        <Datagrid>
            <TextField source="data.username" label="Username"/>
            <TextField source="data.email" label="Email"/>
            <BooleanField source="data.locked" label="Locked"/>
            <BooleanField source="data.expired" label="Expired"/>

            <EditButton/>
            <ShowButton/>
        </Datagrid>
    </List>);
};

const UserTitle = ({ record }) => {
    return <span>Node {record ? `"${record.name}"` : ''}</span>;
};

export const UserShow = (props) => (
    <Show title={<UserTitle />} {...props}>
        <SimpleShowLayout>
            <TextField source="id" />
            <TextField source="name" />
            <TextField source="slug" />
            <TextField source="path" />
            <TextField source="status" />
            <TextField source="weight" />
            <TextField source="revision" />
            <TextField source="version" />
            <BooleanField source="enabled" />
            <BooleanField source="deleted" />
            <DateField label="Creation date" source="created_at" />
            <DateField label="Updated date" source="updated_at" />
        </SimpleShowLayout>
    </Show>
);

const UserNodeTitle = ({ record }) => {
    return <span>Add user</span>;
};

const UserCreateToolbar = props => <Toolbar {...props} >
    <SaveButton label="Save And Show" redirect="show" submitOnEnter={true} />
    <SaveButton label="Save And Add" redirect={false} submitOnEnter={false} raised={false} />
</Toolbar>;

export const UserCreate = (props) => (
    <Create title={<UserNodeTitle />} {...props}>
        <TabbedForm toolbar={<UserCreateToolbar />} redirect="show">
            <FormTab label="User">
                <TextInput label="Username" source="data.username" />
                <TextInput label="Firstname" source="data.firstname" />
                <TextInput label="Lastname" source="data.lastname" />
                <TextInput label="Email" source="data.email" />
                <BooleanInput label="Locked" source="data.locked" />
                <BooleanInput label="Enabled" source="data.enabled" />
                <BooleanInput label="Expired" source="data.expired" />
                <SelectInput label="Gender" source="data.gender" choices={[
                    { id: 'm', name: 'Male' },
                    { id: 'f', name: 'Female' },
                    { id: 'o', name: 'Other' },
                    { id: 'u', name: 'Unknown' },
                ]} />
                <TextInput label="Locale" source="data.locale" />
                <TextInput label="Timezone" source="data.timezone" />
                <DateInput label="Date of birth" source="data.dateofbirth" />
            </FormTab>
            <FormTab label="Node">
                <BooleanInput label="enabled" source="enabled"/>
                <NumberInput label="Weight" source="weight" />
                <TextInput label="Name" source="name"/>
                <TextInput label="Slug" source="slug"/>

                <SelectInput label="Status" source="status" default={0} choices={[
                    { id: 0, name: 'New' },
                    { id: 1, name: 'Draft' },
                    { id: 2, name: 'Completed' },
                    { id: 3, name: 'Validated' },
                ]} />
                <CheckboxGroupInput source="access" choices={[
                    { id: 'node:api:master', name: 'node:api:master' },
                    { id: 'node:api:read', name: 'node:prism:render' },
                    { id: 'node:prism:render', name: 'node:prism:render' },
                ]} />
            </FormTab>

            <FormTab label="Information">
                <DisabledInput label="Version" source="version"/>
                <DisabledInput label="Parent" source="parent_uuid"/>
                <DisabledInput label="Set" source="set_uuid"/>
                <DisabledInput label="Source" source="source"/>
                <DisabledInput label="Updated At" source="updated_at"/>
                <DisabledInput label="Created At" source="created_at"/>
            </FormTab>
        </TabbedForm>
    </Create>
);