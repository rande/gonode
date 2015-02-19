/** @jsx React.DOM */

jest.dontMock('../List.jsx');
jest.dontMock('lodash');
jest.dontMock('react-router');
jest.dontMock('react-admin');

describe('Table', function() {
  it('render table', function() {
    var React = require('react/addons');
    var Admin = require('react-admin');

    var Component = require('../List.jsx');
    var TestUtils = React.addons.TestUtils;

    var StubComponent = Admin.stubRouterContext(Component, {someProp: 'foo'});

    var DomNode = TestUtils.renderIntoDocument(
      <StubComponent />
    );

    // Verify that it's Off by default
    var widget = TestUtils.scryRenderedDOMComponentsWithTag(DomNode, "div");
  });
});
