import React, { Component, PropTypes } from 'react';
import { size } from 'lodash';

import Button from 'components/buttons/Button';
import DropdownButton from 'components/buttons/DropdownButton';
import Dropdown from 'components/forms/fields/Dropdown';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import helpers from 'components/forms/queries/QueryForm/helpers';
import InputField from 'components/forms/fields/InputField';
import KolideAce from 'components/KolideAce';
import queryInterface from 'interfaces/query';
import SelectTargetsDropdown from 'components/forms/fields/SelectTargetsDropdown';
import targetInterface from 'interfaces/target';
import validateQuery from 'components/forms/validators/validate_query';

const baseClass = 'query-form';

const validate = (formData) => {
  const errors = {};
  const {
    error: queryError,
    valid: queryValid,
  } = validateQuery(formData.query);

  if (!queryValid) {
    errors.query = queryError;
  }

  if (!formData.name) {
    errors.name = 'Title must be present';
  }

  const valid = !size(errors);

  return { valid, errors };
};

class QueryForm extends Component {
  static propTypes = {
    fields: PropTypes.shape({
      description: formFieldInterface.isRequired,
      name: formFieldInterface.isRequired,
      query: formFieldInterface.isRequired,
    }).isRequired,
    handleSubmit: PropTypes.func,
    formData: queryInterface,
    onCancel: PropTypes.func,
    onFetchTargets: PropTypes.func,
    onOsqueryTableSelect: PropTypes.func,
    onRunQuery: PropTypes.func,
    onStopQuery: PropTypes.func,
    onTargetSelect: PropTypes.func,
    onUpdate: PropTypes.func,
    queryIsRunning: PropTypes.bool,
    queryType: PropTypes.string,
    selectedTargets: PropTypes.arrayOf(targetInterface),
    targetsCount: PropTypes.number,
  };

  static defaultProps = {
    queryType: 'query',
    targetsCount: 0,
  };

  constructor (props) {
    super(props);

    this.state = { errors: {} };
  }

  onCancel = (evt) => {
    evt.preventDefault();

    const { onCancel: handleCancel } = this.props;

    return handleCancel();
  }

  onLoad = (editor) => {
    editor.setOptions({
      enableLinking: true,
    });

    editor.on('linkClick', (data) => {
      const { type, value } = data.token;
      const { onOsqueryTableSelect } = this.props;

      if (type === 'osquery-token') {
        return onOsqueryTableSelect(value);
      }

      return false;
    });
  }

  onUpdate = (evt) => {
    evt.preventDefault();

    const { fields } = this.props;
    const { onUpdate: handleUpdate } = this.props;
    const formData = {
      name: fields.name.value,
      query: fields.query.value,
    };

    const { valid, errors } = validate(formData);

    if (valid) {
      handleUpdate(fields);

      return false;
    }

    this.setState({
      errors: {
        ...this.state.errors,
        ...errors,
      },
    });

    return false;
  }

  renderButtons = () => {
    const { canSaveAsNew, canSaveChanges } = helpers;
    const {
      fields,
      formData,
      onRunQuery,
      onStopQuery,
      queryIsRunning,
      queryType,
    } = this.props;
    const { onCancel, onSave, onUpdate } = this;

    const dropdownBtnOptions = [{
      disabled: !canSaveChanges(formData, query),
      label: 'Save Changes',
      onClick: onUpdate,
    }, {
      disabled: !canSaveAsNew(formData, query),
      label: 'Save As New...',
      onClick: onSave,
    }];

    let runQueryButton;

    if (queryIsRunning) {
      runQueryButton = (
        <Button
          className={`${baseClass}__stop-query-btn`}
          onClick={onStopQuery}
          variant="alert"
        >
          Stop Query
        </Button>
      );
    } else {
      runQueryButton = (
        <Button
          className={`${baseClass}__run-query-btn`}
          onClick={onRunQuery}
          variant="brand"
        >
          Run Query
        </Button>
      );
    }

    if (queryType === 'label') {
      return (
        <div className={`${baseClass}__button-wrap`}>
          <Button
            className={`${baseClass}__save-changes-btn`}
            onClick={onCancel}
            variant="inverse"
          >
            Cancel
          </Button>
          <Button
            className={`${baseClass}__save-as-new-btn`}
            disabled={!canSaveAsNew(fields, formData)}
            type="submit"
            variant="brand"
          >
            Save Label
          </Button>
        </div>
      );
    }

    return (
      <div className={`${baseClass}__button-wrap`}>
        <DropdownButton
          className={`${baseClass}__save`}
          options={dropdownBtnOptions}
          variant="success"
        >
          Save
        </DropdownButton>

        {runQueryButton}
      </div>
    );
  }

  renderPlatformDropdown = () => {
    const { fields, queryType } = this.props;

    if (queryType !== 'label') {
      return false;
    }

    const { platformOptions } = helpers;

    return (
      <Dropdown
        {...fields.platform}
        options={platformOptions}
      />
    );
  }

  renderTargetsInput = () => {
    const {
      onFetchTargets,
      onTargetSelect,
      queryType,
      selectedTargets,
      targetsCount,
    } = this.props;

    if (queryType === 'label') {
      return false;
    }


    return (
      <div>
        <SelectTargetsDropdown
          onFetchTargets={onFetchTargets}
          onSelect={onTargetSelect}
          selectedTargets={selectedTargets}
          targetsCount={targetsCount}
          label="Select Targets"
        />
      </div>
    );
  }

  render () {
    const { errors } = this.state;
    const { fields, handleSubmit, queryIsRunning, queryType } = this.props;
    const { onLoad, renderPlatformDropdown, renderButtons, renderTargetsInput } = this;

    return (
      <form className={`${baseClass}__wrapper body-wrap`} onSubmit={handleSubmit}>
        <h1>{queryType === 'label' ? 'New Label Query' : 'New Query'}</h1>
        <KolideAce
          {...fields.query}
          error={fields.query.error || errors.query}
          onLoad={onLoad}
          readOnly={queryIsRunning}
          wrapperClassName={`${baseClass}__text-editor-wrapper`}
        />
        {renderTargetsInput()}
        <InputField
          {...fields.name}
          error={fields.name.error || errors.name}
          inputClassName={`${baseClass}__query-title`}
          label={queryType === 'label' ? 'Label title' : 'Query Title'}
        />
        <InputField
          {...fields.description}
          inputClassName={`${baseClass}__query-description`}
          label="Description"
          type="textarea"
        />
        {renderPlatformDropdown()}
        {renderButtons()}
      </form>
    );
  }
}

export default Form(QueryForm, {
  fields: ['description', 'name', 'platform', 'query'],
  validate,
});
