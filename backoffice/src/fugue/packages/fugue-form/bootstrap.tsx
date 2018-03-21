import {h, VNode} from "hyperapp";
import {getAttributes, InputDefinition, MetaRenderFieldOption, defaultRenderField} from "./index";

const renderFormGroup = (definition: InputDefinition, meta: MetaRenderFieldOption): VNode => {
    const attr = getAttributes(definition);

    attr.class = `${attr.class ? attr.class : ''} form-control`;

    if (meta.touched && meta.error) {
        attr.class += ' is-invalid';
    }

    const input = h(definition.nodeName, attr, definition.children);

    return <fieldset class={"form-group"}>
        <label>{definition.label} {meta.touched ? "*" : "" }</label>

        {input}
        {meta.touched && (
            (meta.error && <div class="invalid-feedback">{meta.error}</div>) || (meta.warning && <span>{meta.warning}</span>)
        )}
    </fieldset>;
};

const renderCheckGroup = (definition: InputDefinition, meta: MetaRenderFieldOption): VNode => {
    const attr = getAttributes(definition);

    attr.class = `${attr.class ? attr.class : ''} form-check-input`;

    if (meta.touched && meta.error) {
        attr.class += ' is-invalid';
    }

    const input = h(definition.nodeName, attr, definition.children);

    return <fieldset class={"form-check"}>
        {input}
        <label class="form-check-label">{definition.label} {meta.touched ? "*" : "" }</label>
        {meta.touched && (
            (meta.error && <div class="invalid-feedback">{meta.error}</div>) || (meta.warning && <span>{meta.warning}</span>)
        )}
    </fieldset>;
};

const renderButton = (definition: InputDefinition, meta: MetaRenderFieldOption): VNode => {

    const attr = getAttributes(definition);

    attr.class = `btn ${attr.class ? attr.class : ''}`;

    const input = h(definition.nodeName, attr, definition.children);

    return input;
}

export const renderField = (definition: InputDefinition, meta: MetaRenderFieldOption): VNode => {
    switch (`${definition.nodeName}${definition.type ? '-' + definition.type : ''}`) {
        case 'input-text':
        case 'input-password':
        case 'input-email':
        case 'select':
        case 'textarea':
            return renderFormGroup(definition, meta);

        case 'input-checkbox':
        case 'input-radio':
            return renderCheckGroup(definition, meta);

        case 'input-button':
            return renderButton(definition, meta);

        default:
            console.log('bootstrap: using default renderer', `${definition.nodeName}${definition.type ? '-' + definition.type : ''}`);
            return defaultRenderField(definition, meta);
    }
};