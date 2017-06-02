import React from 'react'

class FeatureDescriptionList extends React.Component {
  render(){
    return(
      <div>
        <section className="section p-t-0">
          <div className="container">
            <div className="row">
              <div className="col-sm-5">
                <div className="feat-description m-t-20">
                  <h4>Simple &amp; Organized</h4>
                  <p className="text-muted">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec non sem consequat, sollicitudin orci ut, placerat ipsum. Curabitur quis neque dolor. Nunc tincidunt eu tellus in fermentum. Donec vel mauris varius, commodo enim a, luctus ex. Integer et neque sagittis, porta diam eget, vestibulum dolor. Nulla odio risus, sodales sit amet placerat eu, luctus quis nisi. Vestibulum sed tempus nibh. Cras quis congue ante. Nunc ac lacus sed elit tincidunt mollis quis vel purus.</p>
                  <a href="" className="btn btn-custom">Learn More</a>
                </div>
              </div>
              <div className="col-sm-6 col-sm-offset-1">
                <img src="/assets/images/mac_1.png" alt="img" className="img-responsive m-t-20" />
              </div>
            </div>
          </div>
        </section>
        <section className="section">
          <div className="container">
            <div className="row">
              <div className="col-sm-6">
                <img src="/assets/images/mac_2.png" alt="img" className="img-responsive" />
              </div>
              <div className="col-sm-5 col-sm-offset-1">
                <div className="feat-description">
                  <h4>Related Lessons are easy to find</h4>
                  <p className="text-muted">Struggling to understand a lesson? No problem! We ensure that related lessons are linked to the current video to help you out.</p>
                  <a href="" className="btn btn-custom">Learn More</a>
                </div>
              </div>
            </div>
          </div>
        </section>
        <section className="section">
          <div className="container">
            <div className="row">
              <div className="col-sm-5">
                <div className="feat-description">
                  <h4>Greenlight / Redlight</h4>
                  <p className="text-muted">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec non sem consequat, sollicitudin orci ut, placerat ipsum. Curabitur quis neque dolor. Nunc tincidunt eu tellus in fermentum. Donec vel mauris varius, commodo enim a, luctus ex. Integer et neque sagittis, porta diam eget, vestibulum dolor. Nulla odio risus, sodales sit amet placerat eu, luctus quis nisi. Vestibulum sed tempus nibh. Cras quis congue ante. Nunc ac lacus sed elit tincidunt mollis quis vel purus.</p>
                  <a href="" className="btn btn-custom">Learn More</a>
                </div>
              </div>
              <div className="col-sm-6 col-sm-offset-1">
                <img src="/assets/images/mac_3.png" alt="img" className="img-responsive" />
              </div>
            </div>
          </div>
        </section>
      </div>
    )
  }
}

export default FeatureDescriptionList;
