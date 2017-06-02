import React from 'react'

class Footer extends React.Component {
  render(){
    return(
      <footer className="footer">
        <div className="container">
          <div className="row">
            <div className="col-lg-3 col-md-3">
                &copy; 2017 Cohesion Education, Inc.
            </div>
            <div className="col-lg-4 col-lg-offset-3 col-md-7">
              <ul className="nav navbar-nav">
                <li><a href="/about">About</a></li>
                <li><a href="/contact">Contact</a></li>
                <li><a href="/terms">Terms</a></li>
                <li><a href="/privacy">Privacy</a></li>
              </ul>
            </div>
            <div className="col-lg-2 col-md-2">
              <ul className="social-icons">
                <li><a href="#"><i className="fa fa-facebook"></i></a></li>
                <li><a href="#"><i className="fa fa-twitter"></i></a></li>
                <li><a href="#"><i className="fa fa-google-plus"></i></a></li>
              </ul>
            </div>
          </div>
        </div>
      </footer>
    )
  }
}

export default Footer;
