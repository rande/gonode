'use strict';

var React = require('react');
var ReactAdmin = require('react-admin');

var NodeNotificationCard = React.createClass({
  propTypes: {
    node: React.PropTypes.object.isRequired
  },

  render: function() {
    return (
      <ReactAdmin.NotificationCard>
        <span className={this.props.node.enabled ? 'glyphicon glyphicon-ok-circle' : 'glyphicon glyphicon-ban-circle'} />
      </ReactAdmin.NotificationCard>
    );
  }
});

module.exports = NodeNotificationCard;
