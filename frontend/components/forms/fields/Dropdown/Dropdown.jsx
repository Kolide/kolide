import React, { Component, PropTypes } from 'react';
import radium from 'radium';
import { noop } from 'lodash';
import componentStyles from './styles';

class Dropdown extends Component {
  static propTypes = {
    containerStyles: PropTypes.object,
    options: PropTypes.arrayOf(PropTypes.shape({
      text: PropTypes.string,
      value: PropTypes.string,
    })),
    onSelect: PropTypes.func,
    style: PropTypes.object,
  };

  static defaultProps = {
    onSelect: noop,
  };

  onOptionClick = (evt) => {
    evt.preventDefault();

    const { onSelect } = this.props;

    onSelect(evt);

    return false;
  }

  renderOption = (option) => {
    const { disabled = false, value, text } = option;
    const { optionWrapperStyles } = componentStyles;

    return (
      <option key={value} style={optionWrapperStyles} value={value} disabled={disabled}>
        {text}
      </option>
    );
  }

  render () {
    const { containerStyles, options, style } = this.props;
    const { onOptionClick, renderOption } = this;
    const { selectWrapperStyles } = componentStyles;

    return (
      <div className="kolide-dropdown-wrapper" style={style}>
        <select className="kolide-dropdown" style={[selectWrapperStyles, containerStyles]} onChange={onOptionClick}>
          {options.map(option => {
            return renderOption(option);
          })}
        </select>
      </div>
    );
  }
}

export default radium(Dropdown);
