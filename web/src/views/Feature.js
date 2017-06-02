import React from 'react'
import PropTypes from 'prop-types'

class Feature extends React.Component {
  static propTypes = {
    iconClassName: PropTypes.string.isRequired,
    title: PropTypes.string.isRequired,
    description: PropTypes.string.isRequired
  }

  static defaultProps =  {
    iconClassName: "",
    title:"",
    description:""
  }

  render(){
    const {iconClassName, title, description} = this.props

    return(
      <div className="features-box">
        <i className={`fa ${iconClassName}`}></i>
        <h4>{title}</h4>
        <p className="text-muted">{description}</p>
      </div>
    )
  }
}

export default Feature;
