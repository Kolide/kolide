import React, { Component, PropTypes } from 'react';
import radium from 'radium';
import componentStyles from './styles';
import avatar from './avatar.svg';
import InputFieldWithIcon from '../fields/InputFieldWithIcon';

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

    const { formData } = this.state;
    const { onSubmit } = this.props;

    return onSubmit(formData);
  }

  render () {
    const { containerStyles, submitButtonStyles, userIconStyles, formStyles } = componentStyles;
    const { onInputChange, onFormSubmit } = this;
    const { formData: { username, password } } = this.state;

    return (
      <form onSubmit={onFormSubmit} style={formStyles}>
        <div style={containerStyles}>
          <img src={avatar} />
          <InputFieldWithIcon
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
        </div>
        <button
          style={submitButtonStyles}
          type="submit"
        >
          Login
        </button>
      </form>
    );
  }
}

export default radium(LoginForm);
