import Header from './Header'
import { shallow } from 'enzyme'

it('renders without crashing', () => {
  const wrapper = shallow(<Header />)
  expect(wrapper.type()).toEqual('section')
  expect(wrapper.hasClass('home')).toBeTruthy()
});
