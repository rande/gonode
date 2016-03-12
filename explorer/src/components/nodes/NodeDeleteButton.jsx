import React, { Component, PropTypes } from 'react';
import { connect }                     from 'react-redux';
import { FormattedMessage }            from 'react-intl';
import ReactCSSTransitionGroup         from 'react-addons-css-transition-group';


class NodeDeleteButton extends Component {
    static displayName = 'NodeDeleteButton';

    static propTypes = {
        size: PropTypes.string.isRequired,
        uuid: PropTypes.string.isRequired
    };

    static defaultProps = {
        size: ''
    };

    constructor(props) {
        super(props);

        this.state = { confirm: false };

        this.handleConfirm = this.handleConfirm.bind(this);
        this.handleCancel  = this.handleCancel.bind(this);
    }

    handleConfirm() {
        this.setState({ confirm: true });
    }

    handleCancel() {
        this.setState({ confirm: false });
    }

    render() {
        const { size }    = this.props;
        const { confirm } = this.state;

        let buttonClasses = 'button';
        if (size !== '') {
            buttonClasses = `${buttonClasses} button-${size}`;
        }

        let buttons;
        if (!confirm) {
            buttons = [(
                <span key="delete" className={buttonClasses} onClick={this.handleConfirm}>
                    <i className="fa fa-trash" />
                    <FormattedMessage id="node.delete.link"/>
                </span>
            )];
        } else {
            buttons = [
                <span key="confirm" className={buttonClasses}>
                    <i className="fa fa-warning" />
                    <FormattedMessage id="node.delete.confirm"/>
                </span>,
                <span key="yes" className={buttonClasses}>
                    <FormattedMessage id="yes"/>
                </span>,
                <span key="noe" className={buttonClasses} onClick={this.handleCancel}>
                    <FormattedMessage id="no"/>
                </span>
            ];
        }

        return (
            <span className="node_delete">
                <ReactCSSTransitionGroup
                    transitionName="node_delete"
                    transitionEnterTimeout={400}
                    transitionLeaveTimeout={1}
                >
                    {buttons}
                </ReactCSSTransitionGroup>
            </span>
        );
    }
}

const mapDispatchToProps = dispatch => ({

});


export default connect(() => ({}), mapDispatchToProps)(NodeDeleteButton);
