import {h, VNode} from "hyperapp";

import {renderField} from "fugue-form/bootstrap";
import {Form, Field} from "fugue-form";
import {Context, Panel, Zone} from "fugue-app";

import {NotConnectedTemplate} from "../Shared/Templates";
import {Footer} from "../Shared/Footer";

import {ApiClient} from "../../core/api";

enum LoginStatus {Loading, Error, Valid, UserInput}

interface LoginState {
    status: LoginStatus
    model: {
        username: string
        password: string
    }
}

interface LoginActions {
    login: () => LoginState
    updateStatus: (value: LoginStatus) => LoginState
    inlineAction: (value: number) => LoginState
}

export const LoginAction = (api: ApiClient) => () => async (ctx: Context) => {
    const actions = (ctx.actions as LoginActions);

    actions.updateStatus(LoginStatus.Loading);

    const result = await api.signin(ctx.state.model.username, ctx.state.model.password);

    if (result) {
        // valid,
        actions.updateStatus(LoginStatus.Valid);

        ctx.appActions.dispatch("user.connection", {result: true});

        return
    }

    actions.updateStatus(LoginStatus.Error);
};

export const UpdateStatus = () => (value: LoginStatus) =>  (ctx: Context) => {
    ctx.updateState({...ctx.state, status: value});
};

export const createLoginPanel = () => (ctx: Context) => (): VNode => {
    const state = {
        status: LoginStatus.UserInput,
        model: {
            username: 'admin',
            password: 'admin',
        },
    };

    // need to put this in the DI.
    const actions = {
        updateStatus: UpdateStatus(),
        login: LoginAction(ctx.get("gonode.api")),
    };

    return <Panel ctx={ctx} actions={actions} state={state}>
        {(ctx: Context) => {
            const state = (ctx.state as LoginState);
            const actions = (ctx.actions as LoginActions);

            return <NotConnectedTemplate>
                <Zone name={"content"} class={"extra-class"}>

                    { state.status === LoginStatus.Valid ? 
                        <div>Login Successful! </div> : null
                    }

                    { state.status === LoginStatus.Loading ? 
                        <div>Please wait while authentication! </div> : null
                    }

                    <Form setState={(state) => ctx.updateState(state)} getState={() => ctx.state} name={"login"} renderField={renderField}>
                        <Field label="Login:" name={"model.username"} type="text" />
                        <Field label="Password:" name={"model.password"} type="password" />
                        <Field type={"button"}
                               label={"Login !"}
                               class="btn-primary"
                               onclick={() => actions.login()}
                               disabled={state.status === LoginStatus.Loading}
                        />
                    </Form>
                </Zone>
                <Footer ctx={ctx} />
            </NotConnectedTemplate>
        }}
    </Panel>
};
