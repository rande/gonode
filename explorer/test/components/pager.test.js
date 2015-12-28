import React          from 'react';
import ReactTestUtils from 'react-addons-test-utils';
import { Link }       from 'react-router';
import expect         from 'expect';
import Pager          from '../../src/components/Pager.jsx';

const shallowRenderer = ReactTestUtils.createRenderer();

describe('pager component', () => {
    const perPageOptions = [5, 10, 15];

    it('should allow to define custom per page values', () => {
        shallowRenderer.render(React.createElement(Pager, {
            perPage:        10,
            page:           1,
            perPageOptions: perPageOptions,
            onChange:       () => {}
        }));

        const render = shallowRenderer.getRenderOutput();

        const perPageSelector = render.props.children[3];

        expect(perPageSelector.type).toEqual('select');
        expect(perPageSelector.props.children.length).toEqual(perPageOptions.length);

        perPageOptions.forEach((perPage, pos) => {
            expect(perPageSelector.props.children[pos].props.value).toEqual(perPage);
        });
    });

    it('should render no previous/next links if none provided', () => {
        shallowRenderer.render(React.createElement(Pager, {
            perPage:        10,
            page:           1,
            perPageOptions: perPageOptions,
            onChange:       () => {}
        }));

        const render = shallowRenderer.getRenderOutput();

        expect(render.props.children[0]).toBe(null);
        expect(render.props.children[4]).toBe(null);
    });

    it('should render previous link if \'previousPage\' provided', () => {
        shallowRenderer.render(React.createElement(Pager, {
            perPage:        10,
            page:           2,
            previousPage:   1,
            perPageOptions: perPageOptions,
            onChange:       () => {}
        }));

        const render = shallowRenderer.getRenderOutput();

        const previousLink = render.props.children[0];

        expect(previousLink).toNotBe(null);
        expect(previousLink.props.className).toEqual('button pager_previous');
        expect(previousLink.props.to).toEqual('/nodes?pp=10&p=1');
    });

    it('should render next link if \'nextPage\' provided', () => {
        shallowRenderer.render(React.createElement(Pager, {
            perPage:        20,
            page:           2,
            nextPage:       3,
            perPageOptions: perPageOptions,
            onChange:       () => {}
        }));

        const render = shallowRenderer.getRenderOutput();

        const nextLink = render.props.children[4];

        expect(nextLink).toNotBe(null);
        expect(nextLink.props.className).toEqual('button pager_next');
        expect(nextLink.props.to).toEqual('/nodes?pp=20&p=3');
    });

    it('should display current page', () => {
        shallowRenderer.render(React.createElement(Pager, {
            perPage:        20,
            page:           13,
            perPageOptions: perPageOptions,
            onChange:       () => {}
        }));

        const render = shallowRenderer.getRenderOutput();

        const currentPage = render.props.children[1];

        expect(currentPage.props.className).toEqual('pager_page');
        expect(currentPage.props.children.props.id).toEqual('pager.page');
        expect(currentPage.props.children.props.values.page).toEqual(13);
    });
});