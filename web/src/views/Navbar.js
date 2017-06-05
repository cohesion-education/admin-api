import React from 'react'
import Auth from '../utils/Auth'

const auth = new Auth()

const Navbar = () =>
  <div className="navbar navbar-custom sticky navbar-fixed-top" role="navigation" id="sticky-nav">
    <div className="container">
      <div className="navbar-header">
        <button type="button" className="navbar-toggle" data-toggle="collapse" data-target=".navbar-collapse">
          <span className="sr-only">Toggle navigation</span>
          <span className="icon-bar"></span>
          <span className="icon-bar"></span>
          <span className="icon-bar"></span>
        </button>
        <a className="navbar-brand logo" href="/">
          <img src="/assets/images/cohesion-logo.png" alt="Cohesion Education"/>
        </a>
      </div>
      <div className="navbar-collapse collapse" id="navbar-menu">
        <ul className="nav navbar-nav navbar-right">
          <li className="active"><a href="#home" className="nav-link">Home</a></li>
          <li><a href="#features" className="nav-link">Features</a></li>
          <li><a href="#pricing" className="nav-link">Plans</a></li>
          <li><a href="/login" className="btn-login" onClick={auth.login}>Login</a></li>
          <li><a href="/register" onClick={auth.login} className="btn btn-white-bordered navbar-btn btn-login">Try for Free</a></li>
        </ul>
      </div>
    </div>
  </div>

export default Navbar
