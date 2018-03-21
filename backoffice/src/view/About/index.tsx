import {h, VNode} from "hyperapp";
import {Context, Panel, Zone} from "fugue-app";
import {NotConnectedTemplate} from "../Shared/Templates";

export const createAboutPanel = () => (ctx: Context) => (): VNode => {
    return <Panel ctx={ctx}>
        {(ctx: Context) => {

            return <NotConnectedTemplate>
                <Zone name={"content"}>
                    <h1>About the page</h1>
                    <p>
                        Hello, this is a solution to create application <br/>
                        based on hyperapp.
                    </p>

                    <p>
                        <a onclick={(e: Event) => {
                            console.log("onclick", ctx);
                            //ctx.appActions.popPanel();
                            e.preventDefault();
                        }} href="/about">Close About Page</a>
                    </p>
                </Zone>
            </NotConnectedTemplate>
        }}
    </Panel>
};