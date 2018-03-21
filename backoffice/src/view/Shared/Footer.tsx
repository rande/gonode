import {h} from "hyperapp";

import {Context, Zone} from "fugue-app";

import {createAboutPanel} from "../About";

export const Footer = ({ctx}: { ctx: Context }) => {
    return <Zone name={"footer"}>
        <div>
            The is the footer of major corporate company
            <br/>
            <a onclick={(e: Event) => {
                ctx.appActions.pushPanel(createAboutPanel);
                e.preventDefault();
            }} href="/about">About page</a>
        </div>
    </Zone>
};