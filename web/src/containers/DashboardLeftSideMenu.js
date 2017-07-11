import React from 'react'
import { connect } from 'react-redux'
import { ListGroup, ListGroupItem } from 'react-bootstrap'

class DashboardLeftSideMenu extends React.Component {
  render (){
    return(
      <ListGroup>
        <ListGroupItem href="/videos">Videos</ListGroupItem>
      </ListGroup>
    )
  }
}

export default connect((state) => state)(DashboardLeftSideMenu)
