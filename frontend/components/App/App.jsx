import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { noop } from 'lodash';
import classnames from 'classnames';

import { fetchCurrentUser } from '../../redux/nodes/auth/actions';
import Footer from '../Footer';
import { getConfig } from '../../redux/nodes/app/actions';
import { authToken } from '../../utilities/local';
import userInterface from '../../interfaces/user';

export class App extends Component {
  static propTypes = {
    children: PropTypes.element,
    dispatch: PropTypes.func,
    showBackgroundImage: PropTypes.bool,
    user: userInterface,
  };

  static defaultProps = {
    dispatch: noop,
  };

  componentWillMount () {
    const { dispatch, user } = this.props;

    if (!user && !!authToken()) {
      dispatch(fetchCurrentUser())
        .catch(() => {
          return false;
        });
    }

    if (user) {
      dispatch(getConfig());
    }

    return false;
  }

  componentWillReceiveProps (nextProps) {
    const { dispatch, user } = nextProps;

    if (this.props.user !== user) {
      dispatch(getConfig());
    }
  }

  render () {
    const { children, showBackgroundImage } = this.props;

    const wrapperStyles = classnames(
      'wrapper',
      { 'wrapper--background': showBackgroundImage }
    );

    return (
      <div className={wrapperStyles}>
        {children}
        <Footer />
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { showBackgroundImage } = state.app;
  const { user } = state.auth;

  return {
    showBackgroundImage,
    user,
  };
};

export default connect(mapStateToProps)(App);
