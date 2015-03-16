'use strict';

var React = require('react');
var B = require('react-bootstrap');
var Router = require('react-router');
var ReactAdmin = require('react-admin');

var NodeInformationCard = require('../helpers/NodeInformationCard.jsx');
var NodeNotificationCard = require('../helpers/NodeNotificationCard.jsx');

module.exports.NotificationElement = React.createClass({
  propTypes: {
    subject: React.PropTypes.string.isRequired,
    type: React.PropTypes.string.isRequired,
    revision: React.PropTypes.number.isRequired,
    date: React.PropTypes.string.isRequired,
    action: React.PropTypes.string.isRequired
  },
  render: function() {
      return (
        <div className="notification-element">
          <Router.Link to="gonodes.edit" params={{uuid: this.props.subject}}>{this.props.name}</Router.Link> <br />
          Action: {this.props.action} - ({this.props.type} / #{this.props.revision})
        </div>
      );
  }
})

module.exports.ListElement = React.createClass({
  propTypes: {
    node: React.PropTypes.object.isRequired
  },

  render: function() {
    var node = this.props.node;
    return (
      <B.Col md={6}>
        <div className="card">
          <NodeInformationCard node={node} />

          <ReactAdmin.IconCard type="circle-thin" />
          <Router.Link to="gonodes.edit" params={node}>{node.name?node.name:node.uuid}</Router.Link> <br />

          Type: {node.type}

          <NodeNotificationCard node={node} />
        </div>
      </B.Col>
    );
  }
});

module.exports.FormElement = React.createClass({
  propTypes: {
    form: React.PropTypes.object.isRequired
  },

  render: function() {
    return (
      <div>
        <p>No handler defined for this node</p>
      </div>
    );
  }
});
