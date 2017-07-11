import React from 'react'
import { connect } from 'react-redux'
import { Grid, Row, Col, Clearfix } from 'react-bootstrap'
import DashboardTopBar from './DashboardTopBar'
import DashboardLeftSideMenu from './DashboardLeftSideMenu'
import DashboardFooter from '../views/DashboardFooter'
import VideoPlayer from '../views/VideoPlayer'

const styles = {
  containerFluid:{
    padding:0,
  },
  header:{
    position:'fixed',
    zIndex:1000,
  },
  main:{
    paddingTop:'7em',
    paddingRight:'10',
    paddingLeft:'10',
    paddingBottom:'1em',
  },
  contentArea:{
    marginTop:'-0.7em',
    paddingBottom:'2em',
  },
}

class Dashboard extends React.Component {

  render (){
    const videoJsOptions = {
      autoplay: false,
      controls: true,
      sources: [{
        src: 'http://www.sample-videos.com/video/mp4/240/big_buck_bunny_240p_1mb.mp4',
        type: 'video/mp4'
      }]
    }

    return(
      <Grid fluid style={styles.containerFluid}>
        <Row>
          <Col xs={12}>
            <DashboardTopBar/>
          </Col>
        </Row>
        <Row style={styles.main}>
          <Col sm={3} xsHidden>
            <DashboardLeftSideMenu />
          </Col>
          <Col sm={9} xs={12} style={styles.contentArea}>
            <h4>Welcome to Cohesion Education!</h4>
            <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aenean malesuada interdum bibendum. Aliquam varius tristique cursus. Suspendisse scelerisque ornare neque, nec semper risus aliquam in. Ut aliquet neque enim, et ultricies justo feugiat eget. Ut ut maximus dolor, ut rutrum nisl. Nam a pharetra diam. Maecenas nec bibendum tellus. Sed vitae magna nisi. Aliquam sed rutrum ligula, nec iaculis nisi. Curabitur eleifend lacus ut leo tincidunt blandit.</p>
            <VideoPlayer { ...videoJsOptions } />
            
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            <DashboardFooter />
          </Col>
        </Row>
      </Grid>
    )
  }
}


export default connect((state) => state)(Dashboard)
