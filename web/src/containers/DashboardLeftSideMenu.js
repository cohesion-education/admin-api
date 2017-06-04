import React from 'react'
import { connect } from 'react-redux'

class DashboardLeftSideMenu extends React.Component {
  render (){
    return(
      <div className="left side-menu">
        <div className="sidebar-inner slimscrollleft">
          <div id="sidebar-menu">
            <ul>
              <li className="text-muted menu-title">User Functions</li>
            </ul>
            <div className="clearfix"></div>
          </div>
          <div className="clearfix"></div>
        </div>
      </div>
    )
  }
}

export default connect((state) => state)(DashboardLeftSideMenu)
