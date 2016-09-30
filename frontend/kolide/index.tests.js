import expect from 'expect';
import Kolide from './index';
import mocks from '../test/mocks';

const {
  invalidForgotPasswordRequest,
  invalidResetPasswordRequest,
  validForgotPasswordRequest,
  validGetUsersRequest,
  validLoginRequest,
  validLogoutRequest,
  validMeRequest,
  validResetPasswordRequest,
  validUser,
} = mocks;

describe('Kolide - API client', () => {
  describe('defaults', () => {
    it('sets the base URL', () => {
      expect(Kolide.baseURL).toEqual('http://localhost:8080/api');
    });
  });

  describe('#getUsers', () => {
    it('calls the appropriate endpoint with the correct parameters', (done) => {
      const bearerToken = 'valid-bearer-token';
      const request = validGetUsersRequest();

      Kolide.getUsers(bearerToken)
        .then((users) => {
          expect(users).toEqual([validUser]);
          expect(request.isDone()).toEqual(true);
          done();
        })
        .catch(done);
    });
  });

  describe('#me', () => {
    it('calls the appropriate endpoint with the correct parameters', (done) => {
      const bearerToken = 'ABC123';
      const request = validMeRequest(bearerToken);

      Kolide.setBearerToken(bearerToken);
      Kolide.me()
        .then((user) => {
          expect(user).toEqual(validUser);
          expect(request.isDone()).toEqual(true);
          done();
        })
        .catch(done);
    });
  });

  describe('#loginUser', () => {
    it('calls the appropriate endpoint with the correct parameters', (done) => {
      const request = validLoginRequest();

      Kolide.loginUser({
        username: 'admin',
        password: 'secret',
      })
        .then((user) => {
          expect(user).toEqual(validUser);
          expect(request.isDone()).toEqual(true);
          done();
        })
        .catch(done);
    });
  });

  describe('#logout', () => {
    it('calls the appropriate endpoint with the correct parameters', (done) => {
      const bearerToken = 'ABC123';
      const request = validLogoutRequest(bearerToken);

      Kolide.setBearerToken(bearerToken);
      Kolide.logout()
        .then(() => {
          expect(request.isDone()).toEqual(true);
          done();
        })
        .catch(done);
    });
  });

  describe('#forgotPassword', () => {
    it('calls the appropriate endpoint with the correct parameters when successful', (done) => {
      const request = validForgotPasswordRequest();
      const email = 'hi@thegnar.co';

      Kolide.forgotPassword({ email })
        .then(() => {
          expect(request.isDone()).toEqual(true);
          done();
        })
        .catch(done);
    });

    it('return errors correctly for unsuccessful requests', (done) => {
      const error = 'Something went wrong';
      const request = invalidForgotPasswordRequest(error);
      const email = 'hi@thegnar.co';

      Kolide.forgotPassword({ email })
        .then(done)
        .catch(errorResponse => {
          const { response } = errorResponse;

          expect(response).toEqual({ error });
          expect(request.isDone()).toEqual(true);
          done();
        });
    });
  });

  describe('#resetPassword', () => {
    const newPassword = 'p@ssw0rd';

    it('calls the appropriate endpoint with the correct parameters when successful', (done) => {
      const passwordResetToken = 'password-reset-token';
      const request = validResetPasswordRequest(newPassword, passwordResetToken);
      const formData = {
        new_password: newPassword,
        password_reset_token: passwordResetToken,
      };

      Kolide.resetPassword(formData)
        .then(() => {
          expect(request.isDone()).toEqual(true);
          done();
        })
        .catch(done);
    });

    it('return errors correctly for unsuccessful requests', (done) => {
      const error = 'Resource not found';
      const passwordResetToken = 'invalid-password-reset-token';
      const request = invalidResetPasswordRequest(newPassword, passwordResetToken, error);
      const formData = {
        new_password: newPassword,
        password_reset_token: passwordResetToken,
      };

      Kolide.resetPassword(formData)
        .then(done)
        .catch(errorResponse => {
          const { response } = errorResponse;

          expect(response).toEqual({ error });
          expect(request.isDone()).toEqual(true);
          done();
        });
    });
  });
});
