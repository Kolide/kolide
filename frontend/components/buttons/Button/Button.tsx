import * as React from 'react';
const classnames = require('classnames');

const baseClass = 'button';

interface IButtonProps {
  className: string;
  disabled: boolean;
  onClick: (evt: React.MouseEvent<HTMLButtonElement>) => boolean;
  size: string;
  text: string;
  type: string;
  variant: string;
}

interface IButtonState {}

class Button extends React.Component<IButtonProps, IButtonState> {
  static defaultProps = {
    variant: 'default',
    size: '',
  };

  handleClick = (evt: React.MouseEvent<HTMLButtonElement>) => {
    const { disabled, onClick } = this.props;

    if (disabled) {
      return false;
    }

    if (onClick) {
      onClick(evt);
    }

    return false;
  }

  render () {
    const { handleClick } = this;
    const { className, disabled, size, text, type, variant } = this.props;
    const fullClassName = classnames(`${baseClass}--${variant}`, className, {
      [baseClass]: variant !== 'unstyled',
      [`${baseClass}--disabled`]: disabled,
      [`${baseClass}--${size}`]: size,
    });

    return (
      <button
        className={fullClassName}
        disabled={disabled}
        onClick={handleClick}
        type={type}
      >
        {text}
      </button>
    );
  }
}

export default Button;
