import React from 'react'
import PropTypes from 'prop-types'

class Testimonial extends React.Component {
  static propTypes = {
    text: PropTypes.string.isRequired,
    avatar: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired
  }

  static defaultProps =  {
    text:"",
    avatar: "",
    name: ""
  }

  render(){
    const {text, avatar, name} = this.props
    return(
      <div className="testimonial-box">
          <h4>{text}</h4>
          <img src={avatar} className="testi-user img-circle" alt="testimonials-user" />
          <p>- {name}</p>
      </div>
    )
  }
}

export default Testimonial;
