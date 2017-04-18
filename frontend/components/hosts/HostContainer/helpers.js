import { filter, includes } from 'lodash';
import moment from 'moment';

const filterHosts = (hosts, label) => {
  if (!label) {
    return hosts;
  }

  if (label.type === 'status' && label.id === 'new') {
    return filter(hosts, h => moment().diff(moment(h.created_at)) <= moment.duration(24, 'hours'));
  }

  const { host_ids: hostIDs, slug, type } = label;

  switch (type) {
    case 'all':
      return hosts;
    case 'status':
      return filter(hosts, { status: slug });
    case 'platform': // Platform labels are implemented the same as custom labels
    case 'custom':
      return filter(hosts, h => includes(hostIDs, h.id));
    default:
      return hosts;
  }
};

const scrollToTop = () => {
  const { window } = global;

  const scrollStep = -window.scrollY / (500 / 15);
  const scrollInterval = setInterval(() => {
    if (window.scrollY !== 0) {
      window.scrollBy(0, scrollStep);
    } else {
      clearInterval(scrollInterval);
    }
  }, 15);
};

export default { filterHosts, scrollToTop };
