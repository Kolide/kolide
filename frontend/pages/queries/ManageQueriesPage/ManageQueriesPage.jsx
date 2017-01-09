import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { filter, get, includes, noop, pull } from 'lodash';
import { push } from 'react-router-redux';

import Button from 'components/buttons/Button';
import entityGetter from 'redux/utilities/entityGetter';
import InputField from 'components/forms/fields/InputField';
import Modal from 'components/modals/Modal';
import NumberPill from 'components/NumberPill';
import Icon from 'components/icons/Icon';
import PackInfoSidePanel from 'components/side_panels/PackInfoSidePanel';
import paths from 'router/paths';
import QueryDetailsSidePanel from 'components/side_panels/QueryDetailsSidePanel';
import QueriesList from 'components/queries/QueriesList';
import queryActions from 'redux/nodes/entities/queries/actions';
import queryInterface from 'interfaces/query';
import { renderFlash } from 'redux/nodes/notifications/actions';

const baseClass = 'manage-queries-page';

export class ManageQueriesPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    queries: PropTypes.arrayOf(queryInterface),
    selectedQuery: queryInterface,
  }

  static defaultProps = {
    dispatch: noop,
  };

  constructor (props) {
    super(props);

    this.state = {
      allQueriesChecked: false,
      checkedQueryIDs: [],
      queriesFilter: '',
      showModal: false,
    };
  }

  componentWillMount() {
    const { dispatch } = this.props;

    dispatch(queryActions.loadAll());

    return false;
  }

  onDeleteQueries = (evt) => {
    evt.preventDefault();

    const { checkedQueryIDs } = this.state;
    const { dispatch } = this.props;
    const { destroy } = queryActions;

    const promises = checkedQueryIDs.map((queryID) => {
      return dispatch(destroy({ id: queryID }));
    });

    return Promise.all(promises)
      .then(() => {
        dispatch(renderFlash('success', 'Queries successfully deleted.'));

        this.setState({ showModal: false });

        return false;
      })
      .catch(() => {
        dispatch(renderFlash('error', 'Something went wrong.'));

        this.setState({ showModal: false });

        return false;
      });
  }

  onCheckAllQueries = (shouldCheck) => {
    if (shouldCheck) {
      const queries = this.getQueries();
      const checkedQueryIDs = queries.map(query => query.id);

      this.setState({ allQueriesChecked: true, checkedQueryIDs });

      return false;
    }

    this.setState({ allQueriesChecked: false, checkedQueryIDs: [] });

    return false;
  }

  onCheckQuery = (checked, id) => {
    const { checkedQueryIDs } = this.state;
    const newCheckedQueryIDs = checked ? checkedQueryIDs.concat(id) : pull(checkedQueryIDs, id);

    this.setState({ allQueriesChecked: false, checkedQueryIDs: newCheckedQueryIDs });

    return false;
  }

  onFilterQueries = (queriesFilter) => {
    this.setState({ queriesFilter });

    return false;
  }

  onSelectQuery = (selectedQuery) => {
    const { dispatch } = this.props;
    const locationObject = {
      pathname: '/queries/manage',
      query: { selectedQuery: selectedQuery.id },
    };

    dispatch(push(locationObject));

    return false;
  }

  onToggleModal = () => {
    const { showModal } = this.state;

    this.setState({ showModal: !showModal });

    return false;
  }

  getQueries = () => {
    const { queriesFilter } = this.state;
    const { queries } = this.props;

    if (!queriesFilter) {
      return queries;
    }

    const lowerQueryFilter = queriesFilter.toLowerCase();

    return filter(queries, (query) => {
      if (!query.name) {
        return false;
      }

      const lowerQueryName = query.name.toLowerCase();

      return includes(lowerQueryName, lowerQueryFilter);
    });
  }

  goToNewQueryPage = () => {
    const { dispatch } = this.props;
    const { NEW_QUERY } = paths;

    dispatch(push(NEW_QUERY));

    return false;
  }

  goToEditQueryPage = (query) => {
    const { dispatch } = this.props;
    const { EDIT_QUERY } = paths;

    dispatch(push(EDIT_QUERY(query)));

    return false;
  }

  renderCTAs = () => {
    const { goToNewQueryPage, onToggleModal } = this;
    const btnClass = `${baseClass}__delete-queries-btn`;
    const checkedQueryCount = this.state.checkedQueryIDs.length;

    if (checkedQueryCount) {
      const queryText = checkedQueryCount === 1 ? 'Query' : 'Queries';

      return (
        <div className={`${baseClass}__ctas`}>
          <p className={`${baseClass}__query-count`}>{checkedQueryCount} {queryText} Selected</p>
          <Button
            className={btnClass}
            onClick={onToggleModal}
            variant="alert"
          >
            Delete
          </Button>
        </div>
      );
    }

    return (
      <Button variant="brand" onClick={goToNewQueryPage}>CREATE NEW QUERY</Button>
    );
  }

  renderModal = () => {
    const { onDeleteQueries, onToggleModal } = this;
    const { showModal } = this.state;

    if (!showModal) {
      return false;
    }

    return (
      <Modal onExit={onToggleModal}>
        <p>Are you sure you want to delete the selected queries?</p>
        <div>
          <Button onClick={onToggleModal} variant="inverse">Cancel</Button>
          <Button onClick={onDeleteQueries} variant="alert">Delete</Button>
        </div>
      </Modal>
    );
  }

  renderSidePanel = () => {
    const { goToEditQueryPage } = this;
    const { selectedQuery } = this.props;

    if (!selectedQuery) {
      // FIXME: Render QueryDetailsSidePanel when Fritz has completed the mock
      return <PackInfoSidePanel />;
    }

    return <QueryDetailsSidePanel onEditQuery={goToEditQueryPage} query={selectedQuery} />;
  }

  render () {
    const { checkedQueryIDs, queriesFilter } = this.state;
    const {
      getQueries,
      onCheckAllQueries,
      onCheckQuery,
      onSelectQuery,
      onFilterQueries,
      renderCTAs,
      renderModal,
      renderSidePanel,
    } = this;
    const { queries: allQueries, selectedQuery } = this.props;
    const queries = getQueries();
    const queriesCount = queries.length;
    const isQueriesAvailable = allQueries.length > 0;

    return (
      <div className={`${baseClass} has-sidebar`}>
        <div className={`${baseClass}__wrapper body-wrap`}>
          <p className={`${baseClass}__title`}>
            <NumberPill number={queriesCount} /> Queries
          </p>
          <div className={`${baseClass}__filter-and-cta`}>
            <div className={`${baseClass}__filter-queries`}>
              <InputField
                name="query-filter"
                onChange={onFilterQueries}
                placeholder="Filter Queries"
                value={queriesFilter}
              />
              <Icon name="search" />
            </div>
            {renderCTAs()}
          </div>
          <QueriesList
            checkedQueryIDs={checkedQueryIDs}
            isQueriesAvailable={isQueriesAvailable}
            onCheckAll={onCheckAllQueries}
            onCheckQuery={onCheckQuery}
            onSelectQuery={onSelectQuery}
            queries={queries}
            selectedQuery={selectedQuery}
          />
        </div>
        {renderSidePanel()}
        {renderModal()}
      </div>
    );
  }
}

const mapStateToProps = (state, { location }) => {
  const queryEntities = entityGetter(state).get('queries');
  const { entities: queries } = queryEntities;
  const selectedQueryID = get(location, 'query.selectedQuery');
  const selectedQuery = selectedQueryID && queryEntities.findBy({ id: selectedQueryID });

  return { queries, selectedQuery };
};

export default connect(mapStateToProps)(ManageQueriesPage);

