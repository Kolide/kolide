import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';

import Icon from 'components/Icon';

const baseClass = 'modal';

class Modal extends Component {
  static propTypes = {
    children: PropTypes.node,
    className: PropTypes.string,
    onExit: PropTypes.func,
    title: PropTypes.oneOfType([PropTypes.string, PropTypes.node]),
  };

  render () {
    const { children, className, onExit, title } = this.props;
    const modalContainerClassName = classnames(`${baseClass}__modal_container`, className);

    return (
      <div className={`${baseClass}__background`}>
        <div className={modalContainerClassName}>
          <div className={`${baseClass}__header`}>
            <span>{title}</span>
            <div className={`${baseClass}__ex`}>
              <button className={`button button--unstyled`} onClick={onExit}>
                <Icon name="x" />
              </button>
            </div>
          </div>
          <div className={`${baseClass}__content`}>
            {children}
          </div>
        </div>
      </div>
    );
  }
}

export default Modal;
