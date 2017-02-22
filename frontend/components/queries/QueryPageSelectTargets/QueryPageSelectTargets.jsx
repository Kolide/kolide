import React, { Component, PropTypes } from 'react';

import campaignInterface from 'interfaces/campaign';
import QueryProgressDetails from 'components/queries/QueryProgressDetails';
import SelectTargetsDropdown from 'components/forms/fields/SelectTargetsDropdown';
import targetInterface from 'interfaces/target';

const baseClass = 'query-page-select-targets';

class QueryPageSelectTargets extends Component {
  static propTypes = {
    campaign: campaignInterface,
    error: PropTypes.string,
    onFetchTargets: PropTypes.func.isRequired,
    onRunQuery: PropTypes.func.isRequired,
    onStopQuery: PropTypes.func.isRequired,
    onTargetSelect: PropTypes.func.isRequired,
    query: PropTypes.string,
    queryIsRunning: PropTypes.bool,
    selectedTargets: PropTypes.arrayOf(targetInterface),
    targetsCount: PropTypes.number,
  };

  render () {
    const {
      error,
      onFetchTargets,
      onTargetSelect,
      selectedTargets,
      targetsCount,
      campaign,
      onRunQuery,
      onStopQuery,
      query,
      queryIsRunning,
    } = this.props;

    return (
      <div className={`${baseClass}__wrapper body-wrap`}>
        <QueryProgressDetails
          campaign={campaign}
          onRunQuery={onRunQuery}
          onStopQuery={onStopQuery}
          query={query}
          queryIsRunning={queryIsRunning}
        />
        <SelectTargetsDropdown
          error={error}
          onFetchTargets={onFetchTargets}
          onSelect={onTargetSelect}
          selectedTargets={selectedTargets}
          targetsCount={targetsCount}
          label="Select Targets"
        />
      </div>
    );
  }
}

export default QueryPageSelectTargets;
