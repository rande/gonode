/** @jsx React.DOM */

jest.dontMock('../MediaImage.jsx');
jest.dontMock('react-router');
jest.dontMock('react-admin');

describe('Table', function() {
  it('render table', function() {
    var React = require('react/addons');
    var Admin = require('react-admin');

    var Components = require('../MediaImage.jsx');
    var TestUtils = React.addons.TestUtils;

    var StubComponent = Admin.stubRouterContext(Components.ListElement, {node: {}});

    var DomNode = TestUtils.renderIntoDocument(
      <StubComponent />
    );

    // Verify that it's Off by default
    var widget = TestUtils.scryRenderedDOMComponentsWithTag(DomNode, "div");
  });
});
