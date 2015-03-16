'use strict';

var Reflux = require('reflux');

module.exports = function(url) {
  return Reflux.createStore({
     init: function() {
       var ws = new WebSocket(url);
       ws.store = this;

       ws.onmessage = function(evt) {
         console.log("message received:", evt.data);

         if (evt.data == "PING") {
           return;
         }

         ws.store.trigger(JSON.parse(evt.data))
       };

       ws.onclose = function(evt) {
         console.log("Connection closed!!");
       };

       ws.onopen = function(evt) {
         console.log("Connection open!!");
         this.send("hello");
       };
     }
  })
}
