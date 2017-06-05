import React from 'react'
import { connect } from 'react-redux'
import PropTypes from 'prop-types'
import Auth from '../utils/Auth'

const auth = new Auth()

class Header extends React.Component {
  static propTypes = {
    title: PropTypes.string.isRequired,
    subtitle: PropTypes.string.isRequired
  }

  static defaultProps =  {
    title:"",
    subtitle:""
  }

  render(){
    const {title, subtitle} = this.props

    return(
      <section className="home bg-img-1" id="home">
        <div className="bg-overlay"></div>
        <div className="container">
          <div className="row">
            <div className="col-sm-12 text-center">
              <div className="home-fullscreen">
                <div className="full-screen">
                  <div className="home-wrapper home-wrapper-alt">
                    <h2 className="font-light text-white">{title}</h2>
                    <h4 className="text-white">{subtitle}</h4>
                    <a href="/register" onClick={auth.login} className="btn btn-custom btn-login">Sign Up</a>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    )
  }
}

export default connect(
  state => ({ ...state.header })
)(Header)
