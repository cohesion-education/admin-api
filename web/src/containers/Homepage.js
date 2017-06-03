import React from 'react'
import { connect } from 'react-redux'
import Navbar from '../views/Navbar'
import Header from './Header'
import Features from './Features'
import FeatureDescriptionList from '../views/FeatureDescriptionList'
import Testimonials from './Testimonials'
import Pricing from './Pricing'
import Footer from '../views/Footer'
import Auth from '../utils/Auth'
import { fetchHomepage } from '../actions'
import '../css/fonts.css'
import '../css/font-awesome.css'
import '../css/homepage.css'

class Homepage extends React.Component {
  auth = new Auth()

  componentDidMount(){
    this.props.dispatch(fetchHomepage())

    console.log(`Homepage.componentDidMount`)
    console.log(`is authenticated? ${this.auth.isAuthenticated()}`)
    this.auth.getProfile((err, profile) => {
      if(err){
        console.log(`failed to get profile: ${err}`)
        return
      }
      console.log(`profile: ${JSON.stringify(profile)}`)
    })
  }

  render (){
    return(
      <div>
        <Navbar />
        <Header />
        <Features />
        <FeatureDescriptionList />
        <Testimonials />
        <Pricing />
        <Footer />
      </div>
    )
  }
}

export default connect((state) => state)(Homepage)
