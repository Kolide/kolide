import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';
import { pick } from 'lodash';

import FormField from 'components/forms/FormField';

const baseClass = 'input-field';

class InputField extends Component {
  static propTypes = {
    autofocus: PropTypes.bool,
    error: PropTypes.string,
    inputClassName: PropTypes.string, // eslint-disable-line react/forbid-prop-types
    inputWrapperClass: PropTypes.string,
    inputOptions: PropTypes.object, // eslint-disable-line react/forbid-prop-types
    name: PropTypes.string,
    onChange: PropTypes.func,
    placeholder: PropTypes.string,
    type: PropTypes.string,
    value: PropTypes.string.isRequired,
  };

  static defaultProps = {
    autofocus: false,
    inputWrapperClass: '',
    inputOptions: {},
    label: null,
    labelClassName: '',
    type: 'text',
    value: '',
  };

  componentDidMount () {
    const { autofocus } = this.props;
    const { input } = this;

    if (autofocus) {
      input.focus();
    }

    return false;
  }

  onInputChange = (evt) => {
    evt.preventDefault();

    const { value } = evt.target;
    const { onChange } = this.props;

    return onChange(value);
  }

  render () {
    const { error, inputClassName, inputOptions, inputWrapperClass, name, placeholder, type, value } = this.props;
    const { onInputChange } = this;
    const shouldShowPasswordClass = type === 'password';
    const inputClasses = classnames(baseClass, inputClassName, {
      [`${baseClass}--password`]: shouldShowPasswordClass,
      [`${baseClass}--error`]: error,
      [`${baseClass}__textarea`]: type === 'textarea',
    });

    const formFieldProps = pick(this.props, ['hint', 'label', 'error', 'name']);

    if (type === 'textarea') {
      return (
        <FormField {...formFieldProps} type="textarea" className={inputWrapperClass}>
          <textarea
            name={name}
            onChange={onInputChange}
            className={inputClasses}
            placeholder={placeholder}
            ref={(r) => { this.input = r; }}
            type={type}
            {...inputOptions}
            value={value}
          />
        </FormField>
      );
    }

    return (
      <FormField {...formFieldProps} type="input" className={inputWrapperClass}>
        <input
          name={name}
          onChange={onInputChange}
          className={inputClasses}
          placeholder={placeholder}
          ref={(r) => { this.input = r; }}
          type={type}
          {...inputOptions}
          value={value}
        />
      </FormField>
    );
  }
}

export default InputField;

