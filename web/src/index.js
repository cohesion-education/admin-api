import React from 'react'
import ReactDOM from 'react-dom'
import ReactGA from 'react-ga'
import Routes from './routes'
import { configureAnchors } from 'react-scrollable-anchor'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap/dist/css/bootstrap-theme.css'

configureAnchors({scrollDuration: 1500})
ReactGA.initialize('UA-92236743-1', { debug:true })
ReactGA.set({ page: window.location.pathname + window.location.search })
ReactGA.pageview(window.location.pathname + window.location.search)

ReactDOM.render(
  <Routes />,
  document.getElementById('root')
);
