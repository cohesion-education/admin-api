import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'


class DashboardTopBar extends React.Component {

  static propTypes = {
    picture: PropTypes.string.isRequired
  }

  static defaultProps =  {
    picture:''
  }

  render (){
    const { picture } = this.props
    console.log(`picture: ${picture}`)

    return(
      <div className="topbar">
        <div className="topbar-left">
          <div className="logo">
            <a href="#" className="open-left">
              <img className="icon-c-collapsed" src="/assets/images/cohesion-c-70x68.png" height="60"/>
              <img className="icon-c-logo" src="/assets/images/cohesion-logo.png" height="60"/>
            </a>
          </div>
        </div>

        <div className="navbar navbar-default" role="navigation">
          <div className="container">
            <div>
              <ul className="nav navbar-nav navbar-right pull-right">
                <li className="dropdown top-menu-item-xs">
                  <a href="#" className="dropdown-toggle profile waves-effect waves-light" data-toggle="dropdown" aria-expanded="true">
                    <img src={picture} alt="user-img" className="img-circle" />
                  </a>
                  <ul className="dropdown-menu">
                    <li><a href="#"><i className="ti-user m-r-10 text-custom"></i> Profile</a></li>
                    <li><a href="#"><i className="ti-settings m-r-10 text-custom"></i> Settings</a></li>
                    <li><a href="#"><i className="ti-lock m-r-10 text-custom"></i> Lock screen</a></li>
                    <li className="divider"></li>
                    <li><a href="/logout"><i className="ti-power-off m-r-10 text-danger"></i> Logout</a></li>
                  </ul>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    )
  }
}

export default connect(
  state => ({ ...state.profile })
)(DashboardTopBar)
