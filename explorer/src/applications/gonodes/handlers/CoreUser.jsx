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
      <B.Col md={6} key={node.uuid}>
        <div className="card card-user">
          <NodeInformationCard node={node} />
          <ReactAdmin.IconCard type="user" />

          <Router.Link to="gonodes.edit" params={node}>{node.name}</Router.Link> <br />

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
        <ReactAdmin.TextInput form={this.props.form} property="data.name" label="Name"/>

        <ReactAdmin.TextInput form={this.props.form} property="data.login" label="Login"/>
        <ReactAdmin.TextInput form={this.props.form} property="data.password" label="Password"/>
      </div>
    );
  }
});

module.exports = {
  ListElement: ListElement,
  FormElement: FormElement
}


