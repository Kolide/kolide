import React, { Component, PropTypes } from 'react';
import radium from 'radium';

import Avatar from '../../../../components/Avatar';
import componentStyles from './styles';
import Dropdown from '../../../../components/forms/fields/Dropdown';
import EditUserForm from '../../../../components/forms/Admin/EditUserForm';
import { userStatusLabel } from './helpers';

class UserBlock extends Component {
  static propTypes = {
    currentUser: PropTypes.object,
    invite: PropTypes.bool,
    onEditUser: PropTypes.func,
    onSelect: PropTypes.func,
    user: PropTypes.object,
  };

  static userActionOptions = (currentUser, user) => {
    const disableActions = currentUser.id === user.id;
    const userEnableAction = user.enabled
      ? { disabled: disableActions, text: 'Disable Account', value: 'disable_account' }
      : { text: 'Enable Account', value: 'enable_account' };
    const userPromotionAction = user.admin
      ? { disabled: disableActions, text: 'Demote User', value: 'demote_user' }
      : { text: 'Promote User', value: 'promote_user' };

    return [
      { text: 'Actions...', value: '' },
      userEnableAction,
      userPromotionAction,
      { text: 'Require Password Reset', value: 'reset_password' },
      { text: 'Modify Details', value: 'modify_details' },
    ];
  };

  constructor (props) {
    super(props);

    this.state = {
      isEdit: false,
    };
  }

  onToggleEditing = (evt) => {
    evt.preventDefault();

    const { isEdit } = this.state;

    this.setState({
      isEdit: !isEdit,
    });

    return false;
  }

  onEditUserFormSubmit = (updatedUser) => {
    const { user, onEditUser } = this.props;

    this.setState({
      isEdit: false,
    });

    return onEditUser(user, updatedUser);
  }

  onUserActionSelect = ({ target }) => {
    const { onSelect, user } = this.props;
    const { value: action } = target;

    if (action === 'modify_details') {
      this.setState({
        isEdit: true,
      });

      return false;
    }

    return onSelect(user, action);
  }

  render () {
    const { currentUser, invite, user } = this.props;
    const {
      avatarStyles,
      nameStyles,
      userDetailsStyles,
      userEmailStyles,
      userHeaderStyles,
      userLabelStyles,
      usernameStyles,
      userPositionStyles,
      userStatusStyles,
      userStatusWrapperStyles,
      userWrapperStyles,
    } = componentStyles(invite);
    const {
      admin,
      email,
      enabled,
      name,
      position,
      username,
    } = user;
    const statusLabel = userStatusLabel(user, invite);
    const userLabel = admin ? 'Admin' : 'User';
    const userActionOptions = UserBlock.userActionOptions(currentUser, user);
    const { isEdit } = this.state;
    const { onEditUserFormSubmit, onToggleEditing } = this;

    if (isEdit) {
      return <EditUserForm onCancel={onToggleEditing} onSubmit={onEditUserFormSubmit} user={user} />;
    }

    return (
      <div style={userWrapperStyles}>
        <div style={userHeaderStyles(admin)}>
          <span style={nameStyles}>{name}</span>
        </div>
        <div style={userDetailsStyles}>
          <Avatar user={user} style={avatarStyles} />
          <div style={userStatusWrapperStyles}>
            <span style={userLabelStyles(admin)}>{userLabel}</span>
            <span style={userStatusStyles(enabled)}>{statusLabel}</span>
            <div style={{ clear: 'both' }} />
          </div>
          <p style={usernameStyles}>{username}</p>
          <p style={userPositionStyles}>{position}</p>
          <p style={userEmailStyles}>{email}</p>
          {!invite && <Dropdown
            options={userActionOptions}
            initialOption={{ text: 'Actions...' }}
            onSelect={this.onUserActionSelect}
          />}
        </div>
      </div>
    );
  }
}

export default radium(UserBlock);
