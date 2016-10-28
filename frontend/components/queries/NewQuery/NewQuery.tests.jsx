import React from 'react';
import expect, { createSpy, restoreSpies } from 'expect';
import { mount } from 'enzyme';
import { noop } from 'lodash';

import { createAceSpy, fillInFormInput } from '../../../test/helpers';
import NewQuery from './index';

describe('NewQuery - component', () => {
  beforeEach(() => {
    createAceSpy();
  });
  afterEach(restoreSpies);

  it('renders the ThemeDropdown', () => {
    const component = mount(
      <NewQuery
        onOsqueryTableSelect={noop}
        onTextEditorInputChange={noop}
        textEditorText="Hello world"
      />
    );

    expect(component.find('ThemeDropdown').length).toEqual(1);
  });

  it('does not render the SaveQueryForm by default', () => {
    const component = mount(
      <NewQuery
        onOsqueryTableSelect={noop}
        onTextEditorInputChange={noop}
        textEditorText="Hello world"
      />
    );

    expect(component.find('SaveQueryForm').length).toEqual(0);
  });

  it('renders the SaveQueryFormModal when "Save Query" is clicked', () => {
    const component = mount(
      <NewQuery
        onOsqueryTableSelect={noop}
        onTextEditorInputChange={noop}
        textEditorText="Hello world"
      />
    );

    expect(component.find('SaveQueryForm').length).toEqual(1);
  });

  it('calls onTargetSelectInputChange when changing the select target input text', () => {
    const onTargetSelectInputChangeSpy = createSpy();
    const component = mount(
      <NewQuery onTargetSelectInputChange={onTargetSelectInputChangeSpy} />
    );
    const selectTargetsInput = component.find('.Select-input input');

    fillInFormInput(selectTargetsInput, 'my target');

    expect(onTargetSelectInputChangeSpy).toHaveBeenCalledWith('my target');
  });

  describe('Query string validations', () => {
    const invalidQuery = 'CREATE TABLE users (LastName varchar(255))';
    const validQuery = 'SELECT * FROM users';

    it('calls onInvalidQuerySubmit when invalid', () => {
      const invalidQuerySubmitSpy = createSpy();
      const component = mount(
        <NewQuery
          onInvalidQuerySubmit={invalidQuerySubmitSpy}
          onOsqueryTableSelect={noop}
          onTextEditorInputChange={noop}
          textEditorText={invalidQuery}
        />
      );
      const form = component.find('SaveQueryForm');
      const inputField = form.find('.save-query-form__input--name');

      fillInFormInput(inputField, 'my query');

      form.simulate('submit');

      expect(invalidQuerySubmitSpy).toHaveBeenCalledWith('Cannot INSERT or CREATE in osquery queries');
    });

    it('calls onNewQueryFormSubmit when valid', () => {
      const onNewQueryFormSubmitSpy = createSpy();
      const component = mount(
        <NewQuery
          onNewQueryFormSubmit={onNewQueryFormSubmitSpy}
          onOsqueryTableSelect={noop}
          onTextEditorInputChange={noop}
          textEditorText={validQuery}
        />
      );
      const form = component.find('SaveQueryForm');
      const inputField = form.find('.save-query-form__input--name');

      fillInFormInput(inputField, 'my query');

      form.simulate('submit');

      expect(onNewQueryFormSubmitSpy).toHaveBeenCalled();
    });
  });
});

