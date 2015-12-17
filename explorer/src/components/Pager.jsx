import React, { Component, PropTypes } from 'react';
import { FormattedMessage }            from 'react-intl';


class Pager extends Component {
    handleChange() {
        const { onChange } = this.props;
        onChange({
            perPage: parseInt(this.refs.perPage.value)
        });
    }

    render() {
        const { perPageOptions, perPage } = this.props;

        return (
            <div className="pager">
                <FormattedMessage id="pager.per_page"/>
                <select ref="perPage" value={perPage} onChange={this.handleChange.bind(this)}>
                    {perPageOptions.map(perPageOption => (
                        <option key={perPageOption} value={perPageOption}>{perPageOption}</option>
                    ))}
                </select>
            </div>
        );
    }
}

Pager.propTypes = {
    perPageOptions: PropTypes.array.isRequired,
    perPage:        PropTypes.number.isRequired,
    page:           PropTypes.number.isRequired,
    onChange:       PropTypes.func.isRequired
};


export default Pager;
