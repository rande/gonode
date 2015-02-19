'use strict';

var React = require('react');
var ReactAdmin = require('react-admin');

var NodeInformationCard = React.createClass({
  propTypes: {
    node: React.PropTypes.object.isRequired
  },

  render: function() {
    return (
      <ReactAdmin.InformationCard>
        {this.props.node.uuid} rev: {this.props.node.revision}
      </ReactAdmin.InformationCard>
    )
  }
});


module.exports = NodeInformationCard;
