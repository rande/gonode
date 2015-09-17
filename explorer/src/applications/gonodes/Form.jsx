'use strict';

var React = require('react');
var B = require('react-bootstrap');
var Router = require('react-router');
var Reflux = require('reflux');
var _ = require('lodash');

var ReactAdmin = require('react-admin');

var passwordRegExp = new RegExp("{([a-zA-Z0-9]*)}(.*)");

var CompositePassword = ReactAdmin.createInput({

    mixins: [ReactAdmin.BootstrapInput],

    getInitialState: function () {
        return {
            algo: '',
            password: ''
        }
    },

    componentDidMount: function () {
        this.loadData(this.getValue());
    },

    componentWillReceiveProps: function (nextProps) {
        this.loadData(this.readValue(nextProps.form.state.node, this.props.property))
    },

    loadData: function (value) {
        var match = passwordRegExp.exec(value);

        if (match && match.length == 3) {
            this.setState({
                algo: match[1],
                password: match[2]
            });
        }
    },

    updateAlgo: function (event) {
        this.setValue("{" + event.target.value + "}" + this.state.password)
    },

    updatePassword: function (event) {
        this.setValue("{" + this.state.algo + "}" + event.target.value)
    },

    render: function () {
        return this.renderInput(
            <span>
        <B.Input type="select" name="algo" value={this.state.algo} onChange={this.updateAlgo}>
            <option value="bcrypt">bcrypt (Secure)</option>
            <option value="plain">plain (Not secure, test only)</option>
            <option value="md5">md5 (Not secure, test only)</option>
        </B.Input>

        <B.Input type="password" name="text" value={this.state.password} onChange={this.updatePassword}/>
      </span>
        )
    }
});

var Form = React.createClass({
    mixins: [Router.State, Reflux.ListenerMixin],

    getInitialState: function () {
        return {
            node: {
                name: "loading node ..."
            },
            errors: {}
        };
    },

    componentDidMount: function () {
        var streamStore = ReactAdmin.Container("gonodes.stores.stream");

        this.listenTo(streamStore, this.onStreamUpdate);

        this.loadData();
    },

    onStreamUpdate: function (data) {
        console.log("receive an update from the node");

        setTimeout(function () {
            if (data.revision > this.state.node.revision) {
                console.log("new version available");

                ReactAdmin.Status.Action("warning", "The object has been updated by an external user, so the current form has been updated!", 5000);

                this.loadData();
            }
        }.bind(this), 500); // wait a bit before reloading, avoid race condition

    },

    refreshView: function () {
        this.setState({
            node: this.state.node
        });
    },

    componentWillReceiveProps: function () {
        this.loadData();
    },

    submit: function () {
        var endpoint = ReactAdmin.Container("gonodes.api.endpoint");

        endpoint.put('/' + this.getParams().uuid, this.state.node, function (res) {
            if (res.ok) {

                ReactAdmin.Status.Action("success", "The subject has been saved!");

                this.setState({
                    node: res.body,
                    errors: {}
                });
            } else {

                ReactAdmin.Status.Action("danger", "An error occurs while saving", 4000);

                this.setState({
                    errors: res.body
                });
            }
        }.bind(this));
    },

    loadData: function () {
        var endpoint = ReactAdmin.Container("gonodes.api.endpoint");
        var update = ReactAdmin.Container("admin.reflux.action.update");

        endpoint.get('/' + this.getParams().uuid, function (error, response) {
            if (!this.isMounted()) {
                return;
            }

            if (!response) {
                console.log("response not defined", error);

                update("KO", "An error occurs while loading data");
                return;
            }

            if (!response.ok) {
                // an error occurs
                update("KO", "An unexpected error occurs while loading data");
                return;
            }

            this.setState({
                node: response.body,
                errors: {}
            });

        }.bind(this));
    },

    render: function () {
        var factory = ReactAdmin.Container("gonodes.factory");

        var CustomForm = false;
        var component = factory.get(this.state.node, 'form');

        // load the form part for the current node
        if (component) {
            CustomForm = React.createElement(component, {form: this})
        }

        return (
            <div className="col-sm-12 col-md-12 main">
                <h2 className="sub-header">Edit {this.state.node.name} (rev: {this.state.node.revision}) </h2>

                <form>
                    <div className="row">
                        <div className="col-sm-6">
                            <ReactAdmin.TextInput form={this} property="name" label="Name" help="Enter the name"/>
                            <ReactAdmin.TextInput form={this} property="slug" label="Slug" help="Enter the slug"/>
                        </div>
                        <div className="col-sm-6">
                            <ReactAdmin.BooleanSelect form={this} property="enabled" label="enabled">
                                <option value="1">Yes</option>
                                <option value="0">No</option>
                            </ReactAdmin.BooleanSelect>

                            <ReactAdmin.NumberSelect form={this} property="status" label="Status">
                                <option value="0">New</option>
                                <option value="1">Draft</option>
                                <option value="2">Completed</option>
                                <option value="3">Validated</option>
                            </ReactAdmin.NumberSelect>

                            <ReactAdmin.NumberInput form={this} property="weight" help="error message " label="Weight"/>
                        </div>
                    </div>

                    {CustomForm}
                </form>

                <B.Button bsStyle="primary" onClick={this.submit}>Save</B.Button>
            </div>
        );
    }
});

module.exports = Form;

