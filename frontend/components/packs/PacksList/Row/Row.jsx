import React, { Component, PropTypes } from 'react';
import classNames from 'classnames';
import { isEqual } from 'lodash';
import moment from 'moment';

import Checkbox from 'components/forms/fields/Checkbox';
import ClickableTableRow from 'components/ClickableTableRow';
import Icon from 'components/icons/Icon';
import packInterface from 'interfaces/pack';

const baseClass = 'packs-list-row';

class Row extends Component {
  static propTypes = {
    checked: PropTypes.bool,
    onCheck: PropTypes.func,
    onSelect: PropTypes.func,
    pack: packInterface.isRequired,
  };

  shouldComponentUpdate (nextProps) {
    return !isEqual(this.props, nextProps);
  }

  handleChange = (shouldCheck) => {
    const { onCheck, pack } = this.props;

    return onCheck(shouldCheck, pack.id);
  }

  handleSelect = () => {
    const { onSelect, pack } = this.props;

    return onSelect(pack);
  }

  renderStatusData = () => {
    const { disabled } = this.props.pack;
    const iconClassName = classNames(`${baseClass}__status-icon`, {
      [`${baseClass}__status-icon--enabled`]: !disabled,
      [`${baseClass}__status-icon--disabled`]: disabled,
    });

    if (disabled) {
      return (
        <td className={`${baseClass}__td`}>
          <Icon className={iconClassName} name="offline" />
          <span className={`${baseClass}__status-text`}>Disabled</span>
        </td>
      );
    }

    return (
      <td className={`${baseClass}__td`}>
        <Icon className={iconClassName} name="success-check" />
        <span className={`${baseClass}__status-text`}>Enabled</span>
      </td>
    );
  }

  render () {
    const { checked, pack } = this.props;
    const { handleChange, handleSelect, renderStatusData } = this;
    const updatedTime = moment(pack.updated_at).format('MM/DD/YY');

    return (
      <ClickableTableRow className={baseClass} onClick={handleSelect}>
        <td className={`${baseClass}__td`}>
          <Checkbox
            name={`select-pack-${pack.id}`}
            onChange={handleChange}
            value={checked}
            wrapperClassName={`${baseClass}__checkbox`}
          />
        </td>
        <td className={`${baseClass}__td ${baseClass}__td-pack-name`}>{pack.name}</td>
        <td className={`${baseClass}__td ${baseClass}__td-query-count`}>{pack.query_count}</td>
        {renderStatusData()}
        <td />
        <td className={`${baseClass}__td`}>{updatedTime}</td>
      </ClickableTableRow>
    );
  }
}

export default Row;

