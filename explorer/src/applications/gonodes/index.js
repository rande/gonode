'use strict'

var React = require('react');
var ReactAdmin = require('react-admin');
var Router = require('react-router');
var _ = require('lodash');

var Form = require('./Form.jsx');
var List = require('./List.jsx');

/**
 * This is used to build the nested view required by React Router
 */
var View = React.createClass({
  render: function() {
    return <Router.RouteHandler />
  }
});

/**
 * Define the routes required to list or edit the node
 */
function getRoutes() {
  return <Router.Route name="gonodes" handler={View} >

    <Router.Route name="gonodes.list"  path="list" handler={List} />
    <Router.Route name="gonodes.edit"  path="edit/:uuid" handler={Form} />

    <Router.DefaultRoute handler={List} />
  </Router.Route>
}

/**
 * Configure form handler
 */
var MediaImageHandler = require('./handlers/MediaImage.jsx');
var CoreUserHandler = require('./handlers/CoreUser.jsx');
var DefaultHandler = require('./handlers/Default.jsx');
var GoNodeFactory = require('./factory');

// the factory defines the related components used in the explorer
var factory = (new GoNodeFactory())
  // default
  .add('default', 'form',                 DefaultHandler.FormElement)
  .add('default', 'list.element',         DefaultHandler.ListElement)
  .add('default', 'notification.element', DefaultHandler.NotificationElement)

  // media stuff
  .add('media.image', 'form',                 MediaImageHandler.FormElement)
  .add('media.image', 'list.element',         MediaImageHandler.ListElement)
  .add('media.image', 'notification.element', DefaultHandler.NotificationElement)

  // core user stuff
  .add('core.user', 'form',                 CoreUserHandler.FormElement)
  .add('core.user', 'list.element',         CoreUserHandler.ListElement)
  .add('core.user', 'notification.element', DefaultHandler.NotificationElement)
;

var location = window.location;
var streamUrl = (location.protocol === "https:" ? "wss://" :  "ws://" +  location.host + "/nodes/stream");

var streamStore = require('./stores/Stream.jsx')(streamUrl);
var historyStore = require('./stores/History.jsx')();


historyStore.listenTo(streamStore, historyStore.addHistory);

streamStore.listen(function(data) {
    var component = factory.get(data.type, 'notification.element');
    if (!component) {
        // TODO: log error ?
        console.log("Empty component for notification:", data);
        return;
    }

    ReactAdmin.Notification.Action(component, data || {})
});

// register the factory into the global container
ReactAdmin.Container("gonodes.factory", factory);
ReactAdmin.Container("gonodes.stores.stream", streamStore)
ReactAdmin.Container("gonodes.stores.history", historyStore)

module.exports = {
  getRoutes: getRoutes
}
