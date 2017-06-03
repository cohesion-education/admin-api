import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import Homepage from './Homepage'
import Auth from '../utils/Auth'

describe("<Homepage /> ", () => {
    let wrapper
    let _store = {
        dispatch: jest.fn(),
        subscribe: jest.fn(),
        getState: jest.fn(() =>
          ({
              auth:new Auth(),
              header:{
                title:'Test Title',
                subtitle:'Test Subtitle'
              },
              features:{
                title:'Features Title',
                subtitle:'Features Subtitle',
                highlights:[{
                  iconClassName: 'test-ico',
                  title:'Highlight Title',
                  description:'Highlight Description'
                }]
              },
              pricing:{
                title: 'Pricing Title',
                subtitle:'Pricing Subtitle',
                list:[
                  {
                    title: 'Price Title',
                    price: '$1 million',
                    duration: 'Forever'
                  }
                ]
              },
              testimonials:[{
                text: 'Super Awesome!',
                avatar:'test-avatar.png',
                name:'Satisfied Parent'
              }]
          })
        )
    }

    beforeAll(() => wrapper = mount(
      <Provider store={_store}>
        <Homepage />
      </Provider>
    ))

    afterEach(() => jest.resetAllMocks())

    it("renders without crashing", () => {
      expect(wrapper.find('Navbar').length).toBe(1)
      expect(wrapper.find('Header').length).toBe(1)
      expect(wrapper.find('FeatureList').length).toBe(1)
      expect(wrapper.find('FeatureDescriptionList').length).toBe(1)
      expect(wrapper.find('TestimonialList').length).toBe(1)
      expect(wrapper.find('PricingList').length).toBe(1)
      expect(wrapper.find('Footer').length).toBe(1)
    })
})
