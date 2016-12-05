import React, { Component, PropTypes } from 'react';
import { Link } from 'react-router';
import classnames from 'classnames';

const baseClass = 'stacked-white-boxes';

class StackedWhiteBoxes extends Component {
  static propTypes = {
    children: PropTypes.element,
    headerText: PropTypes.string,
    className: PropTypes.string,
    leadText: PropTypes.string,
    previousLocation: PropTypes.string,
  };

  renderBackButton = () => {
    const { previousLocation } = this.props;

    if (!previousLocation) return false;

    return (
      <div className={`${baseClass}__back`}>
        <Link to={previousLocation} className={`${baseClass}__back-link`}>╳</Link>
      </div>
    );
  }

  renderHeader = () => {
    const { headerText } = this.props;

    return (
      <div className={`${baseClass}__header`}>
        <p className={`${baseClass}__header-text`}>{headerText}</p>
      </div>
    );
  }

  render () {
    const { children, className, leadText } = this.props;
    const { renderBackButton, renderHeader } = this;

    const boxClass = classnames(
      baseClass,
      className
    );

    return (
      <div className={boxClass}>
        <div className={`${baseClass}__box`}>
          {renderBackButton()}
          {renderHeader()}
          <p className={`${baseClass}__box-text`}>{leadText}</p>
          {children}
        </div>
      </div>
    );
  }
}

export default StackedWhiteBoxes;
