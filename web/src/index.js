import React from 'react'
import ReactDOM from 'react-dom'
import ReactGA from 'react-ga'
import Routes from './routes'
import { configureAnchors } from 'react-scrollable-anchor'
import 'bootstrap/dist/css/bootstrap.css'

configureAnchors({scrollDuration: 1500})
ReactGA.initialize(process.env.REACT_APP_GA_TRACKING_ID, { debug:true })

ReactDOM.render(
  <Routes />,
  document.getElementById('root')
);
