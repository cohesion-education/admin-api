import React from 'react'
import { connect } from 'react-redux'
import DashboardTopBar from './DashboardTopBar'
import DashboardLeftSideMenu from './DashboardLeftSideMenu'
import DashboardFooter from '../views/DashboardFooter'
import Auth from '../utils/Auth'

import '../css/core.css'
import '../css/components.css'
// import '../css/icons.css'
import '../css/pages.css'
import '../css/responsive.css'

class Dashboard extends React.Component {
  auth = new Auth()

  componentDidMount(){
    //this.props.dispatch(fetchHomepage())
    console.log(`Dashboard.componentDidMount`)
    console.log(`is authenticated? ${this.auth.isAuthenticated()}`)
  }

  render (){
    return(
      <div id="wrapper">
        <DashboardTopBar />
        <DashboardLeftSideMenu />
        <div className="content-page">
          <div className="content">
            <div className="container">
              { /* TODO - display child */ }
              <div className="row">
                <div className="col-sm-12">
                  <h4 className="page-title">Welcome to Cohesion Education!</h4>
                  <p className="text-muted page-title-alt"></p>
                </div>
              </div>
            </div>
          </div>
          <DashboardFooter />
        </div>
      </div>
    )
  }
}

export default connect((state) => state)(Dashboard)
