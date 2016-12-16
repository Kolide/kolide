import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';
import { noop } from 'lodash';

import { scheduledQueryStub, queryStub } from 'test/stubs';
import QueriesList from './index';

const scheduledQueries = [
  { ...scheduledQueryStub, id: 1 },
  { ...scheduledQueryStub, id: 2 },
];

describe('QueriesList - component', () => {
  it('renders a QueriesListItem for each scheduled query', () => {
    const component = mount(
      <QueriesList
        allQueries={[]}
        onHidePackForm={noop}
        onSelectQuery={noop}
        scheduledQueries={scheduledQueries}
        selectedScheduledQueryIDs={[]}
      />
    );

    expect(component.find('QueriesListItem').length).toEqual(2);
  });

  it('renders "No queries found" help text when scheduled queries are available but the scheduled queries are filtered out', () => {
    const component = mount(
      <QueriesList
        allQueries={[]}
        isScheduledQueriesAvailable
        onHidePackForm={noop}
        onSelectQuery={noop}
        scheduledQueries={[]}
        selectedScheduledQueryIDs={[]}
      />
    );

    expect(component.text()).toInclude('No queries matched your search criteria');
  });

  it('renders initial help text when no queries have been scheduled yet', () => {
    const component = mount(
      <QueriesList
        allQueries={[]}
        onHidePackForm={noop}
        onSelectQuery={noop}
        scheduledQueries={[]}
        selectedScheduledQueryIDs={[]}
      />
    );

    expect(component.text()).toInclude("First let's add a query");
  });

  it('renders the pack config form when the prop is present', () => {
    const component = mount(
      <QueriesList
        allQueries={[]}
        onHidePackForm={noop}
        onSelectQuery={noop}
        scheduledQueries={[]}
        shouldShowPackForm
        selectedScheduledQueryIDs={[]}
      />
    );

    expect(component.find('PackQueryConfigForm').length).toEqual(1);
  });

  it('sets the query dropdown options in state', () => {
    const componentWithoutOptions = mount(
      <QueriesList
        allQueries={[]}
        onHidePackForm={noop}
        onSelectQuery={noop}
        scheduledQueries={[]}
        shouldShowPackForm
        selectedScheduledQueryIDs={[]}
      />
    );
    const componentWithOptions = mount(
      <QueriesList
        allQueries={[queryStub]}
        onHidePackForm={noop}
        onSelectQuery={noop}
        scheduledQueries={[]}
        shouldShowPackForm
        selectedScheduledQueryIDs={[]}
      />
    );

    expect(componentWithoutOptions.state().queryDropdownOptions).toEqual([]);
    expect(componentWithOptions.state().queryDropdownOptions).toEqual([{ label: queryStub.name, value: String(queryStub.id) }]);
  });
});
