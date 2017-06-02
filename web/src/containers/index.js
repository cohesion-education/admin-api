import { connect } from 'react-redux'
import Navbar from '../views/Navbar'
import Header from '../views/Header'
import FeatureList from '../views/FeatureList'
import TestimonialList from '../views/TestimonialList'
import PricingList from '../views/PricingList'

// function mapStateToProps(state){
//   console.log(`mapping state to props`)
//   Object.keys(state).map(key => console.log(`${key} = ${state[key]}; type: ${typeof state[key]}`))
//   return { ...state }
// }

export const Nav = connect(
  state => ({ auth:state.auth }),
  dispatch => ({})
)(Navbar)


export const HomepageHeader = connect(
  state => ({ ...state.header }),
  dispatch => ({})
)(Header)

export const Features = connect(
  state => ({ ...state.features }),
  dispatch => ({})
)(FeatureList)

export const Testimonials = connect(
  state => ({ ...state.testimonials }),
  dispatch => ({})
)(TestimonialList)

export const Pricing = connect(
  state => ({ ...state.pricing }),
  dispatch => ({})
)(PricingList)
