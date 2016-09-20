import React, { Component, PropTypes } from 'react';
import { isEqual, last } from 'lodash';
import componentStyles from './styles';
import navItems from './navItems';

class AdminSidePanel extends Component {
  static propTypes = {
    user: PropTypes.object,
  };

  constructor (props) {
    super(props);

    this.state = {
      activeTab: 'Hosts',
      activeSubItem: 'Add Hosts',
    };
  }

  setActiveTab = (activeTab) => {
    return (evt) => {
      evt.preventDefault();

      this.setState({ activeTab });
      return false;
    };
  }

  setActiveSubItem = (activeSubItem) => {
    return (evt) => {
      evt.preventDefault();

      this.setState({ activeSubItem });
      return false;
    };
  }

  renderHeader = () => {
    const {
      user: {
        enabled,
        username,
      },
    } = this.props;
    const {
      companyLogoStyles,
      headerStyles,
      orgNameStyles,
      usernameStyles,
      userStatusStyles,
    } = componentStyles;

    return (
      <header style={headerStyles}>
        <img
          alt="Company logo"
          src="http://2a0d7e1b.ngrok.io/assets/ge-logo-6d786c8e9079010a195f208d34ffb9d67e77ceff8c468c5c1e7fb739b086060f.png"
          style={companyLogoStyles}
        />
        <h1 style={orgNameStyles}>General Electric</h1>
        <div style={userStatusStyles(enabled)} />
        <h2 style={usernameStyles}>{username}</h2>
      </header>
    );
  }

  renderNavItem = (navItem, lastChild) => {
    const { activeTab } = this.state;
    const { icon, name, subItems } = navItem;
    const active = activeTab === name;
    const {
      iconStyles,
      navItemBeforeStyles,
      navItemNameStyles,
      navItemStyles,
      navItemWrapperStyles,
    } = componentStyles;
    const { renderSubItems, setActiveTab } = this;

    return (
      <div style={navItemWrapperStyles(lastChild)} key={`nav-item-${name}`}>
        {active && <div style={navItemBeforeStyles} />}
        <li
          onClick={setActiveTab(name)}
          style={navItemStyles(active)}
        >
          <div style={{ position: 'relative' }}>
            <i className={icon} style={iconStyles} />
            <span style={navItemNameStyles}>
              {name}
            </span>
          </div>
          {active && renderSubItems(subItems)}
        </li>
      </div>
    );
  }

  renderNavItems = () => {
    const { renderNavItem } = this;
    const { navItemListStyles } = componentStyles;

    return (
      <ul style={navItemListStyles}>
        {navItems.map((navItem, index, collection) => {
          const lastChild = isEqual(navItem, last(collection));
          return renderNavItem(navItem, lastChild);
        })}
      </ul>
    );
  }

  renderSubItem = (subItem) => {
    const { activeSubItem } = this.state;
    const { name, path } = subItem;
    const active = activeSubItem === name;
    const { setActiveSubItem } = this;
    const {
      subItemBeforeStyles,
      subItemStyles,
      subItemLinkStyles,
    } = componentStyles;

    return (
      <div
        key={`sub-item-${name}`}
        style={{ position: 'relative' }}
      >
        {active && <div style={subItemBeforeStyles} />}
        <li
          onClick={setActiveSubItem(name)}
          style={subItemStyles(active)}
        >
          <span to={path} style={subItemLinkStyles(active)}>{name}</span>
        </li>
      </div>
    );
  }

  renderSubItems = (subItems) => {
    const { subItemsStyles } = componentStyles;
    const { renderSubItem } = this;

    return (
      <ul style={subItemsStyles}>
        {subItems.map(subItem => {
          return renderSubItem(subItem);
        })}
      </ul>
    );
  }

  render () {
    const { navStyles } = componentStyles;
    const { renderHeader, renderNavItems } = this;

    return (
      <nav style={navStyles}>
        {renderHeader()}
        {renderNavItems()}
      </nav>
    );
  }
}

export default AdminSidePanel;
