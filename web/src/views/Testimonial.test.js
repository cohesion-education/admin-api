import Testimonial from './Testimonial'
import { shallow } from 'enzyme'

it('renders without crashing', () => {
  const wrapper = shallow(<Testimonial />)
  expect(wrapper.type()).toEqual('div')
  expect(wrapper.hasClass('testimonial-box')).toBeTruthy()
});
