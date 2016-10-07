import md5 from 'js-md5';
import Kolide from '../../../kolide';

export const CLEAR_AUTH_ERRORS = 'CLEAR_AUTH_ERRORS';
export const LOGIN_REQUEST = 'LOGIN_REQUEST';
export const LOGIN_SUCCESS = 'LOGIN_SUCCESS';
export const LOGIN_FAILURE = 'LOGIN_FAILURE';
export const LOGOUT_REQUEST = 'LOGOUT_REQUEST';
export const LOGOUT_SUCCESS = 'LOGOUT_SUCCESS';
export const LOGOUT_FAILURE = 'LOGOUT_FAILURE';

export const clearAuthErrors = { type: CLEAR_AUTH_ERRORS };
export const loginRequest = { type: LOGIN_REQUEST };
export const loginSuccess = (user) => {
  return {
    type: LOGIN_SUCCESS,
    payload: {
      data: user,
    },
  };
};
export const loginFailure = (error) => {
  return {
    type: LOGIN_FAILURE,
    payload: {
      error,
    },
  };
};

export const fetchCurrentUser = () => {
  return (dispatch) => {
    dispatch(loginRequest);
    return Kolide.me()
      .then(response => {
        const { user } = response;
        const { email } = user;
        const emailHash = md5(email.toLowerCase());

        user.gravatarURL = `https://www.gravatar.com/avatar/${emailHash}`;
        return dispatch(loginSuccess(user));
      })
      .catch(response => {
        dispatch(loginFailure('Unable to authenticate the current user'));
        throw response;
      });
  };
};

// formData should be { username: <string>, password: <string> }
export const loginUser = (formData) => {
  return (dispatch) => {
    return new Promise((resolve, reject) => {
      dispatch(loginRequest);
      return Kolide.loginUser(formData)
        .then(response => {
          const { user } = response;
          const { email } = user;
          const emailHash = md5(email.toLowerCase());

          user.gravatarURL = `https://www.gravatar.com/avatar/${emailHash}`;
          dispatch(loginSuccess(response));
          return resolve(user);
        })
        .catch(response => {
          const { error } = response;
          dispatch(loginFailure(error));
          return reject(error);
        });
    });
  };
};

export const logoutFailure = (error) => {
  return {
    type: LOGOUT_FAILURE,
    error,
  };
};
export const logoutRequest = { type: LOGOUT_REQUEST };
export const logoutSuccess = { type: LOGOUT_SUCCESS };
export const logoutUser = () => {
  return (dispatch) => {
    dispatch(logoutRequest);
    return Kolide.logout()
      .then(() => {
        return dispatch(logoutSuccess);
      })
      .catch((error) => {
        dispatch(logoutFailure('Unable to log out of your account'));
        throw error;
      });
  };
};
