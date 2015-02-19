'use strict'

var React = require('react');
var ReactAdmin = require('react-admin');
var Router = require('react-router');

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
  .add('default', 'form', DefaultHandler.FormElement)
  .add('default', 'list.element', DefaultHandler.ListElement)

  // media stuff
  .add('media.image', 'form', MediaImageHandler.FormElement)
  .add('media.image', 'list.element', MediaImageHandler.ListElement)
  // core user stuff
  .add('core.user', 'form', CoreUserHandler.FormElement)
  .add('core.user', 'list.element', CoreUserHandler.ListElement)
;

// register the factory into the global container
ReactAdmin.Container("gonodes.factory", factory);

module.exports = {
  getRoutes: getRoutes
}
