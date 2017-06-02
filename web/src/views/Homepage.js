import React from 'react'
import { Nav, HomepageHeader, Features, Testimonials, Pricing } from '../containers'
import FeatureDescriptionList from './FeatureDescriptionList'
import Footer from './Footer'
import '../css/fonts.css'
import '../css/font-awesome.css'
import '../css/homepage.css'

const Homepage = () =>
  <div>
    <Nav />
    <HomepageHeader />
    <Features />
    <FeatureDescriptionList />
    <Testimonials />
    <Pricing />
    <Footer />
  </div>


export default Homepage
