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
      <div className="row">
        <div className="col-sm-6">

          <div className="col-sm-6">
            <ReactAdmin.TextInput form={this.props.form} property="data.firstname" label="Firstname"/>
          </div>
          <div className="col-sm-6">
            <ReactAdmin.TextInput form={this.props.form} property="data.lastname" label="Lastname"/>
          </div>

          <ReactAdmin.TextInput form={this.props.form} property="data.email" label="Email"/>
          <ReactAdmin.TextInput form={this.props.form} property="data.login" label="Login"/>
          <ReactAdmin.TextInput form={this.props.form} property="data.newpassword" label="New Password" help="Set a new password for the user"/>

          <ReactAdmin.TextInput form={this.props.form} property="data.locale" label="Locale"/>
          <ReactAdmin.TextInput form={this.props.form} property="data.timezone" label="Timezone"/>
        </div>

        <div className="col-sm-6">
          <ReactAdmin.BooleanInput form={this.props.form} property="data.locked" label="Deleted" />
          <ReactAdmin.BooleanInput form={this.props.form} property="data.enabled" label="Enabled" />
          <ReactAdmin.BooleanInput form={this.props.form} property="data.expired" label="Expired" />

          <ReactAdmin.Radio form={this.props.form} property="data.gender" name="Gender" value="m" label="Male" />
          <ReactAdmin.Radio form={this.props.form} property="data.gender" name="Gender" value="f" label="Female" />
        </div>
      </div>
    );
  }
});

module.exports = {
  ListElement: ListElement,
  FormElement: FormElement
}


