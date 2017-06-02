import Footer from './Footer'
import { shallow } from 'enzyme'

it('renders without crashing', () => {
  const wrapper = shallow(<Footer />)
  expect(wrapper.type()).toEqual('footer')
  expect(wrapper.hasClass('footer')).toBeTruthy()
});
