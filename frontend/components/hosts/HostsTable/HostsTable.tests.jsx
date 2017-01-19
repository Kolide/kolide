import React from 'react';
import expect, { createSpy, restoreSpies } from 'expect';
import { mount } from 'enzyme';

import { hostStub } from 'test/stubs';
import HostsTable from 'components/hosts/HostsTable';

describe('HostsTable - component', () => {
  afterEach(restoreSpies);

  it('calls the onDestroyHost prop when the trash icon button is clicked', () => {
    const spy = createSpy();
    const component = mount(<HostsTable hosts={[hostStub]} onDestroyHost={spy} />);
    const btn = component.find('Button');

    btn.simulate('click');

    expect(spy).toHaveBeenCalled();
  });
});
