import React from 'react';
import expect from 'expect';
import { find } from 'lodash';
import { mount } from 'enzyme';

import ConnectedManageQueriesPage, { ManageQueriesPage } from 'pages/queries/ManageQueriesPage/ManageQueriesPage';
import { connectedComponent, fillInFormInput, reduxMockStore } from 'test/helpers';
import { queryStub } from 'test/stubs';

const store = {
  entities: {
    queries: {
      data: {
        [queryStub.id]: queryStub,
        101: {
          ...queryStub,
          id: 101,
          name: 'My unique query name',
        },
      },
    },
  },
};

describe('ManageQueriesPage - component', () => {
  it('filters the queries list', () => {
    const Component = connectedComponent(ConnectedManageQueriesPage, {
      mockStore: reduxMockStore(store),
    });
    const page = mount(Component).find('ManageQueriesPage');
    const queryFilterInput = page.find({ name: 'query-filter' }).find('input');

    expect(page.node.getQueries().length).toEqual(2);

    fillInFormInput(queryFilterInput, 'My unique query name');

    expect(page.node.getQueries().length).toEqual(1);
  });

  it('renders a QueriesList component', () => {
    const page = mount(connectedComponent(ConnectedManageQueriesPage));

    expect(page.find('QueriesList').length).toEqual(1);
  });

  it('renders the QueryDetailsSidePanel when a query is selected', () => {
    const mockStore = reduxMockStore(store);
    const props = { location: { query: { selectedQuery: queryStub.id } } };
    const Component = connectedComponent(ConnectedManageQueriesPage, { mockStore, props });
    const page = mount(Component).find('ManageQueriesPage');

    expect(page.find('QueryDetailsSidePanel').length).toEqual(1);
  });

  it('updates checkedQueryIDs in state when the check all queries Checkbox is toggled', () => {
    const page = mount(<ManageQueriesPage queries={[queryStub]} />);
    const selectAllQueries = page.find({ name: 'check-all-queries' });

    expect(page.state('checkedQueryIDs')).toEqual([]);

    selectAllQueries.simulate('change');

    expect(page.state('checkedQueryIDs')).toEqual([queryStub.id]);

    selectAllQueries.simulate('change');

    expect(page.state('checkedQueryIDs')).toEqual([]);
  });

  it('updates checkedQueryIDs in state when a query row Checkbox is toggled', () => {
    const page = mount(<ManageQueriesPage queries={[queryStub]} />);
    const queryCheckbox = page.find({ name: `query-checkbox-${queryStub.id}` });

    expect(page.state('checkedQueryIDs')).toEqual([]);

    queryCheckbox.simulate('change');

    expect(page.state('checkedQueryIDs')).toEqual([queryStub.id]);

    queryCheckbox.simulate('change');

    expect(page.state('checkedQueryIDs')).toEqual([]);
  });

  describe('bulk delete action', () => {
    const queries = [queryStub, { ...queryStub, id: 101, name: 'My unique query name' }];

    it('displays the delete action button when a query is checked', () => {
      const page = mount(<ManageQueriesPage queries={queries} />);
      const checkAllQueries = page.find({ name: 'check-all-queries' });

      checkAllQueries.simulate('change');

      expect(page.state('checkedQueryIDs')).toEqual([queryStub.id, 101]);
      expect(page.find('.manage-queries-page__delete-queries-btn').length).toEqual(1);
    });

    it('dispatches the query destroy function when the delete button is clicked', () => {
      const mockStore = reduxMockStore(store);
      const Component = connectedComponent(ConnectedManageQueriesPage, { mockStore });
      const page = mount(Component).find('ManageQueriesPage');
      const checkAllQueries = page.find({ name: 'check-all-queries' });

      checkAllQueries.simulate('change');

      const deleteBtn = page.find('.manage-queries-page__delete-queries-btn');

      deleteBtn.simulate('click');

      const dispatchedActions = mockStore.getActions();

      expect(dispatchedActions).toInclude({ type: 'queries_DESTROY_REQUEST' });
    });
  });

  describe('selecting a query', () => {
    it('updates the URL when a query is selected', () => {
      const mockStore = reduxMockStore(store);
      const Component = connectedComponent(ConnectedManageQueriesPage, { mockStore });
      const page = mount(Component).find('ManageQueriesPage');
      const firstRow = page.find('QueriesListRow').first();

      expect(page.prop('selectedQuery')).toNotExist();

      firstRow.find('ClickableTableRow').first().simulate('click');

      const dispatchedActions = mockStore.getActions();
      const locationChangeAction = find(dispatchedActions, { type: '@@router/CALL_HISTORY_METHOD' });

      expect(locationChangeAction.payload.args).toEqual([{
        pathname: '/queries/manage',
        query: { selectedQuery: queryStub.id },
      }]);
    });

    it('sets the selectedQuery prop', () => {
      const mockStore = reduxMockStore(store);
      const props = { location: { query: { selectedQuery: queryStub.id } } };
      const Component = connectedComponent(ConnectedManageQueriesPage, { mockStore, props });
      const page = mount(Component).find('ManageQueriesPage');

      expect(page.prop('selectedQuery')).toEqual(queryStub);
    });
  });
});
