'use strict';

function GoNodeFactory() {
    this.handlers = {}
}

GoNodeFactory.prototype.get = function (node, component) {
    var handler = this.getHandler(node.type);

    if (component in handler) {
        return handler[component];
    }

    return this.getHandler('default')[component];
}

GoNodeFactory.prototype.add = function (code, type, component) {
    var handler = this.getHandler(code);

    this.handlers[code][type] = component;

    return this;
}

GoNodeFactory.prototype.getHandler = function (code) {
    if (!(code in this.handlers)) {
        this.handlers[code] = {}
    }

    return this.handlers[code];
}

module.exports = GoNodeFactory;
