import React, { PropTypes } from 'react';

import Button from 'components/buttons/Button';
import hostInterface from 'interfaces/host';
import Icon from 'components/icons/Icon';
import PlatformIcon from 'components/icons/PlatformIcon';
import { humanMemory, humanUptime } from './helpers';

const baseClass = 'host-details';

export const STATUSES = {
  online: 'ONLINE',
  offline: 'OFFLINE',
};

const HostDetails = ({ host, onDestroyHost }) => {
  const {
    hostname,
    ip,
    mac,
    memory,
    os_version: osVersion,
    osquery_version: osqueryVersion,
    platform,
    status,
    uptime,
  } = host;

  return (
    <div className={`${baseClass} ${baseClass}--${status}`}>
      <span className={`${baseClass}__delete-host`}>
        <Button onClick={onDestroyHost(host)} variant="unstyled" title="Delete this host">
          <Icon name="trash" className={`${baseClass}__delete-host-icon`} />
        </Button>
      </span>

      <p className={`${baseClass}__status`}>{status}</p>

      <p className={`${baseClass}__hostname`}>{hostname}</p>

      <ul className={`${baseClass}__details-list`}>
        <li className={` ${baseClass}__detail ${baseClass}__detail--os`}>
          <PlatformIcon name={platform} className={`${baseClass}__icon`} />
          <span className={`${baseClass}__host-content`}>{osVersion}</span>
        </li>

        <li className={` ${baseClass}__detail ${baseClass}__detail--cpu`}>
          <Icon name="cpu" className={`${baseClass}__icon`} />
          <span className={`${baseClass}__host-content`}> 1 x 2.4Ghz</span>
        </li>

        <li className={` ${baseClass}__detail ${baseClass}__detail--osquery`}>
          <Icon name="osquery" className={`${baseClass}__icon`} />
          <span className={`${baseClass}__host-content`}>{osqueryVersion}</span>
        </li>

        <li className={` ${baseClass}__detail ${baseClass}__detail--memory`}>
          <Icon name="memory" className={`${baseClass}__icon`} />
          <span className={`${baseClass}__host-content`}>{humanMemory(memory)}</span>
        </li>

        <li className={` ${baseClass}__detail ${baseClass}__detail--uptime`}>
          <Icon name="uptime" className={`${baseClass}__icon`} />
          <span className={`${baseClass}__host-content`}>{humanUptime(uptime)}</span>
        </li>

        <li className={` ${baseClass}__detail ${baseClass}__detail--mac`}>
          <Icon name="mac" className={`${baseClass}__icon`} />
          <span className={`${baseClass}__host-content ${baseClass}__host-content--mono`}>{mac || '04:01:34:EA:54:01'}</span>
        </li>

        <li className={` ${baseClass}__detail ${baseClass}__detail--ip`}>
          <Icon name="world" className={`${baseClass}__icon`} />
          <span className={`${baseClass}__host-content ${baseClass}__host-content--mono`}>{ip || '104.236.116.77'}</span>
        </li>
      </ul>
    </div>
  );
};

HostDetails.propTypes = {
  host: hostInterface.isRequired,
  onDestroyHost: PropTypes.func.isRequired,
};

export default HostDetails;
