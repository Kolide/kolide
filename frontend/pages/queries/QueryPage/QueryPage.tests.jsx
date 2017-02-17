import React from 'react';
import expect, { spyOn, restoreSpies } from 'expect';
import FileSave from 'file-saver';
import { mount } from 'enzyme';
import nock from 'nock';
import { noop } from 'lodash';

import convertToCSV from 'utilities/convert_to_csv';
import * as queryPageActions from 'redux/nodes/components/QueryPages/actions';
import helpers from 'test/helpers';
import hostActions from 'redux/nodes/entities/hosts/actions';
import queryActions from 'redux/nodes/entities/queries/actions';
import ConnectedQueryPage, { QueryPage } from 'pages/queries/QueryPage/QueryPage';
import { hostStub, queryStub } from 'test/stubs';

const { connectedComponent, createAceSpy, fillInFormInput, reduxMockStore } = helpers;
const { defaultSelectedOsqueryTable } = queryPageActions;
const locationProp = { params: {}, location: { query: {} } };

describe('QueryPage - component', () => {
  beforeEach(() => {
    createAceSpy();

    spyOn(hostActions, 'loadAll')
      .andReturn(() => Promise.resolve([]));

    nock('http://localhost:8080')
      .post('/api/v1/kolide/targets', JSON.stringify({
        selected: {
          hosts: [1, 99],
          labels: [],
        },
      }))
      .reply(200, {
        targets_count: 1234,
        targets: [
          {
            id: 3,
            label: 'OS X El Capitan 10.11',
            name: 'osx-10.11',
            platform: 'darwin',
            target_type: 'hosts',
          },
        ],
      });
  });

  afterEach(() => {
    nock.cleanAll();
    restoreSpies();
  });

  const store = {
    components: {
      QueryPages: {
        queryText: 'SELECT * FROM users',
        selectedOsqueryTable: defaultSelectedOsqueryTable,
        selectedTargets: [],
      },
    },
    entities: {
      hosts: {
        data: {
          [hostStub.id]: hostStub,
          99: { ...hostStub, id: 99 },
        },
      },
      queries: { loading: false, data: {} },
      targets: {},
    },
  };
  const mockStore = reduxMockStore(store);

  describe('rendering', () => {
    it('does not render when queries are loading', () => {
      const loadingQueriesStore = {
        ...store,
        entities: {
          ...store.entities,
          queries: { loading: true, data: {} },
        },
      };
      const page = mount(connectedComponent(ConnectedQueryPage, {
        mockStore: reduxMockStore(loadingQueriesStore),
        props: locationProp,
      }));

      expect(page.html()).toNotExist();
    });

    it('renders the QueryForm component', () => {
      const page = mount(connectedComponent(ConnectedQueryPage, { mockStore, props: locationProp }));

      expect(page.find('QueryForm').length).toEqual(1);
    });

    it('renders the QuerySidePanel component', () => {
      const page = mount(connectedComponent(ConnectedQueryPage, { mockStore, props: locationProp }));

      expect(page.find('QuerySidePanel').length).toEqual(1);
    });

    describe('results table', () => {
      const props = { dispatch: noop, query: queryStub, selectedOsqueryTable: defaultSelectedOsqueryTable };
      const queryResult = { org_name: 'Kolide', org_url: 'https://kolide.co' };
      const campaign = {
        id: 1,
        hosts_count: { failed: 0, successful: 1, total: 1 },
        query_results: [queryResult],
        totals: { count: 1 },
      };

      it('renders the results table when there are result totals', () => {
        const Page = mount(<QueryPage {...props} />);

        Page.setState({ campaign });

        expect(Page.find('.query-page__results').length)
          .toEqual(1, 'Expected the results section to render');
      });

      it('does not render the results table when there are no result totals', () => {
        const Page = mount(<QueryPage {...props} />);
        const campaignWithoutResults = { id: 1, hosts_count: { failed: 0, successful: 0, total: 0 } };

        Page.setState({ campaign: campaignWithoutResults });

        expect(Page.find('.query-page__results').length)
          .toEqual(0, 'Expected the results section not to render');
      });
    });
  });

  it('sets selectedTargets based on host_ids', () => {
    const singleHostProps = { params: {}, location: { query: { host_ids: String(hostStub.id) } } };
    const multipleHostsProps = { params: {}, location: { query: { host_ids: [String(hostStub.id), '99'] } } };
    const singleHostPage = mount(connectedComponent(ConnectedQueryPage, { mockStore, props: singleHostProps }));
    const multipleHostsPage = mount(connectedComponent(ConnectedQueryPage, { mockStore, props: multipleHostsProps }));

    expect(singleHostPage.find('QueryPage').prop('selectedTargets')).toEqual([hostStub]);
    expect(multipleHostsPage.find('QueryPage').prop('selectedTargets')).toEqual([hostStub, { ...hostStub, id: 99 }]);
  });

  it('sets targetError in state when the query is run and there are no selected targets', () => {
    const page = mount(connectedComponent(ConnectedQueryPage, { mockStore, props: locationProp }));
    const QueryPageSelectTargets = page.find('QueryPageSelectTargets');
    const runQueryBtn = page.find('.query-page-select-targets__run-query-btn');

    expect(QueryPageSelectTargets.prop('error')).toNotExist();

    runQueryBtn.simulate('click');

    expect(QueryPageSelectTargets.prop('error')).toEqual('You must select at least one target to run a query');
  });

  it('calls the onUpdateQuery prop when the query is updated', () => {
    const query = { id: 1, name: 'My query', description: 'My query description', query: 'select * from users' };
    const locationWithQueryProp = { params: { id: 1 }, location: { query: {} } };
    const mockStoreWithQuery = reduxMockStore({
      components: {
        QueryPages: {
          queryText: 'SELECT * FROM users',
          selectedOsqueryTable: defaultSelectedOsqueryTable,
          selectedTargets: [],
        },
      },
      entities: {
        queries: {
          data: {
            1: query,
          },
        },
      },
    });
    const page = mount(connectedComponent(ConnectedQueryPage, {
      mockStore: mockStoreWithQuery,
      props: locationWithQueryProp,
    }));
    const form = page.find('QueryForm');
    const nameInput = form.find({ name: 'name' }).find('input');
    const saveChangesBtn = form.find('li.dropdown-button__option').first().find('Button');
    fillInFormInput(nameInput, 'new name');
    spyOn(queryActions, 'update').andReturn(() => Promise.resolve({
      description: query.description,
      name: 'new name',
      queryText: 'SELECT * FROM users',
    }));

    form.simulate('submit');
    saveChangesBtn.simulate('click');

    expect(queryActions.update).toHaveBeenCalledWith(query, { name: 'new name' });
  });

  describe('#componentWillReceiveProps', () => {
    it('resets selected targets and removed the campaign when the hostname changes', () => {
      const queryResult = { org_name: 'Kolide', org_url: 'https://kolide.co' };
      const campaign = { id: 1, query_results: [queryResult], hosts_count: { total: 1 } };
      const props = {
        dispatch: noop,
        loadingQueries: false,
        location: { pathname: '/queries/11' },
        query: { query: 'select * from users' },
        selectedOsqueryTable: defaultSelectedOsqueryTable,
        selectedTargets: [hostStub],
      };
      const Page = mount(<QueryPage {...props} />);
      const PageNode = Page.node;

      spyOn(PageNode, 'destroyCampaign');
      spyOn(PageNode, 'removeSocket');
      spyOn(queryPageActions, 'setSelectedTargets');

      Page.setState({ campaign });
      Page.setProps({ location: { pathname: '/queries/new' } });

      expect(queryPageActions.setSelectedTargets).toHaveBeenCalledWith([]);
      expect(PageNode.destroyCampaign).toHaveBeenCalled();
      expect(PageNode.removeSocket).toHaveBeenCalled();
    });
  });

  describe('export as csv', () => {
    it('exports the campaign query results in csv format', () => {
      const queryResult = { org_name: 'Kolide', org_url: 'https://kolide.co' };
      const campaign = {
        id: 1,
        hosts_count: {
          failed: 0,
          successful: 1,
          total: 1,
        },
        query_results: [queryResult],
        totals: {
          count: 1,
        },
      };
      const queryResultsCSV = convertToCSV([queryResult]);
      const fileSaveSpy = spyOn(FileSave, 'saveAs');
      const Page = mount(<QueryPage dispatch={noop} query={queryStub} selectedOsqueryTable={defaultSelectedOsqueryTable} />);
      const filename = 'query_results.csv';
      const fileStub = new global.window.File([queryResultsCSV], filename, { type: 'text/csv' });

      Page.setState({ campaign });
      Page.node.socket = {};

      const QueryResultsTable = Page.find('QueryResultsTable');

      QueryResultsTable.find('Button').simulate('click');

      expect(fileSaveSpy).toHaveBeenCalledWith(fileStub);
    });
  });
});
