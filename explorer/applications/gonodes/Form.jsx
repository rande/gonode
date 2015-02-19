'use strict';

var React = require('react');
var B = require('react-bootstrap');
var Router = require('react-router');
var _ = require('lodash');

var ReactAdmin = require('react-admin');

var passwordRegExp = new RegExp("{([a-zA-Z0-9]*)}(.*)");

var CompositePassword = ReactAdmin.createInput({

  mixins: [ReactAdmin.BootstrapInput],

  getInitialState: function() {
    return {
      algo: '',
      password: ''
    }
  },

  componentDidMount: function() {
    this.loadData(this.getValue());
  },

  componentWillReceiveProps: function(nextProps) {
    this.loadData(this.readValue(nextProps.form.state.node, this.props.property))
  },

  loadData: function(value) {
    var match = passwordRegExp.exec(value);

    if (match && match.length == 3) {
      this.setState({
        algo: match[1],
        password: match[2]
      });
    }
  },

  updateAlgo: function(event) {
    this.setValue("{" + event.target.value + "}" + this.state.password)
  },

  updatePassword: function(event) {
    this.setValue("{" + this.state.algo + "}" + event.target.value)
  },

  render: function() {
    return this.renderInput(
      <span>
        <B.Input type="select" name="algo" value={this.state.algo} onChange={this.updateAlgo} >
          <option value="bcrypt">bcrypt (Secure)</option>
          <option value="plain">plain (Not secure, test only)</option>
          <option value="md5">md5 (Not secure, test only)</option>
        </B.Input>

        <B.Input type="password" name="text" value={this.state.password} onChange={this.updatePassword} />
      </span>
    )
  }
});


var Form = React.createClass({
  mixins: [Router.State],

  getInitialState: function() {
    return {
      node: {
        name: "loading node ..."
      },
      errors: {}
    };
  },

  componentDidMount: function() {
    this.loadData();
  },

  refreshView: function() {
    this.setState({
      node: this.state.node
    });
  },

  submit: function() {
    var endpoint = ReactAdmin.Container("gonodes.api.endpoint");

    endpoint.put('/' + this.getParams().uuid, this.state.node, function(res) {
      if (res.ok) {
        this.setState({
          node: res.body,
          errors: {}
        });
      } else {
        this.setState({
          errors: res.body
        });
      }
    }.bind(this));
  },

  loadData: function() {
    var endpoint = ReactAdmin.Container("gonodes.api.endpoint");

    endpoint.get('/' + this.getParams().uuid, function(error, response) {
      if (!this.isMounted()) {
        return;
      }

      if (!response) {
        throw e;
        console.log("response not defined");
        return;
      }

      if (!response.ok) {
        // an error occurs
        return;
      }

      this.setState({
        node: response.body,
        errors: {}
      });

    }.bind(this));
  },

  render: function() {
    var factory = ReactAdmin.Container("gonodes.factory");

    var CustomForm = false;
    var component = factory.get(this.state.node, 'form');

    // load the form part for the current node
    if (component) {
      CustomForm = React.createElement(component, {form:this})
    }

    return (
      <div className="col-sm-12 col-md-12 main">
        <h2 className="sub-header">Edit {this.state.node.name} (rev: {this.state.node.revision}) </h2>

        <form>
          <div className="row">
            <div className="col-sm-6">
              <ReactAdmin.TextInput form={this} property="data.password" label="Password"/>
            </div>
            <div className="col-sm-6">
              <CompositePassword form={this} property="data.password" label="Password" help="Configure your password"/>
            </div>
          </div>

          <hr />
          <div className="row">
            <div className="col-sm-6">
              <ReactAdmin.TextInput form={this} property="name" label="Name" help="Enter the name"/>
            </div>
            <div className="col-sm-6">
              <ReactAdmin.TextAreaInput form={this} property="name" label="Name"/>
            </div>
          </div>

          <hr />
          <ReactAdmin.TextInput form={this} property="slug" label="Slug" help="Enter the slug"/>

          <hr />
          <div className="row">
            <div className="col-sm-6">
              <ReactAdmin.BooleanRadio form={this} property="enabled" name="enabled" value="true" label="Enabled" help="Enable the node" />
              <ReactAdmin.BooleanRadio form={this} property="enabled" name="enabled" value="false" label="Disabled" help="Disable the node" />
            </div>

            <div className="col-sm-6">
              <ReactAdmin.BooleanSelect form={this} property="enabled" label="enabled">
                  <option value="1">Yes</option>
                  <option value="0">No</option>
              </ReactAdmin.BooleanSelect>
            </div>
          </div>

          <hr />
          <ReactAdmin.BooleanInput form={this} property="deleted" label="Deleted" />

          <hr />
          <div className="row">
            <div className="col-sm-4">
              <ReactAdmin.NumberSelect form={this} property="status" label="Status">
                  <option value="1">Status 1</option>
                  <option value="2">Status 2</option>
                  <option value="3">Status 3</option>
                  <option value="4">Status 4</option>
              </ReactAdmin.NumberSelect>
            </div>

            <div className="col-sm-4">
              <ReactAdmin.NumberRadio form={this} property="status" name="status" value="1" label="Status 1" />
              <ReactAdmin.NumberRadio form={this} property="status" name="status" value="2" label="Status 2" />
              <ReactAdmin.NumberRadio form={this} property="status" name="status" value="3" label="Status 3" />
              <ReactAdmin.NumberRadio form={this} property="status" name="status" value="4" label="Status 4" />
            </div>

            <div className="col-sm-4">
              <ReactAdmin.NumberInput form={this} property="status" help="numeral status" label="Status" />
            </div>
          </div>

          <hr />
          <ReactAdmin.NumberInput form={this} property="weight" help="error message " label="Weight" />

          <hr />
          {CustomForm}
        </form>

        <B.Button bsStyle="primary" onClick={this.submit}>Save</B.Button>
      </div>
    );
  }
});

module.exports = Form;

