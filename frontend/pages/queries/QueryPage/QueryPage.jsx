import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';
import { connect } from 'react-redux';
import FileSaver from 'file-saver';
import { clone, filter, includes, isArray, isEqual, merge } from 'lodash';
import moment from 'moment';
import { push } from 'react-router-redux';

import Kolide from 'kolide';
import campaignHelpers from 'redux/nodes/entities/campaigns/helpers';
import convertToCSV from 'utilities/convert_to_csv';
import debounce from 'utilities/debounce';
import deepDifference from 'utilities/deep_difference';
import entityGetter from 'redux/utilities/entityGetter';
import { formatSelectedTargetsForApi } from 'kolide/helpers';
import hostActions from 'redux/nodes/entities/hosts/actions';
import QueryForm from 'components/forms/queries/QueryForm';
import osqueryTableInterface from 'interfaces/osquery_table';
import queryActions from 'redux/nodes/entities/queries/actions';
import queryInterface from 'interfaces/query';
import QueryPageSelectTargets from 'components/queries/QueryPageSelectTargets';
import QueryResultsTable from 'components/queries/QueryResultsTable';
import QuerySidePanel from 'components/side_panels/QuerySidePanel';
import { renderFlash } from 'redux/nodes/notifications/actions';
import { selectOsqueryTable, setSelectedTargets, setSelectedTargetsQuery } from 'redux/nodes/components/QueryPages/actions';
import targetInterface from 'interfaces/target';
import validateQuery from 'components/forms/validators/validate_query';
import Spinner from 'components/loaders/Spinner';

const baseClass = 'query-page';
const DEFAULT_CAMPAIGN = {
  hosts_count: {
    total: 0,
  },
};

export class QueryPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    errors: PropTypes.shape({
      base: PropTypes.string,
    }),
    hostIDs: PropTypes.oneOfType([PropTypes.array, PropTypes.string]),
    loadingQueries: PropTypes.bool.isRequired,
    location: PropTypes.shape({
      pathname: PropTypes.string,
    }),
    query: queryInterface,
    selectedOsqueryTable: osqueryTableInterface,
    selectedTargets: PropTypes.arrayOf(targetInterface),
  };

  static defaultProps = {
    loadingQueries: false,
  };

  constructor (props) {
    super(props);

    this.state = {
      campaign: DEFAULT_CAMPAIGN,
      queryIsRunning: false,
      queryText: props.query.query,
      targetsCount: 0,
      targetsError: null,
      isQueryFullScreen: false,
      queryPosition: {},
    };

    this.csvQueryName = 'Query Results';
  }

  componentWillMount () {
    const { dispatch, hostIDs } = this.props;

    if (hostIDs) {
      dispatch(hostActions.loadAll());
    }

    return false;
  }

  componentWillReceiveProps (nextProps) {
    const nextPathname = nextProps.location.pathname;
    const { pathname } = this.props.location;

    if (nextPathname !== pathname) {
      this.resetCampaignAndTargets();
    }

    return false;
  }

  componentWillUnmount () {
    this.resetCampaignAndTargets();

    return false;
  }

  onChangeQueryFormField = (fieldName, value) => {
    if (fieldName === 'name') {
      this.csvQueryName = value;
    }

    if (fieldName === 'query') {
      this.setState({ queryText: value });
    }

    return false;
  }

  onExportQueryResults = (evt) => {
    evt.preventDefault();

    const { campaign } = this.state;
    const { query_results: queryResults } = campaign;

    if (queryResults) {
      const csv = convertToCSV(queryResults, (fields) => {
        const result = filter(fields, f => f !== 'host_hostname');

        result.unshift('host_hostname');

        return result;
      });
      const formattedTime = moment(new Date()).format('MM-DD-YY hh-mm-ss');
      const filename = `${this.csvQueryName} (${formattedTime}).csv`;
      const file = new global.window.File([csv], filename, { type: 'text/csv' });

      FileSaver.saveAs(file);
    }

    return false;
  }

  onFetchTargets = (query, targetResponse) => {
    const { dispatch } = this.props;
    const {
      targets_count: targetsCount,
    } = targetResponse;

    dispatch(setSelectedTargetsQuery(query));
    this.setState({ targetsCount });

    return false;
  }

  onOsqueryTableSelect = (tableName) => {
    const { dispatch } = this.props;

    dispatch(selectOsqueryTable(tableName));

    return false;
  }

  onRunQuery = debounce((queryText) => {
    const { dispatch, selectedTargets } = this.props;
    const { error } = validateQuery(queryText);

    if (!selectedTargets.length) {
      this.setState({ targetsError: 'You must select at least one target to run a query' });

      return false;
    }

    if (error) {
      dispatch(renderFlash('error', error));

      return false;
    }

    const { destroyCampaign, removeSocket } = this;
    const selected = formatSelectedTargetsForApi(selectedTargets);

    removeSocket();
    destroyCampaign();

    Kolide.queries.run({ query: queryText, selected })
      .then((campaignResponse) => {
        return Kolide.websockets.queries.run(campaignResponse.id)
          .then((socket) => {
            this.setState({ campaign: campaignResponse });
            this.socket = socket;
            this.setState({ queryIsRunning: true });

            this.socket.onmessage = ({ data }) => {
              const socketData = JSON.parse(data);
              const { previousSocketData } = this;

              if (previousSocketData && isEqual(socketData, previousSocketData)) {
                this.previousSocketData = socketData;

                return false;
              }

              return campaignHelpers.update(this.state.campaign, socketData)
                .then((updatedCampaign) => {
                  const { status } = updatedCampaign;

                  if (status === 'finished') {
                    this.setState({ queryIsRunning: false });
                    removeSocket();

                    return false;
                  }

                  this.previousSocketData = socketData;
                  this.setState({ campaign: updatedCampaign });

                  return false;
                });
            };
          });
      })
      .catch((campaignError) => {
        if (campaignError === 'resource already created') {
          dispatch(renderFlash('error', 'A campaign with the provided query text has already been created'));

          return false;
        }

        dispatch(renderFlash('error', campaignError));

        return false;
      });

    return false;
  })

  onSaveQueryFormSubmit = debounce((formData) => {
    const { dispatch } = this.props;
    const { error } = validateQuery(formData.query);

    if (error) {
      dispatch(renderFlash('error', error));

      return false;
    }

    return dispatch(queryActions.create(formData))
      .then((query) => {
        dispatch(push(`/queries/${query.id}`));
        dispatch(renderFlash('success', 'Query created'));
      })
      .catch(() => false);
  })

  onStopQuery = (evt) => {
    evt.preventDefault();

    const { removeSocket } = this;

    this.setState({ queryIsRunning: false });

    return removeSocket();
  }

  onTargetSelect = (selectedTargets) => {
    const { dispatch } = this.props;

    this.setState({ targetsError: null });

    dispatch(setSelectedTargets(selectedTargets));

    return false;
  }

  onUpdateQuery = (formData) => {
    const { dispatch, query } = this.props;
    const updatedQuery = deepDifference(formData, query);

    dispatch(queryActions.update(query, updatedQuery))
      .then(() => {
        dispatch(renderFlash('success', 'Query updated!'));
      });

    return false;
  };

  onToggleQueryFullScreen = (evt) => {
    const { document: { body } } = global;
    const { isQueryFullScreen, queryPosition } = this.state;
    const rect = evt.target.parentNode.getBoundingClientRect();
    const parent = evt.target.parentNode;
    const grandParent = evt.target.parentNode.parentNode;
    const defaultPosition = {
      top: `${rect.top + body.scrollTop}px`,
      left: `${rect.left + body.scrollLeft}px`,
      right: `${rect.right - rect.left}px`,
      bottom: `${rect.bottom - parent.offsetHeight - rect.top}px`,
      maxWidth: `${parent.offsetWidth}px`,
      minWidth: `${parent.offsetWidth}px`,
      maxHeight: `${parent.offsetHeight}px`,
      position: 'fixed',
    };
    let newPosition = clone(defaultPosition);

    if(!isQueryFullScreen) {
      this.setState({ queryPosition: defaultPosition });
    } else {
      newPosition = queryPosition;
      window.setTimeout(function(){
        parent.style.position = 'static';
      }, 2500);
    }

    merge(parent.style, newPosition);
    // grandParent.style.minHeight = `${parent.offsetHeight}px`;

    this.setState({ isQueryFullScreen: !isQueryFullScreen });
    return false;
  }

  destroyCampaign = () => {
    const { campaign } = this.state;

    if (this.campaign || campaign) {
      this.campaign = null;
      this.setState({ campaign: DEFAULT_CAMPAIGN });
    }

    return false;
  }

  removeSocket = () => {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
      this.previousSocketData = null;
    }

    return false;
  }

  resetCampaignAndTargets = () => {
    const { destroyCampaign, removeSocket } = this;
    const { dispatch } = this.props;

    destroyCampaign();
    dispatch(setSelectedTargets([]));
    removeSocket();

    return false;
  }

  renderResultsTable = () => {
    const { campaign, queryIsRunning, isQueryFullScreen } = this.state;
    const { onExportQueryResults, onToggleQueryFullScreen } = this;
    const loading = queryIsRunning && !campaign.hosts_count.total;
    const resultsClasses = classnames(`${baseClass}__results`, 'body-wrap', {
      [`${baseClass}__results--loading`]: loading,
      [`${baseClass}__results--full-screen`]: isQueryFullScreen,
    });
    let resultBody = '';

    if (!loading && isEqual(campaign, DEFAULT_CAMPAIGN)) {
      return false;
    }

    if (loading) {
      resultBody = <Spinner />;
    } else {
      resultBody = <QueryResultsTable campaign={campaign} onExportQueryResults={onExportQueryResults} isQueryFullScreen={isQueryFullScreen} onToggleQueryFullScreen={onToggleQueryFullScreen} />;
    }

    return (
      <div className={resultsClasses}>
        {resultBody}
      </div>
    );
  }

  renderTargetsInput = () => {
    const { onFetchTargets, onRunQuery, onStopQuery, onTargetSelect } = this;
    const { campaign, queryIsRunning, queryText, targetsCount, targetsError } = this.state;
    const { selectedTargets } = this.props;

    return (
      <QueryPageSelectTargets
        campaign={campaign}
        error={targetsError}
        onFetchTargets={onFetchTargets}
        onRunQuery={onRunQuery}
        onStopQuery={onStopQuery}
        onTargetSelect={onTargetSelect}
        query={queryText}
        queryIsRunning={queryIsRunning}
        selectedTargets={selectedTargets}
        targetsCount={targetsCount}
      />
    );
  }

  render () {
    const {
      onChangeQueryFormField,
      onOsqueryTableSelect,
      onRunQuery,
      onSaveQueryFormSubmit,
      onStopQuery,
      onTextEditorInputChange,
      onUpdateQuery,
      renderResultsTable,
      renderTargetsInput,
    } = this;
    const { queryIsRunning } = this.state;
    const {
      errors,
      loadingQueries,
      query,
      selectedOsqueryTable,
    } = this.props;

    if (loadingQueries) {
      return false;
    }

    return (
      <div className={`${baseClass} has-sidebar`}>
        <div className={`${baseClass}__content`}>
          <div className={`${baseClass}__form body-wrap`}>
            <QueryForm
              formData={query}
              handleSubmit={onSaveQueryFormSubmit}
              onChangeFunc={onChangeQueryFormField}
              onOsqueryTableSelect={onOsqueryTableSelect}
              onRunQuery={onRunQuery}
              onStopQuery={onStopQuery}
              onUpdate={onUpdateQuery}
              queryIsRunning={queryIsRunning}
              serverErrors={errors}
              selectedOsqueryTable={selectedOsqueryTable}
            />
          </div>
          {renderTargetsInput()}
          {renderResultsTable()}
        </div>
        <QuerySidePanel
          onOsqueryTableSelect={onOsqueryTableSelect}
          onTextEditorInputChange={onTextEditorInputChange}
          selectedOsqueryTable={selectedOsqueryTable}
        />
      </div>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  const stateEntities = entityGetter(state);
  const { id: queryID } = ownProps.params;
  const reduxQuery = entityGetter(state).get('queries').findBy({ id: queryID });
  const { queryText, selectedOsqueryTable } = state.components.QueryPages;
  const { errors, loading: loadingQueries } = state.entities.queries;
  const queryStub = { description: '', name: '', query: queryText };
  const query = reduxQuery || queryStub;
  let { selectedTargets } = state.components.QueryPages;
  const { host_ids: hostIDs } = ownProps.location.query;

  // hostIDs are URL params so they are strings
  if (hostIDs && !queryID) {
    const { entities: hosts } = stateEntities.get('hosts');
    let hostFilter;

    if (isArray(hostIDs)) {
      hostFilter = h => includes(hostIDs, String(h.id));
    } else {
      hostFilter = { id: Number(hostIDs) };
    }

    selectedTargets = filter(hosts, hostFilter);
  }

  return {
    errors,
    hostIDs,
    loadingQueries,
    query,
    selectedOsqueryTable,
    selectedTargets,
  };
};

export default connect(mapStateToProps)(QueryPage);
