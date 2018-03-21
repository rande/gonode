import {VNode, h} from "hyperapp";

export interface UpdateStatePanel {
    ref: string
    state: any
}

export interface AppActions {
    dispatch(name: string, arg: Array<any>): AppState
    pushPanel(panel: VNode): AppState
    updateState(data: UpdateStatePanel): AppState
    popPanel(data: UpdateStatePanel): AppState
}

export interface AppState {
    title: string
    values: Map<string, any>
    panels: Map<string, any>
    panelStack: Array<string>
    connected: boolean,
    currentPanel: string
}

export interface Context {
    actions: any
    state: any
    appActions: any
    readonly panelRef: string
    updateState(state: any): void
    get(id: string): any
}

export interface PanelProps {
    ctx: Context
    actions: any
    state: any
}

export const currier = (callback: Function, ...innerArgs: Array<any>) => {
    return (...args: Array<any>) => {
        return callback(...args)(...innerArgs);
    }
};

export const applyCurry = (subject: Object, ...innerArgs: Array<any>) => {
    Object.keys(subject).map((key) => {
        subject[key] = currier(subject[key], ...innerArgs)
    });

    return subject;
};

export const Panel = (props: PanelProps, children: Array<VNode | Function>) => (state: AppState, actions: AppActions) : VNode => {
    // console.log('Render panel', {props, state});

    if (!props.ctx) {
        console.error("Missing the ctx parameter", {props});

        return <div />;
    }

    console.log("Create ctx", props.ctx, props.actions, props.state, state.values[props.ctx.panelRef]);

    let ctx = props.ctx;
    ctx.actions = applyCurry(props.actions || {}, props.ctx);
    ctx.state = state.values[ctx.panelRef] || props.state;
    ctx.updateState = (callback) => {
        let panelState = {};

        if (typeof callback === 'function') {
            panelState = callback(state.values[ctx.panelRef])
        } else {
            panelState = callback;
        }

        actions.updateState({ref: ctx.panelRef, state: panelState});
    };

    let finalChildren = children
        .filter((child) => typeof child === 'function')
        .map((child) => {
            return ({...args}) => (state: AppState, actions: AppActions): VNode => {
                const callback = (child as Function);

                return callback(ctx);
            }
        });

    if (finalChildren.length === 0) {
        return <div />;
    }

    return (finalChildren[0] as Function)(ctx);
};

export const Zone = (props: any, children: Array<VNode>) => {
    let finalProps =  {
        ...props,
        class: `${props.class ? props.class + ' ' : ''}zone-${props.name}`
    };

    return <zone {...finalProps}>
        { children }
    </zone>
};

export const renderZone = (name: string, children: Array<VNode>): VNode | void => {
    // remove non zone element
    let zone = children
        .filter((child: VNode) => {
            return child.nodeName === 'zone' && child.attributes && child.attributes['name'] === name;
        })
        .shift();

    if (!zone) {
        return;
    }

    return h("div", zone.attributes, zone.children);
};