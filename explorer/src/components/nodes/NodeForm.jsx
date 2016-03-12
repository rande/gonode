import React, { PropTypes } from 'react';
import { reduxForm }        from 'redux-form';


const NodeForm = ({
    fields: { name, type, slug },
    handleSubmit,
    resetForm,
    submitting
}) => (
    <div>
        <form onSubmit={handleSubmit} className="form">
            <div className="form_group">
                <label>Name</label>
                <div>
                    <input type="text" placeholder="node name" {...name}/>
                </div>
            </div>
            <div className="form_group">
                <label>Type</label>
                <div>
                    <input type="text" placeholder="node type" {...type}/>
                </div>
            </div>
            <div className="form_group">
                <label>Slug</label>
                <div>
                    <input type="text" placeholder="node slug" {...slug}/>
                </div>
            </div>
            <button className="button button-large" onClick={handleSubmit}>Submit</button>
        </form>
    </div>
);

NodeForm.displayName = 'NodeForm';

NodeForm.propTypes = {
    fields:       PropTypes.object.isRequired,
    handleSubmit: PropTypes.func.isRequired,
    resetForm:    PropTypes.func.isRequired,
    submitting:   PropTypes.bool.isRequired
};


export default reduxForm({
    form:   'node',
    fields: ['name', 'type', 'slug']
})(NodeForm);
