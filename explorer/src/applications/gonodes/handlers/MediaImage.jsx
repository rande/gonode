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
        <div className="card card-media">
          <NodeInformationCard node={node} />
          <ReactAdmin.IconCard type="picture-o" />

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

  refresh: function(event) {
    ReactAdmin.WriteValue(this.props.form.state.node, "meta.source_status", 0 /* Init */)

    this.props.form.submit();
  },

  render: function() {
    var node = this.props.form.state.node;

    var RefreshButton = false;
    if (node.uuid && node.meta.source_status != 1 /* Update */) {
        RefreshButton = <B.Button bsStyle="primary" onClick={this.refresh}><i className="fa fa-refresh"></i>Refresh</B.Button>
    }

    var Preview = false;

    if (node.uuid && node.meta.source_status == 2 /* Done */) {
        Preview = <img src={ "/nodes/" + node.uuid + "?raw"} width="200px" />
    }

    return (
      <div>
        <div className="col-sm-6">
          <ReactAdmin.TextInput form={this.props.form} property="data.name" label="Name"/>

          <ReactAdmin.TextInput form={this.props.form} property="data.reference" label="Reference"/>
          <ReactAdmin.TextInput form={this.props.form} property="data.source_url" label="Source URL"/>
        </div>

        <div className="col-sm-6">

          <p>
            {Preview}
            <ul>
              <li>Status: {node.meta.source_status}</li>
              <li>Error: {node.meta.source_error}</li>
              <li>{RefreshButton}</li>
            </ul>
            <ul>
              <li>Dimension: {node.meta.width}x{node.meta.height}</li>
              <li>Size: {node.meta.size}</li>
              <li>ContentType: {node.meta.content_type}</li>
              <li>Length: {node.meta.length}</li>
            </ul>
          </p>
        </div>
      </div>
    );
  }
});

module.exports = {
  ListElement: ListElement,
  FormElement: FormElement
}


