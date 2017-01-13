import React from 'react';
import expect, { createSpy, restoreSpies } from 'expect';
import { mount } from 'enzyme';

import ConfigOptionsForm from 'components/forms/ConfigOptionsForm';
import { configOptionStub } from 'test/stubs';
import { fillInFormInput } from 'test/helpers';

describe('ConfigOptionsForm - form', () => {
  afterEach(restoreSpies);

  it('renders a ConfigOptionForm for each completed config option', () => {
    const formWithOneOption = mount(<ConfigOptionsForm configNameOptions={[]} completedOptions={[configOptionStub]} />);
    const formWithTwoOptions = mount(<ConfigOptionsForm configNameOptions={[]} completedOptions={[configOptionStub, configOptionStub]} />);

    expect(formWithOneOption.find('ConfigOptionForm').length).toEqual(1);
    expect(formWithTwoOptions.find('ConfigOptionForm').length).toEqual(2);
  });

  it('calls the onFormUpdate prop with the old and new option when the option is updated', () => {
    const spy = createSpy();
    const form = mount(<ConfigOptionsForm configNameOptions={[]} completedOptions={[configOptionStub]} onFormUpdate={spy} />);
    const configOptionFormInput = form.find('ConfigOptionForm').find('InputField');

    fillInFormInput(configOptionFormInput.find('input'), 'updated value');

    expect(spy).toHaveBeenCalledWith(configOptionStub, { ...configOptionStub, value: 'updated value' });
  });
});
