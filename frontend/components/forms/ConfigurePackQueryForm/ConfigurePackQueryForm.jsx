import React, { Component, PropTypes } from 'react';

import Icon from 'components/icons/Icon';
import Button from 'components/buttons/Button';
import Dropdown from 'components/forms/fields/Dropdown';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import InputField from 'components/forms/fields/InputField';
import validate from 'components/forms/ConfigurePackQueryForm/validate';

const baseClass = 'configure-pack-query-form';
const fieldNames = ['query_id', 'interval', 'logging_type', 'platform', 'version'];
const platformOptions = [
  { label: 'All', value: 'all' },
  { label: 'Windows', value: 'windows' },
  { label: 'Linux', value: 'linux' },
  { label: 'macOS', value: 'darwin' },
];
const loggingTypeOptions = [
  { label: 'Differential', value: 'differential' },
  { label: 'Differential (Ignore Removals)', value: 'differential_ignore_removals' },
  { label: 'Snapshot', value: 'snapshot' },
];
const minOsqueryVersionOptions = [
  { label: '1.8.1 +', value: '1.8.1' },
  { label: '1.8.2 +', value: '1.8.2' },
  { label: '2.0.0 +', value: '2.0.0' },
  { label: '2.1.1 +', value: '2.1.1' },
  { label: '2.1.2 +', value: '2.1.2' },
  { label: '2.2.0 +', value: '2.2.0' },
  { label: '2.2.1 +', value: '2.2.1' },
];

class ConfigurePackQueryForm extends Component {
  static propTypes = {
    fields: PropTypes.shape({
      interval: formFieldInterface.isRequired,
      logging_type: formFieldInterface.isRequired,
      platform: formFieldInterface.isRequired,
      version: formFieldInterface.isRequired,
    }).isRequired,
    handleSubmit: PropTypes.func,
    onCancel: PropTypes.func,
  };

  onCancel = (evt) => {
    evt.preventDefault();

    const { onCancel: handleCancel } = this.props;

    return handleCancel();
  }

  render () {
    const { fields, handleSubmit } = this.props;
    const { onCancel } = this;

    return (
      <form className={baseClass} onSubmit={handleSubmit}>
        <h2 className={`${baseClass}__title`}>configuration</h2>
        <div className={`${baseClass}__fields`}>
          <InputField
            {...fields.interval}
            inputWrapperClass={`${baseClass}__form-field ${baseClass}__form-field--interval`}
            placeholder="- - -"
            label="Interval"
            hint="Seconds"
            type="number"
          />
          <Dropdown
            {...fields.platform}
            options={platformOptions}
            placeholder="- - -"
            label="Platform"
            wrapperClassName={`${baseClass}__form-field ${baseClass}__form-field--platform`}
          />
          <Dropdown
            {...fields.version}
            options={minOsqueryVersionOptions}
            placeholder="- - -"
            label={['minimum ', <Icon name="osquery" key="min-osquery-vers" />, ' version']}
            wrapperClassName={`${baseClass}__form-field ${baseClass}__form-field--osquer-vers`}
          />
          <Dropdown
            {...fields.logging_type}
            options={loggingTypeOptions}
            placeholder="- - -"
            label="Logging"
            wrapperClassName={`${baseClass}__form-field ${baseClass}__form-field--logging`}
          />
          <div className={`${baseClass}__btn-wrapper`}>
            <Button className={`${baseClass}__cancel-btn`} onClick={onCancel} variant="inverse">
              Cancel
            </Button>
            <Button className={`${baseClass}__submit-btn`} type="submit" variant="brand">
              Save
            </Button>
          </div>
        </div>
      </form>
    );
  }
}

export default Form(ConfigurePackQueryForm, {
  fields: fieldNames,
  validate,
});
