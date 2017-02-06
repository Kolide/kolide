import React, { PropTypes } from 'react';

import ConfigurePackQueryForm from 'components/forms/ConfigurePackQueryForm';
import queryInterface from 'interfaces/query';
import scheduledQueryInterface from 'interfaces/scheduled_query';
import SearchPackQuery from './SearchPackQuery';
import SecondarySidePanelContainer from '../SecondarySidePanelContainer';

const baseClass = 'schedule-query-side-panel';

const ScheduleQuerySidePanel = ({
  allQueries,
  onConfigurePackQuerySubmit,
  onUpdateScheduledQuery,
  onSelectQuery,
  selectedQuery,
  selectedScheduledQuery,
}) => {
  const renderForm = () => {
    if (!selectedQuery) {
      return false;
    }

    const formData = selectedScheduledQuery || {};

    formData.query_id = selectedQuery.id;

    const handleSubmit = selectedScheduledQuery ? onUpdateScheduledQuery : onConfigurePackQuerySubmit;

    return (
      <ConfigurePackQueryForm
        formData={formData}
        handleSubmit={handleSubmit}
      />
    );
  };

  return (
    <SecondarySidePanelContainer className={baseClass}>
      <SearchPackQuery
        allQueries={allQueries}
        onSelectQuery={onSelectQuery}
        selectedQuery={selectedQuery}
      />
      {renderForm()}
    </SecondarySidePanelContainer>
  );
};

ScheduleQuerySidePanel.propTypes = {
  allQueries: PropTypes.arrayOf(queryInterface),
  onConfigurePackQuerySubmit: PropTypes.func,
  onSelectQuery: PropTypes.func,
  onUpdateScheduledQuery: PropTypes.func,
  selectedQuery: queryInterface,
  selectedScheduledQuery: scheduledQueryInterface,
};

export default ScheduleQuerySidePanel;
