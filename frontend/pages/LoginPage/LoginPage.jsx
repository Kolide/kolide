import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import componentStyles from './styles';
import Icon from '../../components/icons/Icon';
import { loadBackground, resizeBackground } from '../../utilities/backgroundImage';
import local from '../../utilities/local';
import LoginForm from '../../components/forms/LoginForm';
import { loginUser } from '../../redux/nodes/auth/actions';

export class LoginPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    error: PropTypes.string,
    loading: PropTypes.bool,
    user: PropTypes.object,
  };

  componentWillMount () {
    const { dispatch } = this.props;
    const { window } = global;

    if (local.getItem('auth_token')) {
      return dispatch(push('/'));
    }

    loadBackground();
    window.onresize = resizeBackground;
  }

  onSubmit = (formData) => {
    const { dispatch } = this.props;
    dispatch(loginUser(formData))
      .then(() => {
        dispatch(push('/login_successful'));
      });
  }

  render () {
    const { containerStyles, formWrapperStyles, whiteTabStyles } = componentStyles;
    const { onSubmit } = this;

    return (
      <div style={containerStyles}>
        <div style={formWrapperStyles}>
          <Icon name="kolideText" />
          <div style={whiteTabStyles} />
          <LoginForm onSubmit={onSubmit} />
        </div>
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { error, loading, user } = state.auth;

  return {
    error,
    loading,
    user,
  };
};

export default connect(mapStateToProps)(LoginPage);
