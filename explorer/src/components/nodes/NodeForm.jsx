import React, { Component, PropTypes } from 'react';
import { reduxForm }                   from 'redux-form';


class NodeForm extends Component {
    render() {
        const {
            fields: {
                name,
                type,
                slug
            },
            handleSubmit,
            resetForm,
            submitting
        } = this.props;

        return (
            <div>
                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label>Name</label>
                        <div>
                            <input type="text" placeholder="node name" {...name}/>
                        </div>
                    </div>
                    <div className="form-group">
                        <label>Type</label>
                        <div>
                            <input type="text" placeholder="node type" {...type}/>
                        </div>
                    </div>
                    <div className="form-group">
                        <label>Slug</label>
                        <div>
                            <input type="text" placeholder="node slug" {...slug}/>
                        </div>
                    </div>
                    <button className="button" onClick={handleSubmit}>Submit</button>
                </form>
            </div>
        );
    }
}


export default reduxForm({
    form:   'node',
    fields: ['name', 'type', 'slug']
})(NodeForm);
