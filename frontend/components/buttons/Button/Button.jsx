import React, { Component, PropTypes } from 'react';
import radium from 'radium';

import componentStyles from './styles';

class Button extends Component {
  static propTypes = {
    className: PropTypes.string,
    disabled: PropTypes.bool,
    onClick: PropTypes.func,
    style: PropTypes.object, // eslint-disable-line react/forbid-prop-types
    text: PropTypes.string,
    type: PropTypes.string,
    variant: PropTypes.string,
  };

  static defaultProps = {
    style: {},
    variant: 'default',
  };

  handleClick = (evt) => {
    const { disabled, onClick } = this.props;

    if (disabled) return false;

    onClick(evt);

    return false;
  }

  render () {
    const { handleClick } = this;
    const { className, style, text, type, variant } = this.props;

    return (
      <button
        className={className}
        onClick={handleClick}
        style={[componentStyles[variant], style]}
        type={type}
      >
        {text}
      </button>
    );
  }
}

export default radium(Button);
