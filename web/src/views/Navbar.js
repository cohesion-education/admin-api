import React from 'react'
import Auth from '../utils/Auth'
import { Navbar, Nav, NavItem } from 'react-bootstrap'
import injectSheet from 'react-jss'
import classNames from 'classnames'
import logo from '../images/cohesion-logo.png'

const auth = new Auth()

const styles = {
  navbar:{
    backgroundColor: '#ffffff',
    width: '100%',
    borderRadius: '0px',
    zIndex: '999',
    marginBottom: '0px',
    webkitTransition: 'all 0.5s ease-in-out',
    mozTransition: 'all 0.5s ease-in-out',
    oTransition: 'all 0.5s ease-in-out',
    transition: 'all 0.5s ease-in-out',
  },
  logo:{
    marginTop:'0.5em',
    height: '4.5em',
  },
  navLinks:{
    padding:'0',
    fontSize:'15px',
    margin:'5px 0px',
    paddingTop:'20px',
    paddingBottom:'10px',
    lineHeight:'20px',
    color:'#337ab7',
    textDecoration:'none',
    '&:hover':{
      backgroundColor: '#ffffff !important',
      color: '#5d9cec !important',
    }
  },
  registerButton:{
    backgroundColor: 'transparent !important',
    border: '2px solid #5d9cec !important',
    borderRadius: '50px',
    letterSpacing: '0.6px',
    marginTop:'20px',
    padding:'0px 0px',
  },
  '@media (min-width: 768px)': {
    navLinks: {
      paddingTop:'10px',
      paddingBottom:'10px',
    },
    registerButton:{
      padding:'0px',
    },
  }
}

const HomepageNavbar = ({classes}) =>
  <Navbar fluid fixedTop className={classes.navbar}>
    <Navbar.Header>
      <Navbar.Brand>
        <img src={logo} alt='Cohesion Education' className={classes.logo}/>
      </Navbar.Brand>
    </Navbar.Header>
    <Nav pullRight>
      <NavItem href='#features' className={classes.navLinks}>Features</NavItem>
      <NavItem href='#pricing' className={classes.navLinks}>Plans</NavItem>
      <NavItem href='/login' onClick={auth.login} className={classes.navLinks}>Login</NavItem>
      <NavItem href='/register' onClick={auth.login} className={classNames(classes.navLinks, classes.registerButton)}>Try for Free</NavItem>
    </Nav>
  </Navbar>

export default injectSheet(styles)(HomepageNavbar)
