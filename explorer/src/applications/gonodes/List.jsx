'use strict';

var React = require('react');
var B = require('react-bootstrap');
var Router = require('react-router');
var RB = require('react-router-bootstrap')
var ReactAdmin = require('react-admin');

var NodeInformationCard = require('./helpers/NodeInformationCard.jsx');
var NodeNotificationCard = require('./helpers/NodeNotificationCard.jsx');

var List = ReactAdmin.createTable({
    getDefaultProps: function () {
        return {
            className: "col-sm-12 col-md-12 main",
            per_page: 32,
            index: "gonodes.list"
        }
    },

    refreshGrid: function (filters) {
        var endpoint = ReactAdmin.Container("gonodes.api.endpoint");

        if (!endpoint) {
            // need to migrate this to a ReFlux store
            return;
        }

        endpoint.get(this.getFilters(filters), function (error, res) {
            if (!this.isMounted()) {
                return;
            }

            if (!res) {
                console.log(error);
            }

            if (res.ok) {
                this.setState({
                    page: res.body.page,
                    per_page: res.body.per_page,
                    base_query: {},
                    next: res.body.next,
                    previous: res.body.previous,
                    elements: res.body.elements
                });
            }
        }.bind(this))
    },

    renderRow: function (node) {
        var factory = ReactAdmin.Container("gonodes.factory");

        var component = factory.get(node, 'list.element');

        // load the form part for the current node
        return React.createElement(component, {node: node, key: "form-list-" + node.uuid})
    }
});

module.exports = List;
