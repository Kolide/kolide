import expect from 'expect';

import helpers from 'kolide/helpers';

describe('Kolide API - helpers', () => {
  describe('#setupData', () => {
    const formData = {
      email: 'hi@gnar.dog',
      full_name: 'Gnar Dog',
      kolide_server_url: 'https://gnar.kolide.co',
      org_logo_url: 'https://thegnar.co/assets/logo.png',
      org_name: 'The Gnar Co.',
      password: 'p@ssw0rd',
      password_confirmation: 'p@ssw0rd',
      username: 'gnardog',
    };

    it('formats the form data to send to the server', () => {
      expect(helpers.setupData(formData)).toEqual({
        kolide_server_url: 'https://gnar.kolide.co',
        org_info: {
          org_logo_url: 'https://thegnar.co/assets/logo.png',
          org_name: 'The Gnar Co.',
        },
        admin: {
          admin: true,
          email: 'hi@gnar.dog',
          full_name: 'Gnar Dog',
          password: 'p@ssw0rd',
          password_confirmation: 'p@ssw0rd',
          username: 'gnardog',
        },
      });
    });
  });
});
