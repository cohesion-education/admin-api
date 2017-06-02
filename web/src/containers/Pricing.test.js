import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { Pricing } from './'


jest.mock('../views/PricingList')

describe("<Pricing /> Container ", () => {
    let wrapper
    let _store = {
        dispatch: jest.fn(),
        subscribe: jest.fn(),
        getState: jest.fn(() =>
          ({
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
              }
          })
        )
    }

    beforeAll(() => wrapper = mount(
      <Provider store={_store}>
        <Pricing />
      </Provider>
    ))

    afterEach(() => jest.resetAllMocks())

    it("renders PricingList with props", () => {
      expect(wrapper
        .find('PricingList')
        .props()
        .title
      ).toBe('Pricing Title')

      expect(wrapper
        .find('PricingList')
        .props()
        .subtitle
      ).toBe('Pricing Subtitle')

      expect(wrapper
        .find('PricingList')
        .props()
        .list
        .length
      ).toBe(1)
    })
})
