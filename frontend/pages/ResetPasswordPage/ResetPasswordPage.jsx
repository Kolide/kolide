import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { noop } from 'lodash';
import { push } from 'react-router-redux';

import debounce from '../../utilities/debounce';
import { resetPassword } from '../../redux/nodes/components/ResetPasswordPage/actions';
import ResetPasswordForm from '../../components/forms/ResetPasswordForm';
import StackedWhiteBoxes from '../../components/StackedWhiteBoxes';
import { updateUser, performRequiredPasswordReset } from '../../redux/nodes/auth/actions';
import userInterface from '../../interfaces/user';

export class ResetPasswordPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    token: PropTypes.string,
    user: userInterface,
  };

  static defaultProps = {
    dispatch: noop,
  };

  componentWillMount () {
    const { dispatch, token, user } = this.props;

    if (!user && !token) {
      return dispatch(push('/login'));
    }

    return false;
  }

  onSubmit = debounce((formData) => {
    const { dispatch, token, user } = this.props;

    if (user) {
      return this.loggedInUser(formData);
    }

    const resetPasswordData = {
      ...formData,
      password_reset_token: token,
    };

    return dispatch(resetPassword(resetPasswordData))
      .then(() => {
        return dispatch(push('/login'));
      });
  })

  handleLeave = (location) => {
    const { dispatch } = this.props;

    return dispatch(push(location));
  }

  loggedInUser = (formData) => {
    const { dispatch, user } = this.props;
    const { new_password: password } = formData;
    const passwordUpdateParams = { password };

    return dispatch(performRequiredPasswordReset(passwordUpdateParams))
      .then(() => { return dispatch(push('/')); });
  }

  render () {
    const { handleLeave, onSubmit } = this;

    return (
      <StackedWhiteBoxes
        headerText="Reset Password"
        leadText="Create a new password using at least one letter, one numeral and seven characters."
        onLeave={handleLeave}
      >
        <ResetPasswordForm handleSubmit={onSubmit} />
      </StackedWhiteBoxes>
    );
  }
}

const mapStateToProps = (state) => {
  const { ResetPasswordPage: componentState } = state.components;
  const { user } = state.auth;

  return {
    ...componentState,
    user,
  };
};

export default connect(mapStateToProps)(ResetPasswordPage);
