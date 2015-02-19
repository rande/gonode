'use strict';

var React = require('react');
var Router = require('react-router');
var B = require('react-bootstrap');
var RB = require('react-router-bootstrap');
var ReactAdmin = require('react-admin');

var Dashboard = require('layouts/Dashboard.jsx');

var GoNodesApplication = require('applications/gonodes');

ReactAdmin.Container()
  .set("gonodes.api.endpoint", new ReactAdmin.EndPoint('/nodes', {'Accept':'application/json'}));

/**
 *  Define the global layout of your application
 *  The Router.RouteHandler call is required to render sub view defined
 *  on each defined applications.
 */
var Header = React.createClass({
  render: function() {
    return (
      <B.Navbar inverse={true} fixedTop={true} fluid={true} role="navigation">
        <div className="navbar-header">
          <button type="button" className="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
            <span className="sr-only">Toggle navigation</span>
            <span className="icon-bar"></span>
            <span className="icon-bar"></span>
            <span className="icon-bar"></span>
          </button>
          <Router.Link to="homepage" className="navbar-brand">React Admin</Router.Link>
        </div>

        <div id="navbar" className="navbar-collapse collapse">
          <B.Nav activeKey={1} right={true} navbar={true}>
            <RB.NavItemLink eventKey={1} to="gonodes.list">Nodes</RB.NavItemLink>
          </B.Nav>

          <form className="navbar-form navbar-right">
            <input type="text" bsClass="form-control" placeholder="Search..." />
          </form>
        </div>
      </B.Navbar>
    );
  }
});

var App = React.createClass({
  render: function () {
    return (
      <div>
        <Header />
        <div className="container-fluid">
          <div className="row">
            <Router.RouteHandler />
          </div>
        </div>
      </div>
    );
  }
});

/**
 * This is the default 404 page when a page does not exist
 * Feel free to build your own ...
 */
var DebugRouter = React.createClass({
  mixins: [Router.State],
  render: function() {
    var message = this.getPathname();

    return <div className="col-sm-12 col-md-12 main">
      <h2>Oh Oh ... Page does not exist </h2>
      <h3>underwood error code is 404</h3>
      <p>
        {message}
      </p>

    </div>
  }
});

/**
 * Register the different view or applications to the router
 * You can append any routes you want
 */
var routes = (
  <Router.Route name="homepage" path="/" handler={App}>
    <Router.DefaultRoute handler={Dashboard}/>
    <Router.NotFoundRoute handler={DebugRouter}/>

    {GoNodesApplication.getRoutes()}

  </Router.Route>
);

/**
 * Start the application, the app id is set in the index.html page
 */
Router.run(routes, function (Handler) {
  React.render(<Handler/>, document.getElementById("app"));
});
