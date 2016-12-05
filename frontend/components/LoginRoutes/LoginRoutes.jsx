import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { RouteTransition } from 'react-router-transition';

import { hideBackgroundImage, showBackgroundImage } from 'redux/nodes/app/actions';
import LoginPage from 'pages/LoginPage';
import redirectLocationInterface from 'interfaces/redirect_location';
import Footer from 'components/Footer';

export class LoginRoutes extends Component {
  static propTypes = {
    children: PropTypes.element,
    dispatch: PropTypes.func,
    location: redirectLocationInterface,
  };

  componentWillMount () {
    const { dispatch } = this.props;

    dispatch(showBackgroundImage);
  }

  componentWillUnmount () {
    const { dispatch } = this.props;

    dispatch(hideBackgroundImage);
  }

  render () {
    const { children, location: { pathname } } = this.props;

    return (
      <div className="login-routes">
        <LoginPage pathname={pathname} />
        <RouteTransition
          pathname={pathname}
          atEnter={{
            scale: 1.3,
            opacity: 0,
          }}
          atLeave={{
            scale: 1.3,
            opacity: 0,
          }}
          atActive={{
            scale: 1,
            opacity: 1,
          }}
          mapStyles={styles => ({
            opacity: styles.opacity,
            transform: `scale(${styles.scale})`,
          })}
        >
          {children}
        </RouteTransition>
        <Footer />
      </div>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  const { location } = ownProps;

  return { location };
};

export default connect(mapStateToProps)(LoginRoutes);
