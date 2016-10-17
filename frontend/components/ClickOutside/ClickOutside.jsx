import React, { Component } from 'react';
import ReactDOM from 'react-dom'
import { noop } from 'lodash';
import radium from 'radium';

import { handleClickOutside } from './helpers';

export default (WrappedComponent, { onOutsideClick = noop }) => {
  class ClickOutside extends Component {
    componentDidMount () {
      const { componentInstance } = this;
      const clickHandler = onOutsideClick(componentInstance);
      const componentNode = ReactDOM.findDOMNode(componentInstance);

      this.handleAction = handleClickOutside(clickHandler, componentNode);

      global.document.addEventListener('mousedown', this.handleAction);
      global.document.addEventListener('touchStart', this.handleAction);
    }

    componentWillUnmount () {
      global.document.removeEventListener('mousedown', this.handleAction);
      global.document.removeEventListener('touchStart', this.handleAction);
    }

    setInstance = (instance) => {
      this.componentInstance = instance;
    }

    render () {
      const { setInstance } = this;
      return <WrappedComponent {...this.props } ref={setInstance}/>
    }
  }

  return radium(ClickOutside);
};
