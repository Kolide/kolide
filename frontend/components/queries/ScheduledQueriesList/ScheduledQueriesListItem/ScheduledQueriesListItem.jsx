import React, { Component, PropTypes } from 'react';

import Checkbox from 'components/forms/fields/Checkbox';
import Icon from 'components/icons/Icon';
import PlatformIcon from 'components/icons/PlatformIcon';
import { isEmpty, isEqual } from 'lodash';
import scheduledQueryInterface from 'interfaces/scheduled_query';

class ScheduledQueriesListItem extends Component {
  static propTypes = {
    checked: PropTypes.bool,
    disabled: PropTypes.bool,
    onSelect: PropTypes.func.isRequired,
    scheduledQuery: scheduledQueryInterface.isRequired,
  };

  shouldComponentUpdate (nextProps) {
    if (isEqual(nextProps, this.props)) {
      return false;
    }

    return true;
  }

  onCheck = (value) => {
    const { onSelect, scheduledQuery } = this.props;

    return onSelect(value, scheduledQuery.id);
  }

  loggingTypeString = () => {
    const { scheduledQuery: { snapshot, removed } } = this.props;

    if (snapshot) {
      return 'camera';
    }

    if (removed) {
      return 'plus-minus';
    }

    return 'bold-plus';
  }

  renderPlatformIcon = () => {
    const { scheduledQuery: { platform } } = this.props;
    const platformArr = platform ? platform.split(',') : [];

    if (isEmpty(platformArr) || platformArr.includes('all')) {
      return <PlatformIcon name="" />;
    }

    return platformArr.map((pltf, idx) => <PlatformIcon name={pltf} key={`${idx}-${pltf}`} />);
  }

  render () {
    const { checked, disabled, scheduledQuery } = this.props;
    const { onCheck, renderPlatformIcon } = this;
    const { id, name, interval, shard, version } = scheduledQuery;
    const { loggingTypeString } = this;

    return (
      <tr>
        <td>
          <Checkbox
            disabled={disabled}
            name={`scheduled-query-checkbox-${id}`}
            onChange={onCheck}
            value={checked}
          />
        </td>
        <td className="scheduled-queries-list__query-name">{name}</td>
        <td>{interval}</td>
        <td>{renderPlatformIcon()}</td>
        <td>{version ? `${version}+` : 'Any'}</td>
        <td>{shard}</td>
        <td><Icon name={loggingTypeString()} /></td>
      </tr>
    );
  }
}

export default ScheduledQueriesListItem;

