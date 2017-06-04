import React from 'react'

const InternalServerError = () =>
  <div className="wrapper-page">
    <div className="ex-page-content text-center">
      <div className="text-error">
        <span className="text-primary">5</span>
        <i className="ti-face-sad text-pink"></i>
        <i className="ti-face-sad text-info"></i>
      </div>

      <h2>Internal Server Error.</h2>

      <p className="text-muted">
        Why not try refreshing your page? or you can contact <a href="#">support</a>
      </p>
      <a className="btn btn-default waves-effect waves-light" href="/"> Return Home</a>
    </div>
  </div>

export default InternalServerError
