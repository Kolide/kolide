import React, { Component, PropTypes } from 'react';
import { Link } from 'react-router';
import radium from 'radium';
import avatar from './avatar.svg';
import componentStyles from './styles';
import GradientButton from '../../buttons/GradientButton';
import InputFieldWithIcon from '../fields/InputFieldWithIcon';
import paths from '../../../router/paths';

class LoginForm extends Component {
  static propTypes = {
    onSubmit: PropTypes.func,
  };

  constructor (props) {
    super(props);

    this.state = {
      formData: {
        username: null,
        password: null,
      },
    };
  }

  onInputChange = (formField) => {
    return ({ target }) => {
      const { formData } = this.state;
      const { value } = target;

      this.setState({
        formData: {
          ...formData,
          [formField]: value,
        },
      });
    };
  }

  onFormSubmit = (evt) => {
    evt.preventDefault();

    if (this.canSubmit()) {
      const { formData } = this.state;
      const { onSubmit } = this.props;

      return onSubmit(formData);
    }

    return false;
  }

  validate = () => {
    const {
      errors,
      formData: { username, password },
    } = this.state;

    if (!validatePresence(username)) {
      this.setState({
        errors: {
          ...errors,
          username: 'Username or email field must be completed',
        },
      });

      return false;
    }

    if (!validatePresence(password)) {
      this.setState({
        errors: {
          ...errors,
          password: 'Password field must be completed',
        },
      });

      return false;
    }

    return true;
  }

  render () {
    const {
      containerStyles,
      forgotPasswordStyles,
      forgotPasswordWrapperStyles,
      formStyles,
      submitButtonStyles,
    } = componentStyles;
    const { onInputChange, onFormSubmit } = this;
    const canSubmit = this.canSubmit();

    return (
      <form onSubmit={onFormSubmit} style={formStyles}>
        <div style={containerStyles}>
          <img alt="Avatar" src={avatar} />
          <InputFieldWithIcon
            autofocus
            iconName="user"
            name="username"
            onChange={onInputChange('username')}
            placeholder="Username or Email"
          />
          <InputFieldWithIcon
            iconName="lock"
            name="password"
            onChange={onInputChange('password')}
            placeholder="Password"
            type="password"
          />
          <div style={forgotPasswordWrapperStyles}>
            <Link style={forgotPasswordStyles} to={paths.FORGOT_PASSWORD}>Forgot Password?</Link>
          </div>
        </div>
        <GradientButton
          onClick={onFormSubmit}
          style={submitButtonStyles}
          text="Login"
          type="submit"
        />
      </form>
    );
  }
}

export default radium(LoginForm);
