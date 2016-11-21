import React, { Component, PropTypes } from 'react';
import { noop } from 'lodash';
import AceEditor from 'react-ace';
import classnames from 'classnames';

import hostHelpers from 'components/hosts/HostDetails/helpers';
import targetInterface from 'interfaces/target';

const baseClass = 'target-details';

class TargetDetails extends Component {
  static propTypes = {
    target: targetInterface,
    className: PropTypes.string,
    handleBackToResults: PropTypes.func,
  };

  static defaultProps = {
    handleBackToResults: noop,
  };

  renderHost = () => {
    const { className, handleBackToResults, target } = this.props;
    const {
      display_text: displayText,
      ip,
      mac,
      memory,
      osqueryVersion,
      osVersion,
      platform,
      status,
    } = target;
    const hostBaseClass = 'host-target';
    const isOnline = status === 'online';
    const isOffline = status === 'offline';
    const statusClassName = classnames(
      `${hostBaseClass}__status`,
      { [`${hostBaseClass}__status--is-online`]: isOnline },
      { [`${hostBaseClass}__status--is-offline`]: isOffline },
    );

    return (
      <div className={`${hostBaseClass} ${className}`}>
        <button className={`button button--unstyled ${hostBaseClass}__back`} onClick={handleBackToResults}>
          <i className="kolidecon kolidecon-chevronleft" />Back
        </button>

        <p className={`${hostBaseClass}__display-text`}>
          <i className={`${hostBaseClass}__icon kolidecon-fw kolidecon-single-host`} />
          <span>{displayText}</span>
        </p>
        <p className={statusClassName}>
          {isOnline && <i className={`${hostBaseClass}__icon ${hostBaseClass}__icon--online kolidecon-fw kolidecon-success-check`} />}
          {isOffline && <i className={`${hostBaseClass}__icon ${hostBaseClass}__icon--offline kolidecon-fw kolidecon-offline`} />}
          <span>{status}</span>
        </p>
        <table className={`${baseClass}__table`}>
          <tbody>
            <tr>
              <th>IP Address</th>
              <td>{ip}</td>
            </tr>
            <tr>
              <th>MAC Address</th>
              <td>{mac}</td>
            </tr>
            <tr>
              <th>Platform</th>
              <td>
                <i className={hostHelpers.platformIconClass(platform)} />
                <span className={`${hostBaseClass}__platform-text`}>{platform}</span>
              </td>
            </tr>
            <tr>
              <th>Operating System</th>
              <td>{osVersion}</td>
            </tr>
            <tr>
              <th>Osquery Version</th>
              <td>{osqueryVersion}</td>
            </tr>
            <tr>
              <th>Memory</th>
              <td>{hostHelpers.humanMemory(memory)}</td>
            </tr>
          </tbody>
        </table>
        <div className={`${hostBaseClass}__labels-wrapper`}>
          <p className={`${hostBaseClass}__labels-header`}>
            <i className={`${hostBaseClass}__icon kolidecon-fw kolidecon-label`} />
            <span>Labels</span>
          </p>
          <ul className={`${hostBaseClass}__labels-list`}>
            <li>Engineering</li>
            <li>DevOps</li>
            <li>ElCapDev</li>
            <li>Workstation</li>
          </ul>
        </div>
      </div>
    );
  }

  renderLabel = () => {
    const { handleBackToResults, className, target } = this.props;
    const {
      count,
      description,
      display_text: displayText,
      online,
      query,
    } = target;
    const labelBaseClass = 'label-target';

    return (
      <div className={`${labelBaseClass} ${className}`}>
        <button className={`button button--unstyled ${labelBaseClass}__back`} onClick={handleBackToResults}>
          <i className="kolidecon kolidecon-chevronleft" />Back
        </button>

        <p className={`${labelBaseClass}__display-text`}>
          <i className={`${labelBaseClass}__icon kolidecon-fw kolidecon-label`} /> {displayText}
        </p>

        <p className={`${labelBaseClass}__hosts`}>
          <span className={`${labelBaseClass}__hosts-count`}>{count} HOSTS</span>
          <span className={`${labelBaseClass}__hosts-online`}> ({online}% ONLINE)</span>
        </p>

        <p className={`${labelBaseClass}__description`}>{description}</p>

        <div className={`${labelBaseClass}__editor`}>
          <AceEditor
            editorProps={{ $blockScrolling: Infinity }}
            mode="kolide"
            minLines={4}
            maxLines={4}
            fontSize={13}
            name="label-query"
            readOnly
            setOptions={{ wrap: true }}
            showGutter={false}
            showPrintMargin={false}
            theme="kolide"
            value={query}
            width="100%"
          />
        </div>
      </div>
    );
  }

  render () {
    const { target } = this.props;

    if (!target) {
      return false;
    }

    const { target_type: targetType } = target;
    const { renderHost, renderLabel } = this;

    if (targetType === 'labels') {
      return renderLabel();
    }

    return renderHost();
  }
}

export default TargetDetails;
