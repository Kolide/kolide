import React, { Component, PropTypes } from 'react';
import moment from 'moment';

import Checkbox from 'components/forms/fields/Checkbox';
import ClickableTableRow from 'components/ClickableTableRow';
import { isEqual } from 'lodash';
import queryInterface from 'interfaces/query';

class QueriesListRow extends Component {
  static propTypes = {
    checked: PropTypes.bool,
    onSelect: PropTypes.func.isRequired,
    query: queryInterface.isRequired,
  };

  shouldComponentUpdate (nextProps) {
    if (isEqual(nextProps, this.props)) {
      return false;
    }

    return true;
  }

  onCheck = (value) => {
    const { onCheck: handleCheck, query } = this.props;

    return handleCheck(value, query.id);
  }

  onSelect = () => {
    const { onSelect: handleSelect, query } = this.props;

    return handleSelect(query);
  }

  render () {
    const { checked, query } = this.props;
    const { onCheck, onSelect } = this;
    const { author_name: authorName, id, name, updated_at: updatedAt } = query;
    const lastModifiedDate = moment(updatedAt).format('MM/DD/YY');

    return (
      <ClickableTableRow onClick={onSelect}>
        <td>
          <Checkbox
            name={`query-checkbox-${id}`}
            onChange={onCheck}
            value={checked}
          />
        </td>
        <td>{name}</td>
        <td>{authorName}</td>
        <td>{lastModifiedDate}</td>
      </ClickableTableRow>
    );
  }
}

export default QueriesListRow;


