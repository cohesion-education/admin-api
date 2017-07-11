import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import DashboardTopBar from './DashboardTopBar'

describe("<DashboardTopBar /> ", () => {
    let wrapper
    let _store = {
        dispatch: jest.fn(),
        subscribe: jest.fn(),
        getState: jest.fn(() =>
          ({
            profile:{
              picture:'test-picture'
            }
          })
        )
    }

    beforeAll(() => wrapper = mount(
      <Provider store={_store}>
        <DashboardTopBar />
      </Provider>
    ))

    afterEach(() => jest.resetAllMocks())

    it("renders without crashing", () => {
      expect(wrapper
        .find('img')
        .length
      ).toBe(2)
    })
})
