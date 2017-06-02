import Pricing from './Pricing'
import { shallow } from 'enzyme'

it('renders without crashing', () => {
  const wrapper = shallow(<Pricing />)
  expect(wrapper.type()).toEqual('article')
  expect(wrapper.hasClass('pricing-column')).toBeTruthy()
});
