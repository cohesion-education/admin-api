import React from 'react'
import { connect } from 'react-redux'
import Navbar from '../views/Navbar'
import Header from './Header'
import Features from './Features'
import FeatureDescriptionList from '../views/FeatureDescriptionList'
import Testimonials from './Testimonials'
import Pricing from './Pricing'
import Footer from '../views/Footer'
import history from '../history'
import { fetchHomepage } from '../actions'
import '../css/fonts.css'
import '../css/font-awesome.css'
import '../css/homepage.css'


class Homepage extends React.Component {
  componentWillMount(){
    this.props.dispatch(fetchHomepage())
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
