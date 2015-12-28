import React, { Component, PropTypes } from 'react';
import { Link }                        from 'react-router';
import { FormattedMessage }            from 'react-intl';


class Pager extends Component {
    handleChange() {
        const { onChange } = this.props;
        onChange({
            perPage: parseInt(this.refs.perPage.value)
        });
    }

    render() {
        const {
            perPageOptions,
            perPage,
            page,
            previousPage,
            nextPage
        } = this.props;

        let previousPageButton = null;
        if (previousPage) {
            previousPageButton = (
                <Link to={`/nodes?pp=${perPage}&p=${previousPage}`} className="button pager_previous">
                    <i className="fa fa-chevron-left"/>
                </Link>
            );
        }

        let nextPageButton = null;
        if (nextPage) {
            nextPageButton = (
                <Link to={`/nodes?pp=${perPage}&p=${nextPage}`} className="button pager_next">
                    <i className="fa fa-chevron-right"/>
                </Link>
            );
        }

        return (
            <div className="pager">
                {previousPageButton}
                <span className="pager_page">
                    <FormattedMessage id="pager.page" values={{ page }}/>
                </span>
                <FormattedMessage id="pager.per_page"/>
                <select ref="perPage" value={perPage} onChange={this.handleChange.bind(this)}>
                    {perPageOptions.map(perPageOption => (
                        <option key={perPageOption} value={perPageOption}>{perPageOption}</option>
                    ))}
                </select>
                {nextPageButton}
            </div>
        );
    }
}

Pager.propTypes = {
    perPageOptions: PropTypes.array.isRequired,
    perPage:        PropTypes.number.isRequired,
    page:           PropTypes.number.isRequired,
    previousPage:   PropTypes.number,
    nextPage:       PropTypes.number,
    onChange:       PropTypes.func.isRequired
};


export default Pager;
