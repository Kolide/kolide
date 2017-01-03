import React from 'react';
import expect, { createSpy, restoreSpies } from 'expect';
import { mount } from 'enzyme';

import ConfirmInviteForm from 'components/forms/ConfirmInviteForm';
import { fillInFormInput } from 'test/helpers';

describe('ConfirmInviteForm - component', () => {
  afterEach(restoreSpies);

  const handleSubmitSpy = createSpy();
  const inviteToken = 'abc123';
  const formData = { invite_token: inviteToken };
  const form = mount(<ConfirmInviteForm formData={formData} handleSubmit={handleSubmitSpy} />);

  const nameInput = form.find({ name: 'name' }).find('input');
  const passwordConfirmationInput = form.find({ name: 'password_confirmation' }).find('input');
  const passwordInput = form.find({ name: 'password' }).find('input');
  const submitBtn = form.find('button');
  const usernameInput = form.find({ name: 'username' }).find('input');

  it('renders', () => {
    expect(form.length).toEqual(1);
  });

  it('calls the handleSubmit prop with the invite_token when valid', () => {
    fillInFormInput(nameInput, 'Gnar Dog');
    fillInFormInput(usernameInput, 'gnardog');
    fillInFormInput(passwordInput, 'p@ssw0rd');
    fillInFormInput(passwordConfirmationInput, 'p@ssw0rd');
    submitBtn.simulate('click');

    expect(handleSubmitSpy).toHaveBeenCalledWith({
      ...formData,
      name: 'Gnar Dog',
      username: 'gnardog',
      password: 'p@ssw0rd',
      password_confirmation: 'p@ssw0rd',
    });
  });

  describe('name input', () => {
    it('changes form state on change', () => {
      fillInFormInput(nameInput, 'Gnar Dog');

      expect(form.state().formData).toInclude({ name: 'Gnar Dog' });
    });

    it('validates the field must be present', () => {
      fillInFormInput(nameInput, '');
      form.find('button').simulate('click');

      expect(form.state().errors).toInclude({ name: 'Full name must be present' });
    });
  });

  describe('username input', () => {
    it('changes form state on change', () => {
      fillInFormInput(usernameInput, 'gnardog');

      expect(form.state().formData).toInclude({ username: 'gnardog' });
    });

    it('validates the field must be present', () => {
      fillInFormInput(usernameInput, '');
      submitBtn.simulate('click');

      expect(form.state().errors).toInclude({ username: 'Username must be present' });
    });
  });

  describe('password input', () => {
    it('changes form state on change', () => {
      fillInFormInput(passwordInput, 'p@ssw0rd');

      expect(form.state().formData).toInclude({ password: 'p@ssw0rd' });
    });

    it('validates the field must be present', () => {
      fillInFormInput(passwordInput, '');
      form.find('button').simulate('click');

      expect(form.state().errors).toInclude({ password: 'Password must be present' });
    });
  });

  describe('password_confirmation input', () => {
    it('changes form state on change', () => {
      fillInFormInput(passwordConfirmationInput, 'p@ssw0rd');

      expect(form.state().formData).toInclude({ password_confirmation: 'p@ssw0rd' });
    });

    it('validates the password_confirmation matches the password', () => {
      fillInFormInput(passwordInput, 'p@ssw0rd');
      fillInFormInput(passwordConfirmationInput, 'another-password');
      form.find('button').simulate('click');

      expect(form.state().errors).toInclude({
        password_confirmation: 'Password confirmation does not match password',
      });
    });

    it('validates the field must be present', () => {
      fillInFormInput(passwordConfirmationInput, '');
      form.find('button').simulate('click');

      expect(form.state().errors).toInclude({ password_confirmation: 'Password confirmation must be present' });
    });
  });
});

