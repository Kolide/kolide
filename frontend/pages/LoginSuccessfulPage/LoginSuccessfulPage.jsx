import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import componentStyles from './styles';
import Icon from '../../components/icons/Icon';
import { removeBackground } from '../../utilities/backgroundImage';
import local from '../../utilities/local';

export class LoginSuccessfulPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
  };

  componentWillUnmount () {
    removeBackground();
  }

  render () {
    const { containerStyles, loginSuccessStyles, subtextStyles, whiteBoxStyles } = componentStyles;

    return (
      <div style={containerStyles}>
        <Icon name="kolideText" />
        <div style={whiteBoxStyles}>
          <Icon name="check" />
          <p style={loginSuccessStyles}>Login successful</p>
          <p style={subtextStyles}>Hold on to your butts.</p>
        </div>
      </div>
    );
  }
}

export default connect()(LoginSuccessfulPage);
