import {h, VNode} from "hyperapp";

import {renderZone} from 'fugue-app';

export const NotConnectedTemplate = (props: any, children: Array<VNode>) => (actions: any, state: any): VNode => {
    return <div class="container">
        <div class="row">
            <div class="col">
                { renderZone('header', children) }
            </div>
        </div>
        <div class="row">
            <div class="col">
                { renderZone('content', children) }
            </div>
        </div>
        <div class="row">
            <div class="col">
                { renderZone('footer', children) }
            </div>
        </div>
    </div>
};

export const ConnectedTemplate = (props: any, children: Array<VNode>) => (actions: any, state: any): VNode => {
    return <div class="container">
        <div class="row">
            <div class="col">
                { renderZone('header', children) }
            </div>
        </div>
        <div class="row">
            <div class="col">
                { renderZone('content', children) }
            </div>
        </div>
        <div class="row">
            <div class="col">
                { renderZone('footer', children) }
            </div>
        </div>
    </div>
};

export const ErrorTemplate = (props: any, children: Array<VNode>) => (actions: any, state: any): VNode => {

    return <div class="container">
        <div class="row">
            <div class="col">
                { renderZone('header', children) }
            </div>
        </div>
        <div class="row">
            <div class="col">
                { renderZone('content', children) }
            </div>
        </div>
        <div class="row">
            <div class="col">
                { renderZone('footer', children) }
            </div>
        </div>
    </div>
};