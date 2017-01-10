import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { filter, includes, noop, size, find } from 'lodash';
import { push } from 'react-router-redux';

import EditPackFormWrapper from 'components/packs/EditPackFormWrapper';
import hostActions from 'redux/nodes/entities/hosts/actions';
import hostInterface from 'interfaces/host';
import labelActions from 'redux/nodes/entities/labels/actions';
import labelInterface from 'interfaces/label';
import packActions from 'redux/nodes/entities/packs/actions';
import ScheduleQuerySidePanel from 'components/side_panels/ScheduleQuerySidePanel';
import packInterface from 'interfaces/pack';
import queryActions from 'redux/nodes/entities/queries/actions';
import queryInterface from 'interfaces/query';
import ScheduledQueriesListWrapper from 'components/queries/ScheduledQueriesListWrapper';
import { renderFlash } from 'redux/nodes/notifications/actions';
import scheduledQueryActions from 'redux/nodes/entities/scheduled_queries/actions';
import stateEntityGetter from 'redux/utilities/entityGetter';

const baseClass = 'edit-pack-page';

export class EditPackPage extends Component {
  static propTypes = {
    allQueries: PropTypes.arrayOf(queryInterface),
    dispatch: PropTypes.func,
    isEdit: PropTypes.bool,
    isLoadingPack: PropTypes.bool,
    isLoadingScheduledQueries: PropTypes.bool,
    pack: packInterface,
    packHosts: PropTypes.arrayOf(hostInterface),
    packID: PropTypes.string,
    packLabels: PropTypes.arrayOf(labelInterface),
    scheduledQueries: PropTypes.arrayOf(queryInterface),
  };

  static defaultProps = {
    dispatch: noop,
  };

  constructor (props) {
    super(props);

    this.state = {
      targetsCount: 0,
    };
  }

  componentDidMount () {
    const {
      allQueries,
      dispatch,
      isLoadingPack,
      pack,
      packHosts,
      packID,
      packLabels,
      scheduledQueries,
    } = this.props;
    const { load } = packActions;
    const { loadAll } = queryActions;

    if (!pack && !isLoadingPack) {
      dispatch(load(packID));
    }

    if (pack) {
      if (!packHosts || packHosts.length !== pack.host_ids.length) {
        dispatch(hostActions.loadAll());
      }

      if (!packLabels || packLabels.length !== pack.label_ids.length) {
        dispatch(labelActions.loadAll());
      }
    }

    if (!size(scheduledQueries)) {
      dispatch(scheduledQueryActions.loadAll({ id: packID }));
    }

    if (!size(allQueries)) {
      dispatch(loadAll());
    }

    return false;
  }

  componentWillReceiveProps ({ dispatch, pack, packHosts, packLabels }) {
    if (pack && !this.props.pack) {
      if (!packHosts || packHosts.length !== pack.host_ids.length) {
        dispatch(hostActions.loadAll());
      }

      if (!packLabels || packLabels.length !== pack.label_ids.length) {
        dispatch(labelActions.loadAll());
      }
    }

    return false;
  }

  onCancelEditPack = () => {
    const { dispatch, isEdit, packID } = this.props;

    if (!isEdit) {
      return false;
    }

    return dispatch(push(`/packs/${packID}`));
  }

  onFetchTargets = (query, targetsResponse) => {
    const { targets_count: targetsCount } = targetsResponse;

    this.setState({ targetsCount });

    return false;
  }

  onSelectQuery = (query) => {
    const { allQueries } = this.props;
    const selectedQuery = find(allQueries, { id: Number(query) });
    this.setState({ selectedQuery });

    return false;
  }

  onToggleEdit = () => {
    const { dispatch, isEdit, packID } = this.props;

    if (isEdit) {
      return dispatch(push(`/packs/${packID}`));
    }

    return dispatch(push(`/packs/${packID}/edit`));
  }

  handlePackFormSubmit = (formData) => {
    const { dispatch, pack } = this.props;
    const { update } = packActions;

    return dispatch(update(pack, formData));
  }

  handleRemoveScheduledQueries = (scheduledQueryIDs) => {
    const { destroy } = scheduledQueryActions;
    const { dispatch } = this.props;

    const promises = scheduledQueryIDs.map((id) => {
      return dispatch(destroy({ id }));
    });

    return Promise.all(promises)
      .then(() => {
        dispatch(renderFlash('success', 'Scheduled queries removed'));
      });
  }

  handleConfigurePackQuerySubmit = (formData) => {
    const { create } = scheduledQueryActions;
    const { dispatch, packID } = this.props;
    const scheduledQueryData = {
      ...formData,
      snapshot: formData.logging_type === 'snapshot',
      pack_id: packID,
    };

    dispatch(create(scheduledQueryData))
      .then(() => {
        dispatch(renderFlash('success', 'Query scheduled!'));
      })
      .catch(() => {
        dispatch(renderFlash('error', 'Unable to schedule your query.'));
      });

    return false;
  }

  render () {
    const {
      handleConfigurePackQuerySubmit,
      handlePackFormSubmit,
      handleRemoveScheduledQueries,
      handleScheduledQueryFormSubmit,
      onCancelEditPack,
      onFetchTargets,
      onSelectQuery,
      onToggleEdit,
    } = this;
    const { targetsCount, selectedQuery } = this.state;
    const { allQueries, isEdit, isLoadingScheduledQueries, pack, packHosts, packLabels, scheduledQueries } = this.props;

    const packTargets = [...packHosts, ...packLabels];

    if (!pack || isLoadingScheduledQueries) {
      return false;
    }

    return (
      <div className={`${baseClass} has-sidebar`}>
        <div className={`${baseClass}__content`}>
          <EditPackFormWrapper
            className={`${baseClass}__pack-form body-wrap`}
            handleSubmit={handlePackFormSubmit}
            isEdit={isEdit}
            onCancelEditPack={onCancelEditPack}
            onEditPack={onToggleEdit}
            onFetchTargets={onFetchTargets}
            pack={pack}
            packTargets={packTargets}
            targetsCount={targetsCount}
          />
          <ScheduledQueriesListWrapper
            onRemoveScheduledQueries={handleRemoveScheduledQueries}
            onScheduledQueryFormSubmit={handleScheduledQueryFormSubmit}
            scheduledQueries={scheduledQueries}
          />
        </div>
        <ScheduleQuerySidePanel
          onConfigurePackQuerySubmit={handleConfigurePackQuerySubmit}
          allQueries={allQueries}
          onSelectQuery={onSelectQuery}
          selectedQuery={selectedQuery}
        />
      </div>
    );
  }
}

const mapStateToProps = (state, { params, route }) => {
  const entityGetter = stateEntityGetter(state);
  const isLoadingPack = state.entities.packs.loading;
  const { id: packID } = params;
  const pack = entityGetter.get('packs').findBy({ id: packID });
  const { entities: allQueries } = entityGetter.get('queries');
  const scheduledQueries = entityGetter.get('scheduled_queries').where({ pack_id: packID });
  const isLoadingScheduledQueries = state.entities.scheduled_queries.loading;
  const isEdit = route.path === 'edit';
  const packHosts = pack ? filter(state.entities.hosts.data, (host) => {
    return includes(pack.host_ids, host.id);
  }) : [];
  const packLabels = pack ? filter(state.entities.labels.data, (label) => {
    return includes(pack.label_ids, label.id);
  }) : [];

  return {
    allQueries,
    isEdit,
    isLoadingPack,
    isLoadingScheduledQueries,
    pack,
    packHosts,
    packID,
    packLabels,
    scheduledQueries,
  };
};

export default connect(mapStateToProps)(EditPackPage);
