import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { Button, Navbar, Nav, NavItem, NavDropdown, MenuItem } from 'react-bootstrap'
import logo from '../images/cohesion-logo.png'

const styles = {
  nav:{
    backgroundColor: '#e2e2e2',
    borderColor: '#cecece',
  },
  logo:{
    height:'4.5em',
  },
  navLinks:{
    marginTop: '18px',
    fontSize: '1.2em',
    marginRight: '5px',
  },
  profilePicture:{
    border: '2px solid #edf0f0',
    height: '48px',
    width: '48px',
    cursor:'pointer',
  }
}

class DashboardTopBar extends React.Component {

  static propTypes = {
    picture: PropTypes.string.isRequired
  }

  static defaultProps =  {
    picture:''
  }

  render (){
    const { picture } = this.props

    return(
      <Navbar fluid style={styles.nav} fixedTop>
        <Navbar.Header>
          <Navbar.Brand>
            <img src={logo} alt='Cohesion Education' style={styles.logo}/>
          </Navbar.Brand>
        </Navbar.Header>
        <Nav pullRight>
          <Navbar.Text>
            <Navbar.Link href="/profile">
              <img src={picture} alt="user-img" className="img-circle" style={styles.profilePicture}/>
            </Navbar.Link>
          </Navbar.Text>
          <Button href="/logout" bsStyle="primary" style={styles.navLinks}>Logout</Button>
        </Nav>
      </Navbar>
    )
  }
}

export default connect(
  state => ({ ...state.profile })
)(DashboardTopBar)
