import React, { PropTypes } from 'react';
import { noop } from 'lodash';
import radium from 'radium';

import componentStyles from './styles';
import ElipsisMenu from '../../../components/buttons/ElipsisMenu';
import { humanMemory, humanUptime } from './helpers';

const {
  containerStyles,
  contentSeparatorStyles,
  disableIconStyles,
  elipsisChidrenWrapperStyles,
  elipsisPositionStyles,
  hostContentItemStyles,
  hostnameStyles,
  queryIconStyles,
  statusStyles,
  verticleRuleStyles,
} = componentStyles;
export const STATUSES = {
  online: 'ONLINE',
  offline: 'OFFLINE',
};

const HostDetails = ({ host, onQueryClick = noop, onDisableClick = noop }) => {
  const {
    hostname,
    ip,
    mac,
    memory,
    os_version: osVersion,
    platform,
    status,
    uptime,
  } = host;

  return (
    <div style={containerStyles(status)}>
      <ElipsisMenu positionStyles={elipsisPositionStyles}>
        <div style={elipsisChidrenWrapperStyles}>
          <div onClick={onQueryClick(host)} style={{ cursor: 'pointer', width: '60px' }}>
            <i className="kolidecon-query" style={queryIconStyles} />
            <div>Query</div>
          </div>
          <div style={verticleRuleStyles} />
          <div onClick={onDisableClick(host)} style={{ cursor: 'pointer', width: '60px' }}>
            <i className="kolidecon-ex" style={disableIconStyles} />
            <div>Disable</div>
          </div>
        </div>
      </ElipsisMenu>
      <div style={statusStyles(status)}>
        {status}
      </div>
      <p style={hostnameStyles}>{hostname}</p>
      <div style={contentSeparatorStyles}>
        <div>
          <span style={hostContentItemStyles}>{osVersion}</span>
        </div>
        <div>
          <span style={[hostContentItemStyles, { textTransform: 'capitalize' }]}>{platform}</span>
        </div>
        <div>
          <span style={hostContentItemStyles}>{humanMemory(memory)}</span>
          <span style={hostContentItemStyles}>{humanUptime(uptime)}</span>
        </div>
        <div>
          <span style={hostContentItemStyles}>{mac}</span>
        </div>
        <div>
          <span style={hostContentItemStyles}>{ip}</span>
        </div>
      </div>
      <div style={contentSeparatorStyles}>
        <div>
          <span style={[hostContentItemStyles, { textTransform: 'capitalize' }]}>Tags go here</span>
        </div>
      </div>
    </div>
  );
};

HostDetails.propTypes = {
  host: PropTypes.shape({
    hostname: PropTypes.string,
    ip: PropTypes.string,
    mac: PropTypes.string,
    memory: PropTypes.number,
    platform: PropTypes.string,
    uptime: PropTypes.number,
  }).isRequired,
  onDisableClick: PropTypes.func,
  onQueryClick: PropTypes.func,
};

export default radium(HostDetails);
