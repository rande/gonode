import {h, VNode} from "hyperapp";

import getValue from "get-value";
import setValue from "set-value";

interface FugueHTMLInputElement extends HTMLInputElement {
    definition: InputDefinition
}

interface Context {
    getState(): any
    setState(state: any): any
    name: string
    renderField?(definition: InputDefinition, meta: MetaRenderFieldOption): VNode
}

interface NodeHandler {
    (node: VNode): VNode;
}

export interface MetaRenderFieldOption {
    touched: boolean
    error: string
    warning: string
}

// public interface with optional parameters (DX)
export interface FieldProps {
    name?: string
    type: string | Function
    component?: string | Function
    label?: string
    placeholder?: string
    meta?: MetaRenderFieldOption
    value?: string | number
    validators?: Function[]
    onclick?: Function | null
    disabled?: boolean
    class?: string
    // format(value: any, name: string): any,
    // normalize(value, previousValue, allValues, previousAllValues): value
    // validate: Array<Validator>
}

// internal interface with default valid values.
interface InnerFieldProps {
    name: string
    type: string | Function
    component: string | Function
    label: string
    placeholder: string
    meta: MetaRenderFieldOption
    value: string | number
    validators: Function[]
    onclick: Function | null
    disabled: boolean
    class: string
    // format(value: any, name: string): any,
    // normalize(value, previousValue, allValues, previousAllValues): value
    // validate: Array<Validator>
}

export interface InputDefinition {
    name: string
    nodeName: string
    value: any
    type: string | null
    label: string
    placeholder: string
    children: VNode[]
    events: {
        oncreate: Function | null
        oninput: Function | null
        onclick: Function | null
        onchange: Function | null
    }
    validators: Function[]
    extra: Object
    class: string
}

export const writeValue = (state: any, key: string, value: any): void => {
    return setValue(state, key, value, {});
};

export const readValue = (state: any, key: string): any => {
    return getValue(state, key, null)
};

const visitor = (node: VNode | Function, callback: NodeHandler): VNode => {
    if (typeof node === "function") {
        node = (node() as VNode);
    }

    node = callback(node);

    if (!node.children) {
        return node;
    }

    node.children = node.children.map((cnode: VNode | string): VNode | string => {
        if (typeof cnode === "string") {
            return cnode;
        }

        return visitor(cnode, callback);
    });

    return node;
};

export const getAttributes = (definition: InputDefinition): any => {
    const attributes = definition.extra;

    const validAttributes = ['class', 'value', 'type', 'placeholder'];

    validAttributes.map((a) => {
        if (!definition[a]) {
            return;
        }

        attributes[a] = definition[a];
    });

    Object.keys(definition.events).map((name) => {
        attributes[name] = definition.events[name];
    });

    return attributes;
};

export const defaultRenderField = (definition: InputDefinition, meta: MetaRenderFieldOption): VNode => {

    // if (meta.touched && meta.error) {
    //     definition.class += ' is-invalid';
    // }

    const input = h(definition.nodeName, getAttributes(definition), definition.children);

    return <div>
        <label>{definition.label} {meta.touched ? "*" : "" }</label>

        {input}
        {meta.touched && (
            (meta.error && <div>{meta.error}</div>) || (meta.warning && <span>{meta.warning}</span>)
        )}
    </div>;
};

export const normalizeFieldProps = (attributes: any) : InnerFieldProps => {
    return {
        name: '',
        type: 'text',
        component: 'input',
        label: '',
        placeholder: '',
        meta: {
            touched: false,
            error: '',
            warning: '',
        },
        value: '',
        class: '',
        disabled: false,
        validators: [],
        ...attributes
    }
};

const getMeta = (state: any, name: string) : MetaRenderFieldOption => {
    let meta = readValue(state, 'meta.' + name);

    if (!meta) {
        meta = {
            touched: false,
            error: '',
            warning: ''
        }
    }

    return meta;
};

export const Field = (props: FieldProps, children: VNode[]): VNode => {
    return <field {...props}>{children}</field>
};


const inputHandler = (ctx: Context, name: string) => (evt: Event) => {
    const target = (evt.target as FugueHTMLInputElement);

    let value : any;

    let meta = readValue(ctx.getState(), 'meta.' + name);
    if (!meta) {
        meta = {
            touched: true,
            error: '',
            warning: ''
        }
    }

    if (target.type === 'checkbox') {
        value = target.checked ? target.value : '';
    } else {
        value = target.value;
    }

    meta.error = '';

    const definition = (target.definition as InputDefinition);

    if (definition) {
        definition.validators.map((v) => {
            const error = v(value);

            if (meta.error !== '') {
                return;
            }

            meta.error = error ? error.msg : ''; // first error win.
        });
    }

    let state = writeValue(ctx.getState(), name, value);

    state = writeValue(state, 'meta.' + name, meta);

    // console.log("Form > inputHandler", {state, value, evt});

    ctx.setState(state);
};

const createInputDefinition = (def: any) : InputDefinition => {

    return {
        nodeName: 'input',
        value: null,
        type: null,
        label: '',
        placeholder: '',
        children: [],
        events: {
            oncreate: null,
            oninput: null,
            onclick: null,
            onchange: null,
        },
        extra: {},
        class: '',
        validators: [],
        ...def
    }
};

const getInputDefinition = (node: VNode<FieldProps>, ctx: Context): InputDefinition | boolean => {
    const attributes = normalizeFieldProps(node.attributes);

    const extra = {
        disabled: attributes.disabled
    };

    switch (attributes.type) {
        case 'email':
        case 'password':
        case 'text':
            return createInputDefinition({
                name: attributes.name,
                label: attributes.label,
                class: attributes.class,
                value: readValue(ctx.getState(), attributes.name),
                type: attributes.type,
                placeholder: attributes.placeholder,
                events: {
                    oninput: inputHandler(ctx, attributes.name),
                },
                extra: extra,
                validators: attributes.validators,
            });

        case 'textarea':
            return createInputDefinition({
                name: attributes.name,
                label: attributes.label,
                class: attributes.class,
                nodeName: 'textarea',
                value: readValue(ctx.getState(), attributes.name),
                placeholder: attributes.placeholder,
                events: {
                    oninput: inputHandler(ctx, attributes.name),
                },
                children: [
                    readValue(ctx.getState(), attributes.name)
                ],
                extra: extra,
                validators: attributes.validators,
            });

        case 'select-multiple':
        case 'select':
            let values = readValue(ctx.getState(), attributes.name);

            if (typeof values === "number") {
                values = [values.toString()];
            }

            if (typeof values === "string") {
                values = [values];
            }

            if (attributes.type == 'select-multiple') {
                extra['multiple'] = 'multiple';
            }

            return createInputDefinition({
                name: attributes.name,
                label: attributes.label,
                class: attributes.class,
                nodeName: 'select',
                value: readValue(ctx.getState(), attributes.name),
                events: {
                    oninput: inputHandler(ctx, attributes.name),
                },
                children: (node.children as Array<VNode>).map((option: VNode) => {
                    if (!option.attributes) {
                        option.attributes = {};
                    }

                    if (!values.includes(option.attributes['value'])) {
                        return option;
                    }

                    option.attributes['selected'] = 'selected';

                    return option;
                }),
                extra: extra,
                validators: attributes.validators,
            });

        case 'checkbox':
        case 'radio':
            if (readValue(ctx.getState(), attributes.name) == attributes.value) {
                extra['checked'] = 'checked';
            }

            return createInputDefinition({
                name: attributes.name,
                label: attributes.label,
                value: attributes.value,
                class: attributes.class,
                type: attributes.type,
                events: {
                    onchange: inputHandler(ctx, attributes.name),
                },
                extra: extra,
                validators: attributes.validators,
            });

        case 'button':
            let submitHandler = (callback: any) => (e: Event) => {
                e.preventDefault();

                // console.log("event from form", e);

                callback(ctx.getState());
            };

            return createInputDefinition({
                name: attributes.name,
                value: attributes.label,
                class: attributes.class,
                type: attributes.type,
                events: {
                    onclick: submitHandler(attributes.onclick),
                },
                extra: extra,
                validators: attributes.validators,
            });

        default:
            return false

    }
};

export const Form = (props: Context, children: Array<VNode>) => {
    const handler = (node: VNode): VNode => {
        if (node.nodeName !== 'field') {
            // console.debug("Form > ignore node, not handled node", {nodeName: node.nodeName, node: node});

            return node;
        }

        const field = (node as VNode<FieldProps>);

        if (!field.attributes) {
            // console.debug("Form > ignore node, missing `attributes` property", {nodeName: field.nodeName, field: field});

            return field;
        }

        // if (field.attributes && !field.attributes['name']) {
        //     console.debug("Form > ignore valid node, missing `name` attribute", {nodeName: field.nodeName, type: field.attributes.type, field: field});
        //
        //     return field;
        // }

        // console.debug("Form > Configure input", {type: field.attributes.type, field: field});

        const innerInputDefinition = getInputDefinition(field, props);

        if (innerInputDefinition === false || innerInputDefinition === true) { // not a VNode
            return field;
        }

        innerInputDefinition.events['oncreate'] = (element: FugueHTMLInputElement) => {
            element.definition = innerInputDefinition;
        };

        switch (field.attributes.component) {
            // case 'input':
            //     return innerInput;
            default:
            //     if (typeof field.attributes.component === 'function') {
            //         return field.attributes.component(innerInput, getMeta(props.getState(), field.attributes['name']));
            //     }

                const renderField = props.renderField ? props.renderField : defaultRenderField;

                return renderField(innerInputDefinition, getMeta(props.getState(), innerInputDefinition.name));
        }
    };

    return <form novalidate="novalidate" autocomplete="off">
        {children.map((cnode: VNode) => visitor(cnode, handler))}
    </form>
};
