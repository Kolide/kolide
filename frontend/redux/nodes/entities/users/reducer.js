import {
  REQUIRE_PASSWORD_RESET_FAILURE,
  REQUIRE_PASSWORD_RESET_REQUEST,
  REQUIRE_PASSWORD_RESET_SUCCESS,
  PERFORM_REQUIRED_PASSWORD_RESET_FAILURE,
  PERFORM_REQUIRED_PASSWORD_RESET_REQUEST,
  PERFORM_REQUIRED_PASSWORD_RESET_SUCCESS,
} from './actions';
import config, { initialState } from './config';

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case REQUIRE_PASSWORD_RESET_REQUEST:
      return {
        ...state,
        errors: {},
        loading: true,
      };
    case REQUIRE_PASSWORD_RESET_SUCCESS:
      return {
        ...state,
        errors: {},
        loading: false,
        data: {
          ...state.data,
          [payload.user.id]: payload.user,
        },
      };
    case REQUIRE_PASSWORD_RESET_FAILURE:
      return {
        ...state,
        loading: false,
        errors: payload.errors,
      };
    default:
      return config.reducer(state, { type, payload });
  }
};
