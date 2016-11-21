import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import InputFieldWithIcon from 'components/forms/fields/InputFieldWithIcon';
import validate from 'components/forms/ResetPasswordForm/validate';

const baseClass = 'reset-password-form';
const formFields = ['new_password', 'new_password_confirmation'];

class ResetPasswordForm extends Component {
  static propTypes = {
    handleSubmit: PropTypes.func,
    fields: PropTypes.shape({
      new_password: formFieldInterface.isRequired,
      new_password_confirmation: formFieldInterface.isRequired,
    }),
  };

  render () {
    const { fields, handleSubmit } = this.props;

    return (
      <form onSubmit={handleSubmit} className={baseClass}>
        <InputFieldWithIcon
          {...fields.new_password}
          autofocus
          iconName="lock"
          placeholder="New Password"
          className={`${baseClass}__input`}
          type="password"
        />
        <InputFieldWithIcon
          {...fields.new_password_confirmation}
          iconName="lock"
          placeholder="Confirm Password"
          className={`${baseClass}__input`}
          type="password"
        />
        <Button
          onClick={handleSubmit}
          className={`${baseClass}__btn`}
          text="Reset Password"
          type="submit"
          variant="gradient"
        />
      </form>
    );
  }
}

export default Form(ResetPasswordForm, {
  fields: formFields,
  validate,
});
