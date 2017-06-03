import PricingPlan from './PricingPlan'
import { shallow } from 'enzyme'

it('renders without crashing', () => {
  const wrapper = shallow(<PricingPlan />)
  expect(wrapper.type()).toEqual('article')
  expect(wrapper.hasClass('pricing-column')).toBeTruthy()
});
