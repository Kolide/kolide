import React from 'react';
import expect, { createSpy, restoreSpies } from 'expect';
import { mount } from 'enzyme';

import { queryStub, scheduledQueryStub } from 'test/stubs';
import { fillInFormInput } from 'test/helpers';
import QueriesListWrapper from './index';

const allQueries = [queryStub];
const scheduledQueries = [
  scheduledQueryStub,
  { ...scheduledQueryStub, id: 100, name: 'mac hosts' },
];

describe('QueriesListWrapper - component', () => {
  afterEach(restoreSpies);

  it('renders the "Remove Query" button when queries have been selected', () => {
    const component = mount(
      <QueriesListWrapper
        allQueries={allQueries}
        scheduledQueries={scheduledQueries}
      />
    );

    component.find('Checkbox').last().find('input').simulate('change');

    const addQueryBtn = component.find('Button').findWhere(b => b.prop('text') === 'Add New Query');
    const removeQueryBtn = component.find('Button').findWhere(b => b.prop('text') === 'Remove Query');

    expect(addQueryBtn.length).toEqual(0);
    expect(removeQueryBtn.length).toEqual(1);
  });

  it('calls the onRemoveScheduledQueries prop', () => {
    const spy = createSpy();
    const component = mount(
      <QueriesListWrapper
        allQueries={allQueries}
        onRemoveScheduledQueries={spy}
        scheduledQueries={[scheduledQueryStub]}
      />
    );

    component.find('Checkbox').last().find('input').simulate('change');

    const removeQueryBtn = component.find('Button').findWhere(b => b.prop('text') === 'Remove Query');

    removeQueryBtn.simulate('click');

    expect(spy).toHaveBeenCalledWith([scheduledQueryStub.id]);
  });

  it('filters queries', () => {
    const component = mount(
      <QueriesListWrapper
        allQueries={allQueries}
        scheduledQueries={scheduledQueries}
      />
    );

    const searchQueriesInput = component.find({ name: 'search-queries' });
    const QueriesList = component.find('QueriesList');

    expect(QueriesList.prop('scheduledQueries')).toEqual(scheduledQueries);

    fillInFormInput(searchQueriesInput, 'something that does not match');

    expect(QueriesList.prop('scheduledQueries')).toEqual([]);
  });
});
