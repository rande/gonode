'use strict';

var Reflux = require('reflux');

module.exports = function() {
  return Reflux.createStore({
     init: function() {
       this.history = []
     },

     addHistory: function(data) {
       this.history.unshift(data);

       if (this.history.length > 32) {
         this.history.pop();
       }

       this.trigger(this.history);
     }
  })
}
