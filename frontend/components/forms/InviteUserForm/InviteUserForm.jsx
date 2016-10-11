import React, { Component, PropTypes } from 'react';
import radium from 'radium';
import componentStyles from './styles';
import Button from '../../buttons/Button';
import InputFieldWithIcon from '../fields/InputFieldWithIcon';
import validatePresence from '../validators/validate_presence';
import validEmail from '../validators/valid_email';

class InviteUserForm extends Component {
  static propTypes = {
    error: PropTypes.string,
    invitedBy: PropTypes.shape({
      id: PropTypes.number,
    }),
    onCancel: PropTypes.func,
    onSubmit: PropTypes.func,
  };

  constructor (props) {
    super(props);

    this.state = {
      errors: {
        email: null,
        admin: null,
      },
      formData: {
        email: null,
        admin: 'false',
      },
    };
  }

  componentWillReceiveProps (nextProps) {
    const { error } = nextProps;
    const { errors } = this.state;

    if (this.props.error !== error) {
      this.setState({
        errors: {
          ...errors,
          email: error,
        },
      });
    }
  }

  onInputChange = (formField) => {
    return ({ target }) => {
      const { errors, formData } = this.state;
      const { value } = target;

      this.setState({
        errors: {
          ...errors,
          [formField]: null,
        },
        formData: {
          ...formData,
          [formField]: value,
        },
      });
    };
  }

  onFormSubmit = (evt) => {
    evt.preventDefault();
    const valid = this.validate();

    if (valid) {
      const { formData: { admin, email } } = this.state;
      const { invitedBy, onSubmit } = this.props;

      return onSubmit({
        admin: admin === 'true',
        email,
        invited_by: invitedBy.id,
      });
    }

    return false;
  }

  validate = () => {
    const {
      errors,
      formData: { email },
    } = this.state;

    if (!validatePresence(email)) {
      this.setState({
        errors: {
          ...errors,
          email: 'Email field must be completed',
        },
      });

      return false;
    }

    if (!validEmail(email)) {
      this.setState({
        errors: {
          ...errors,
          email: `${email} is not a valid email`,
        },
      });

      return false;
    }

    return true;
  }

  render () {
    const { buttonStyles, buttonWrapperStyles, radioElementStyles, roleTitleStyles } = componentStyles;
    const { errors, formData: { admin } } = this.state;
    const { onCancel } = this.props;
    const { onFormSubmit, onInputChange } = this;

    return (
      <form onSubmit={onFormSubmit}>
        <InputFieldWithIcon
          autofocus
          error={errors.email}
          name="email"
          onChange={onInputChange('email')}
          placeholder="Email"
        />
        <div style={radioElementStyles}>
          <p style={roleTitleStyles}>role</p>
          <input
            checked={admin === 'false'}
            onChange={onInputChange('admin')}
            type="radio"
            value="false"
          /> USER (default)
          <br />
          <input
            checked={admin === 'true'}
            onChange={onInputChange('admin')}
            type="radio"
            value="true"
          /> ADMIN
        </div>
        <div style={buttonWrapperStyles}>
          <Button
            onClick={onCancel}
            style={buttonStyles}
            text="Cancel"
            variant="inverse"
          />
          <Button
            style={buttonStyles}
            text="Invite"
            type="submit"
          />
        </div>
      </form>
    );
  }
}

export default radium(InviteUserForm);
