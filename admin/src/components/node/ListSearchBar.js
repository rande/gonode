// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

import React from 'react';
import MenuItem from 'material-ui/MenuItem';
import { reduxForm, Field, submit } from 'redux-form'
import {TextField, SelectField, Checkbox} from 'redux-form-material-ui'
import { connect } from 'react-redux';


let ListSearchBar = (props) => {
    return (
        <div>
            <form onSubmit={props.onSubmit}>
                <Field name="type" floatingLabelText="Type" component={SelectField} >
                    {props.services.map((elm, pos) => {
                        return <MenuItem key={pos} value={elm.code} primaryText={elm.name} />
                    })}
                </Field>
                <Field name="name" floatingLabelText="Name" component={TextField} />
                <Field name="enabled" label="Enabled" component={Checkbox} />
            </form>
        </div>
    );
};

ListSearchBar.propTypes = {
    services: React.PropTypes.array,
    search: React.PropTypes.object,
};

const mapStateToProps = state => ({
    initialValues: {
        ...state.nodeApp.search.params,
        page: state.nodeApp.search.page,
        per_page: state.nodeApp.search.per_page,
    },

});

const mapDispatchToProps = dispatch => ({
    onSubmit: (values) => {
        dispatch(authenticateUser(values.login, values.password));
    },
    onClick: () => {
        dispatch((submit('node.search')))
    }
});

ListSearchBar = reduxForm({
  form: 'node.search'
})(ListSearchBar);

ListSearchBar = connect(mapStateToProps, mapDispatchToProps)(ListSearchBar);

export default ListSearchBar;