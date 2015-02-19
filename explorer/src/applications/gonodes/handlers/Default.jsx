'use strict';

var React = require('react');
var B = require('react-bootstrap');
var Router = require('react-router');
var ReactAdmin = require('react-admin');

var NodeInformationCard = require('../helpers/NodeInformationCard.jsx');
var NodeNotificationCard = require('../helpers/NodeNotificationCard.jsx');


var ListElement = React.createClass({
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

var FormElement = React.createClass({
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

module.exports = {
  ListElement: ListElement,
  FormElement: FormElement
}


