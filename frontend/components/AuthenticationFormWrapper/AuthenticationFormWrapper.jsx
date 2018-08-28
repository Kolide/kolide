import React, { Component, PropTypes } from 'react';

const baseClass = 'auth-form-wrapper';

class AuthenticationFormWrapper extends Component {
  static propTypes = {
    children: PropTypes.node,
  };

  render () {
    const { children } = this.props;

    return (
      <div className={baseClass}>
        <img alt="Kolide vertical logo" src="/assets/images/kolide-logo-vertical.svg" className={`${baseClass}__logo`} />
        {children}
      </div>
    );
  }
}

export default AuthenticationFormWrapper;
