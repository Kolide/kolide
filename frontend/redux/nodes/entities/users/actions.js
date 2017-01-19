import Kolide from 'kolide';

import { formatErrorResponse } from 'redux/nodes/entities/base/helpers';

import config from './config';
const { extendedActions } = config;

// Actions for admin to require password reset for a user
export const REQUIRE_PASSWORD_RESET_REQUEST = 'REQUIRE_PASSWORD_RESET_REQUEST';
export const REQUIRE_PASSWORD_RESET_SUCCESS = 'REQUIRE_PASSWORD_RESET_SUCCESS';
export const REQUIRE_PASSWORD_RESET_FAILURE = 'REQUIRE_PASSWORD_RESET_FAILURE';

export const requirePasswordResetRequest = { type: REQUIRE_PASSWORD_RESET_REQUEST };

export const requirePasswordResetSuccess = (user) => {
  return {
    type: REQUIRE_PASSWORD_RESET_SUCCESS,
    payload: { user },
  };
};

export const requirePasswordResetFailure = (errors) => {
  return {
    type: REQUIRE_PASSWORD_RESET_FAILURE,
    payload: { errors },
  };
};

export const enableUser = (user, { enabled }) => {
  const { successAction, updateFailure, updateRequest, updateSuccess } = extendedActions;

  return (dispatch) => {
    dispatch(updateRequest);

    return Kolide.users.enable(user, { enabled })
      .then((userResponse) => {
        return dispatch(successAction(userResponse, updateSuccess));
      })
      .catch((response) => {
        const errorsObject = formatErrorResponse(response);

        dispatch(updateFailure(errorsObject));

        throw errorsObject;
      });
  };
};

export const requirePasswordReset = (user, { require }) => {
  return (dispatch) => {
    dispatch(requirePasswordResetRequest);

    return Kolide.requirePasswordReset(user, { require })
      .then((updatedUser) => {
        dispatch(requirePasswordResetSuccess(updatedUser));

        return updatedUser;
      })
      .catch((response) => {
        const errorsObject = formatErrorResponse(response);
        dispatch(requirePasswordResetFailure(errorsObject));

        throw response;
      });
  };
};

export const updateAdmin = (user, { admin }) => {
  const { successAction, updateFailure, updateRequest, updateSuccess } = extendedActions;

  return (dispatch) => {
    dispatch(updateRequest);

    return Kolide.users.updateAdmin(user, { admin })
      .then((userResponse) => {
        return dispatch(successAction(userResponse, updateSuccess));
      })
      .catch((response) => {
        const errorsObject = formatErrorResponse(response);

        dispatch(updateFailure(errorsObject));

        throw errorsObject;
      });
  };
};

export default { ...config.actions, enableUser, requirePasswordReset, updateAdmin };
