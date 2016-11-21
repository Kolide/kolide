import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import InputField from 'components/forms/fields/InputField';

const baseClass = 'edit-user-form';
const fieldNames = ['email', 'name', 'position', 'username'];

class EditUserForm extends Component {
  static propTypes = {
    onCancel: PropTypes.func,
    handleSubmit: PropTypes.func,
    fields: PropTypes.shape({
      email: formFieldInterface.isRequired,
      name: formFieldInterface.isRequired,
      position: formFieldInterface.isRequired,
      username: formFieldInterface.isRequired,
    }).isRequired,
  };

  render () {
    const { fields, handleSubmit, onCancel } = this.props;

    return (
      <form className={baseClass} onSubmit={handleSubmit}>
        <InputField
          {...fields.name}
          label="Name"
          labelClassName={`${baseClass}__label`}
          inputWrapperClass={`${baseClass}__input-wrap ${baseClass}__input-wrap--first`}
          inputClassName={`${baseClass}__input`}
        />
        <InputField
          {...fields.username}
          label="Username"
          labelClassName={`${baseClass}__label`}
          inputWrapperClass={`${baseClass}__input-wrap`}
          inputClassName={`${baseClass}__input ${baseClass}__input--username`}
        />
        <InputField
          {...fields.position}
          label="Position"
          labelClassName={`${baseClass}__label`}
          inputWrapperClass={`${baseClass}__input-wrap`}
          inputClassName={`${baseClass}__input`}
        />
        <InputField
          {...fields.email}
          inputWrapperClass={`${baseClass}__input-wrap`}
          label="Email"
          labelClassName={`${baseClass}__label`}
          inputClassName={`${baseClass}__input ${baseClass}__input--email`}
        />
        <div className={`${baseClass}__btn-wrap`}>
          <Button
            className={`${baseClass}__form-btn ${baseClass}__form-btn--submit`}
            text="Submit"
            type="submit"
            variant="brand"
          />
          <Button
            className={`${baseClass}__form-btn`}
            onClick={onCancel}
            text="Cancel"
            variant="inverse"
          />
        </div>
      </form>
    );
  }
}

export default Form(EditUserForm, {
  fields: fieldNames,
});
