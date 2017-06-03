import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import Features from './Features'

describe("<Features /> Container ", () => {
    let wrapper
    let _store = {
        dispatch: jest.fn(),
        subscribe: jest.fn(),
        getState: jest.fn(() =>
          ({
              features:{
                title:'Features Title',
                subtitle:'Features Subtitle',
                highlights:[{
                  iconClassName: 'test-ico',
                  title:'Highlight Title',
                  description:'Highlight Description'
                }]
              }
          })
        )
    }

    beforeAll(() => wrapper = mount(
      <Provider store={_store}>
        <Features />
      </Provider>
    ))

    afterEach(() => jest.resetAllMocks())

    it("renders FeatureList with props", () => {
      expect(wrapper
        .find('FeatureList')
        .props()
        .title
      ).toBe('Features Title')

      expect(wrapper
        .find('FeatureList')
        .props()
        .subtitle
      ).toBe('Features Subtitle')
    })
})
