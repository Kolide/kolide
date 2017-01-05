import expect, { createSpy, restoreSpies } from 'expect';
import { find } from 'lodash';

import { formatErrorResponse } from 'redux/nodes/entities/base/helpers';
import reduxConfig from 'redux/nodes/entities/base/reduxConfig';
import { reduxMockStore } from 'test/helpers';
import schemas from 'redux/nodes/entities/base/schemas';

const store = {
  entities: {
    invites: {
      errors: {},
      data: {},
      loading: false,
    },
    users: {
      errors: {},
      data: {},
      loading: false,
    },
  },
};
const invite = { id: 1, name: 'Gnar Dog', email: 'hi@thegnar.co' };
const user = { id: 1, email: 'hi@thegnar.co' };

describe('reduxConfig', () => {
  afterEach(restoreSpies);

  describe('dispatching the create action', () => {
    describe('successful create call', () => {
      const mockStore = reduxMockStore(store);
      const createFunc = createSpy().andCall(() => {
        return Promise.resolve([user]);
      });

      const config = reduxConfig({
        createFunc,
        entityName: 'users',
        schema: schemas.USERS,
      });
      const { actions, reducer } = config;

      it('calls the createFunc', () => {
        mockStore.dispatch(actions.create());

        expect(createFunc).toHaveBeenCalled();
      });

      it('dispatches the correct actions', () => {
        mockStore.dispatch(actions.create());

        const dispatchedActions = mockStore.getActions();
        const dispatchedActionTypes = dispatchedActions.map((action) => { return action.type; });

        expect(dispatchedActionTypes).toInclude('users_CREATE_REQUEST');
        expect(dispatchedActionTypes).toInclude('users_CREATE_SUCCESS');
        expect(dispatchedActionTypes).toNotInclude('users_CREATE_FAILURE');
      });

      it('adds the returned user to state', () => {
        const createSuccessAction = {
          type: 'users_CREATE_SUCCESS',
          payload: {
            data: {
              users: {
                [user.id]: user,
              },
            },
          },
        };
        const initialState = {
          loading: false,
          entities: {},
          errors: {},
        };
        const newState = reducer(initialState, createSuccessAction);

        expect(newState.data[user.id]).toEqual(user);
      });
    });

    describe('unsuccessful create call', () => {
      const mockStore = reduxMockStore(store);
      const errors = [
        { name: 'first_name',
          reason: 'is not valid',
        },
        { name: 'last_name',
          reason: 'must be changed or something',
        },
      ];
      const errorResponse = {
        message: {
          message: 'Validation Failed',
          errors,
        },
      };
      const formattedErrors = formatErrorResponse(errorResponse);
      const createFunc = createSpy().andCall(() => {
        return Promise.reject(errorResponse);
      });
      const config = reduxConfig({
        createFunc,
        entityName: 'users',
        schema: schemas.USERS,
      });
      const { actions, reducer } = config;

      it('calls the createFunc', () => {
        mockStore.dispatch(actions.create());

        expect(createFunc).toHaveBeenCalled();
      });

      it('dispatches the correct actions', () => {
        mockStore.dispatch(actions.create());

        const dispatchedActions = mockStore.getActions();
        const dispatchedActionTypes = dispatchedActions.map((action) => { return action.type; });

        expect(dispatchedActionTypes).toInclude('users_CREATE_REQUEST');
        expect(dispatchedActionTypes).toNotInclude('users_CREATE_SUCCESS');

        const createFailureAction = find(dispatchedActions, { type: 'users_CREATE_FAILURE' });

        expect(createFailureAction.payload).toEqual({
          errors: formattedErrors,
        });
      });

      it('adds the returned errors to state', () => {
        const createFailureAction = {
          type: 'users_CREATE_FAILURE',
          payload: {
            errors: formattedErrors,
          },
        };
        const initialState = {
          loading: false,
          entities: {},
          errors: {},
        };
        const newState = reducer(initialState, createFailureAction);

        expect(newState.errors).toEqual(formattedErrors);
      });
    });
  });

  describe('dispatching the destroy action', () => {
    describe('successful destroy call', () => {
      const mockStore = reduxMockStore(store);
      const destroyFunc = createSpy().andCall(() => {
        return Promise.resolve();
      });

      const config = reduxConfig({
        destroyFunc,
        entityName: 'invites',
        schema: schemas.INVITES,
      });
      const { actions, reducer } = config;

      it('calls the destroyFunc', () => {
        mockStore.dispatch(actions.destroy({ inviteID: invite.id }));

        expect(destroyFunc).toHaveBeenCalled();
      });

      it('dispatches the correct actions', () => {
        mockStore.dispatch(actions.destroy({ inviteID: invite.id }));

        const dispatchedActions = mockStore.getActions();
        const dispatchedActionTypes = dispatchedActions.map((action) => { return action.type; });

        expect(dispatchedActionTypes).toInclude('invites_DESTROY_REQUEST');
        expect(dispatchedActionTypes).toInclude('invites_DESTROY_SUCCESS');
        expect(dispatchedActionTypes).toNotInclude('invites_DESTROY_FAILURE');
      });

      it('removes the returned invite from state', () => {
        const destroySuccessAction = {
          type: 'invites_DESTROY_SUCCESS',
          payload: {
            id: 1,
          },
        };
        const initialState = {
          data: {
            [invite.id]: invite,
            2: { id: 2, name: 'Jason Meller' },
          },
          errors: {},
          loading: false,
        };
        const newState = reducer(initialState, destroySuccessAction);

        expect(newState.data).toEqual({
          2: { id: 2, name: 'Jason Meller' },
        });
      });
    });

    describe('unsuccessful destroy call', () => {
      const mockStore = reduxMockStore(store);
      const errors = [
        {
          name: 'base',
          reason: 'Unable to create user',
        },
      ];
      const errorResponse = {
        message: {
          message: 'Validation Failed',
          errors,
        },
      };
      const destroyFunc = createSpy().andCall(() => {
        return Promise.reject(errorResponse);
      });
      const formattedErrors = formatErrorResponse(errorResponse);
      const config = reduxConfig({
        destroyFunc,
        entityName: 'users',
        schema: schemas.USERS,
      });
      const { actions, reducer } = config;

      it('calls the createFunc', () => {
        mockStore.dispatch(actions.destroy());

        expect(destroyFunc).toHaveBeenCalled();
      });

      it('dispatches the correct actions', () => {
        mockStore.dispatch(actions.destroy());

        const dispatchedActions = mockStore.getActions();
        const dispatchedActionTypes = dispatchedActions.map((action) => { return action.type; });
        const destroyFailureAction = find(dispatchedActions, { type: 'users_DESTROY_FAILURE' });

        expect(dispatchedActionTypes).toInclude('users_DESTROY_REQUEST');
        expect(dispatchedActionTypes).toNotInclude('users_DESTROY_SUCCESS');

        expect(destroyFailureAction.payload).toEqual({
          errors: formattedErrors,
        });
      });

      it('adds the returned errors to state', () => {
        const destroyFailureAction = {
          type: 'users_DESTROY_FAILURE',
          payload: {
            errors: formattedErrors,
          },
        };
        const initialState = {
          loading: false,
          entities: {},
          errors: {},
        };
        const newState = reducer(initialState, destroyFailureAction);

        expect(newState.errors).toEqual(formattedErrors);
      });
    });
  });

  describe('dispatching the load action', () => {
    describe('successful load call', () => {
      const mockStore = reduxMockStore(store);
      const loadFunc = createSpy().andCall(() => {
        return Promise.resolve(user);
      });

      const config = reduxConfig({
        entityName: 'users',
        loadFunc,
        schema: schemas.USERS,
      });
      const { actions, reducer } = config;

      it('calls the loadFunc', () => {
        mockStore.dispatch(actions.load());

        expect(loadFunc).toHaveBeenCalled();
      });

      it('dispatches the correct actions', () => {
        mockStore.dispatch(actions.load());

        const dispatchedActions = mockStore.getActions();
        const dispatchedActionTypes = dispatchedActions.map((action) => { return action.type; });

        expect(dispatchedActionTypes).toInclude('users_LOAD_REQUEST');
        expect(dispatchedActionTypes).toInclude('users_LOAD_SUCCESS');
        expect(dispatchedActionTypes).toNotInclude('users_LOAD_FAILURE');
      });

      it('adds the returned user to state', () => {
        const loadSuccessAction = {
          type: 'users_LOAD_SUCCESS',
          payload: {
            data: {
              users: {
                [user.id]: user,
              },
            },
          },
        };
        const initialState = {
          loading: false,
          entities: {},
          errors: {},
        };
        const newState = reducer(initialState, loadSuccessAction);

        expect(newState.data[user.id]).toEqual(user);
      });
    });

    describe('unsuccessful load call', () => {
      const mockStore = reduxMockStore(store);
      const errors = [
        {
          name: 'base',
          reason: 'Unable to load users',
        },
      ];
      const errorResponse = {
        message: {
          message: 'Something went wrong',
          errors,
        },
      };
      const formattedErrors = formatErrorResponse(errorResponse);
      const loadFunc = createSpy().andCall(() => {
        return Promise.reject(errorResponse);
      });
      const config = reduxConfig({
        entityName: 'users',
        loadFunc,
        schema: schemas.USERS,
      });
      const { actions, reducer } = config;

      it('calls the loadFunc', () => {
        mockStore.dispatch(actions.load());

        expect(loadFunc).toHaveBeenCalled();
      });

      it('dispatches the correct actions', () => {
        mockStore.dispatch(actions.load());

        const dispatchedActions = mockStore.getActions();
        const dispatchedActionTypes = dispatchedActions.map((action) => { return action.type; });
        const loadFailureAction = find(dispatchedActions, { type: 'users_LOAD_FAILURE' });

        expect(dispatchedActionTypes).toInclude('users_LOAD_REQUEST');
        expect(dispatchedActionTypes).toNotInclude('users_LOAD_SUCCESS');
        expect(loadFailureAction.payload).toEqual({
          errors: formattedErrors,
        });
      });

      it('adds the returned errors to state', () => {
        const loadFailureAction = {
          type: 'users_LOAD_FAILURE',
          payload: {
            errors: formattedErrors,
          },
        };
        const initialState = {
          loading: false,
          entities: {},
          errors: {},
        };
        const newState = reducer(initialState, loadFailureAction);

        expect(newState.errors).toEqual(formattedErrors);
      });
    });
  });

  describe('dispatching the loadAll action', () => {
    describe('successful load call', () => {
      const mockStore = reduxMockStore(store);
      const loadAllFunc = createSpy().andCall(() => {
        return Promise.resolve([user]);
      });

      const config = reduxConfig({
        entityName: 'users',
        loadAllFunc,
        schema: schemas.USERS,
      });
      const { actions, reducer } = config;

      it('calls the loadAllFunc', () => {
        mockStore.dispatch(actions.loadAll());

        expect(loadAllFunc).toHaveBeenCalled();
      });

      it('dispatches the correct actions', () => {
        mockStore.dispatch(actions.loadAll());

        const dispatchedActions = mockStore.getActions();
        const dispatchedActionTypes = dispatchedActions.map((action) => { return action.type; });

        expect(dispatchedActionTypes).toInclude('users_LOAD_REQUEST');
        expect(dispatchedActionTypes).toInclude('users_LOAD_SUCCESS');
        expect(dispatchedActionTypes).toNotInclude('users_LOAD_FAILURE');
      });

      it('adds the returned user to state', () => {
        const loadSuccessAction = {
          type: 'users_LOAD_SUCCESS',
          payload: {
            data: {
              users: {
                [user.id]: user,
              },
            },
          },
        };
        const initialState = {
          loading: false,
          entities: {},
          errors: {},
        };
        const newState = reducer(initialState, loadSuccessAction);

        expect(newState.data[user.id]).toEqual(user);
      });
    });

    describe('unsuccessful loadAll call', () => {
      const mockStore = reduxMockStore(store);
      const errors = [
        {
          name: 'base',
          reason: 'Unable to load users',
        },
      ];
      const errorResponse = {
        message: {
          message: 'Cannot get users',
          errors,
        },
      };
      const formattedErrors = formatErrorResponse(errorResponse);
      const loadAllFunc = createSpy().andCall(() => {
        return Promise.reject(errorResponse);
      });
      const config = reduxConfig({
        entityName: 'users',
        loadAllFunc,
        schema: schemas.USERS,
      });
      const { actions, reducer } = config;

      it('calls the loadAllFunc', () => {
        mockStore.dispatch(actions.loadAll());

        expect(loadAllFunc).toHaveBeenCalled();
      });

      it('dispatches the correct actions', () => {
        mockStore.dispatch(actions.loadAll());

        const dispatchedActions = mockStore.getActions();
        const dispatchedActionTypes = dispatchedActions.map((action) => { return action.type; });
        const loadAllFailureAction = find(dispatchedActions, { type: 'users_LOAD_FAILURE' });

        expect(dispatchedActionTypes).toInclude('users_LOAD_REQUEST');
        expect(dispatchedActionTypes).toNotInclude('users_LOAD_SUCCESS');

        expect(loadAllFailureAction.payload).toEqual({
          errors: formattedErrors,
        });
      });

      it('adds the returned errors to state', () => {
        const loadAllFailureAction = {
          type: 'users_LOAD_FAILURE',
          payload: {
            errors: formattedErrors,
          },
        };
        const initialState = {
          loading: false,
          entities: {},
          errors: {},
        };
        const newState = reducer(initialState, loadAllFailureAction);

        expect(newState.errors).toEqual(formattedErrors);
      });
    });
  });
});
