import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';

import configInterface from 'interfaces/config';
import Icon from 'components/icons/Icon';
import userInterface from 'interfaces/user';
import UserMenu from 'components/side_panels/SiteNavHeader/UserMenu';

class SiteNavHeader extends Component {
  static propTypes = {
    config: configInterface,
    onLogoutUser: PropTypes.func,
    onRedirectTo: PropTypes.func,
    user: userInterface,
  };

  constructor (props) {
    super(props);

    this.state = { userMenuOpened: false };
  }

  componentDidMount = () => {
    const { document } = global;
    const { closeUserMenu } = this;

    document.addEventListener('mousedown', closeUserMenu, false);
  }

  setHeaderNav = (ref) => {
    this.headerNav = ref;

    return false;
  }

  closeUserMenu = ({ target }) => {
    const { headerNav } = this;

    if (headerNav && !headerNav.contains(target)) {
      this.setState({ userMenuOpened: false });
    }

    return false;
  }

  toggleUserMenu = (evt) => {
    evt.preventDefault();

    const { userMenuOpened } = this.state;

    this.setState({ userMenuOpened: !userMenuOpened });

    return false;
  }

  render () {
    const {
      config: {
        org_name: orgName,
        org_logo_url: orgLogoURL,
      },
      onLogoutUser,
      onRedirectTo,
      user,
    } = this.props;

    const { userMenuOpened } = this.state;
    const { setHeaderNav, toggleUserMenu } = this;
    const { enabled, username } = user;

    const headerBaseClass = 'site-nav-header';

    const headerToggleClass = classnames(
      `${headerBaseClass}__button`,
      'button',
      'button--unstyled',
      { [`${headerBaseClass}__button--open`]: userMenuOpened }
    );

    const userStatusClass = classnames(
      `${headerBaseClass}__user-status`,
      { [`${headerBaseClass}__user-status--enabled`]: enabled }
    );

    return (
      <header className={headerBaseClass}>
        <button className={headerToggleClass} onClick={toggleUserMenu} ref={setHeaderNav}>
          <div className={`${headerBaseClass}__org`}>
            <img
              alt="Company logo"
              src={orgLogoURL}
              className={`${headerBaseClass}__logo`}
            />
            <h1 className={`${headerBaseClass}__org-name`}>{orgName}</h1>
            <div className={userStatusClass} />
            <h2 className={`${headerBaseClass}__username`}>{username}</h2>
            <Icon name="chevrondown" className={`${headerBaseClass}__org-chevron`} />
          </div>

          <UserMenu
            isOpened={userMenuOpened}
            onLogout={onLogoutUser}
            onRedirectTo={onRedirectTo}
            user={user}
          />
        </button>
      </header>
    );
  }
}

export default SiteNavHeader;
