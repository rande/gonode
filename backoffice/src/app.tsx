import {ActionsType, app, View} from "hyperapp";

import {Container, ContainerBuilder, Definition} from "fugue-ioc/di";

import {AppState, AppActions, Context, UpdateStatePanel} from "fugue-app";
import ActionDispatcher from "fugue-app/dispatcher";

import {createLoginPanel} from './view/Login';
import {createAboutPanel} from './view/About';

import {ApiClient} from "./core/api";

export const kuker = (data: any) => {
    try {
        window.postMessage({
            ...data,
            kuker: true,
            time: (new Date()).getTime(),
        }, '*');
    } catch(e) {
        console.log(e);
    }
};

export const deepCopy = (state: AppState) : AppState => {
    let panels = state.panels;

    delete state.panels;

    state = JSON.parse(JSON.stringify(state));

    state.panels = panels;

    return state;
};

let createContext: () => Context;

const actions: ActionsType<AppState, AppActions> = {
    dispatch: () => {},
    pushPanel: (panelCreator: Function) => (state: AppState) => {
        let ctx = createContext();

        state = deepCopy(state);

        state.panels[ctx.panelRef] = panelCreator(ctx);
        state.values[ctx.panelRef] = ctx.state; // this state should be empty as not yet called by the panel closure.

        state.currentPanel = ctx.panelRef;

        state.panelStack.push(ctx.panelRef);

        kuker({
            type: 'pushPanel',
            state: state,
            icon: 'fa-plus-square',
            color: '#002099'
        });

        return state;
    },
    popPanel: () => (state: AppState) => {
        state = deepCopy(state);

        let ref = (state.panelStack.pop() as string);

        delete state.panels[ref];
        delete state.values[ref];

        state.currentPanel = state.panelStack[state.panelStack.length - 1];

        kuker({
            type: 'popPanel',
            state: state,
            icon: 'fa-minus-square',
            color: '#da7c00'
        });

        return state;
    },
    updateState: (data: UpdateStatePanel) => (state: AppState) => {
        if (data.ref == undefined) {
            console.error("updateState global, panelRef is undefined");

            return state;
        }

        if (data.state == undefined) {
            console.error("updateState global, state is undefined");

            return state;
        }

        state = deepCopy(state);

        state.values[data.ref] = data.state;

        kuker({
            type: 'updateState',
            state: state,
            icon: 'fa-sync-alt',
            color: '#bada55'
        });

        return state;
    }
};

const state: AppState = {
    title: "App Administration",
    values: new Map<string, any>(), // store the state for the each panel
    panels: new Map<string, any>(), // store the panel (views)
    currentPanel: '',
    connected: false,
    panelStack: [],
};

const view: View<AppState, AppActions> = (state: AppState, actions: AppActions) => {
    console.log('Refresh UI');

    return state.panels[state.currentPanel];
};

const defaultConfiguration = {
    api: {
        url: "http://localhost:2508/api"
    }
}

function loadApp(target: Element | null, configuration = defaultConfiguration) {

    // get the configuration somewhere

    const builder = new ContainerBuilder();
    const container = new Container();

    // core services
    builder.set("actions.dispatcher", new Definition(ActionDispatcher, []))
    
    // view constructors
    builder.set("panels.login", new Definition(createLoginPanel, []));
    builder.set("panels.about", new Definition(createAboutPanel, []));

    // action reaction
    builder.set("actions.user.connection", new Definition(configuration.api.url, [], ['dispatcher.action']))

    // external services
    builder.set("gonode.api", new Definition(ApiClient, []));
    

    builder.build(container);

    let appActions = app<AppState, AppActions>(
        state,
        actions,
        view,
        target
    );

    createContext = (): Context => {
        return {
            actions: false,
            state: false,
            appActions: appActions,
            panelRef: Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15),
            updateState: (state: any) => {},
            get: (name: string): any => {
                let [service, err] = container.get(name);

                if (err) {
                    console.error("Unable to retrieve the service", {name});
                }

                return service;
            }
        };
    };

    let [service, err] = container.get('panels.login');

    if (err) {
        console.log("Unable to load the valid login panel");
    }

    appActions.pushPanel(service);
}


loadApp(document.body);
