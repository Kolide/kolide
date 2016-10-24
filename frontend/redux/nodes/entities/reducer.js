import { combineReducers } from 'redux';

import hosts from './hosts/reducer';
import invites from './invites/reducer';
import targets from './targets/reducer';
import users from './users/reducer';

export default combineReducers({
  hosts,
  invites,
  targets,
  users,
});
