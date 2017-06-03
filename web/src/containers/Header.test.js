import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import Header from './Header'

describe("<Header /> Container ", () => {
    let wrapper
    let _store = {
        dispatch: jest.fn(),
        subscribe: jest.fn(),
        getState: jest.fn(() =>
          ({
              header:{
                title:'Test Title',
                subtitle:'Test Subtitle'
              }
          })
        )
    }

    beforeAll(() => wrapper = mount(
      <Provider store={_store}>
        <Header />
      </Provider>
    ))

    afterEach(() => jest.resetAllMocks())

    it("renders header with props", () => {
      expect(wrapper
        .find('Header')
        .props()
        .title
      ).toBe('Test Title')

      expect(wrapper
        .find('Header')
        .props()
        .subtitle
      ).toBe('Test Subtitle')
    })
})
