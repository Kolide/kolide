import './public-path';

import ReactDOM from 'react-dom';

import routes from './router';
import './index.scss';

if (typeof window !== 'undefined') {
  const { document } = global;
  const app = document.getElementById('app');

  ReactDOM.render(routes, app);
}
