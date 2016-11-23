import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import formDataInterface from 'interfaces/registration_form_data';

class ConfirmationPage extends Component {
  static propTypes = {
    className: PropTypes.string,
    formData: formDataInterface,
    handleSubmit: PropTypes.func,
  };

  onSubmit = (evt) => {
    evt.preventDefault();

    const { handleSubmit } = this.props;

    return handleSubmit();
  }

  render () {
    const {
      className,
      formData: {
        email,
        full_name: fullName,
        kolide_server_url: kolideWebAddress,
        org_name: orgName,
        username,
      },
    } = this.props;
    const { onSubmit } = this;

    return (
      <div className={className}>
        <i className="kolidecon kolidecon-success-check" />
        <table>
          <caption>Administrator Configuration</caption>
          <tbody>
            <tr>
              <th>Full Name:</th>
              <td>{fullName}</td>
            </tr>
            <tr>
              <th>Username:</th>
              <td>{username}</td>
            </tr>
            <tr>
              <th>Email:</th>
              <td>{email}</td>
            </tr>
            <tr>
              <th>Organization:</th>
              <td>{orgName}</td>
            </tr>
            <tr>
              <th>Kolide URL:</th>
              <td>{kolideWebAddress}</td>
            </tr>
          </tbody>
        </table>
        <Button
          onClick={onSubmit}
          text="Submit"
          variant="gradient"
        />
      </div>
    );
  }
}

export default ConfirmationPage;

