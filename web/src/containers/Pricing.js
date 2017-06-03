import React from 'react'
import { connect } from 'react-redux'
import PropTypes from 'prop-types'
import ScrollableAnchor from 'react-scrollable-anchor'
import PricingPlan from '../views/PricingPlan'

class PricingList extends React.Component {
  static propTypes = {
    title: PropTypes.string.isRequired,
    subtitle: PropTypes.string.isRequired,
    list: PropTypes.array.isRequired
  }

  static defaultProps =  {
    title:"",
    subtitle: "",
    list: []
  }

  render(){
    const {title, subtitle, list} = this.props
    return(
      <ScrollableAnchor id='pricing'>
        <section className="section" key="pricing">
          <div className="container">
            <div className="row">
              <div className="col-sm-12 text-center">
                <h3 className="title">{title}</h3>
                <p className="text-muted sub-title">{subtitle}</p>
              </div>
            </div>

            <div className="row">
              <div className="col-lg-10 col-lg-offset-1">
                <div className="row">
                  {list.map((pricingPlan, i) =>
                    <PricingPlan key={i} {...pricingPlan} />
                  )}
                </div>
              </div>
            </div>
          </div>
        </section>
      </ScrollableAnchor>
    )
  }
}

export default connect(
  state => ({ ...state.pricing })
)(PricingList)
