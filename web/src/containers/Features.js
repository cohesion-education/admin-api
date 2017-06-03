import React from 'react'
import { connect } from 'react-redux'
import PropTypes from 'prop-types'
import ScrollableAnchor from 'react-scrollable-anchor'
import Feature from '../views/Feature'

class FeatureList extends React.Component {
  static propTypes = {
    title: PropTypes.string.isRequired,
    subtitle: PropTypes.string.isRequired,
    highlights: PropTypes.array.isRequired
  }

  static defaultProps =  {
    title:"",
    subtitle:"",
    highlights:[]
  }

  render(){
    const {title, subtitle, highlights} = this.props

    return(
      <ScrollableAnchor id='features'>
        <section className='section' key='features'>
          <div className="container">
            <div className="row">
              <div className="col-sm-12 text-center">
                <h3 className="title">{title}</h3>
                <p className="text-muted sub-title">{subtitle}</p>
              </div>
            </div>
            <div className="row">
              {highlights.map((feature, i) =>
                <div key={i} className="col-sm-4">
                  <Feature {...feature}/>
                </div>
              )}
            </div>
          </div>
        </section>
      </ScrollableAnchor>
    )
  }
}

export default connect(
  state => ({ ...state.features }),
  dispatch => ({})
)(FeatureList)
