import Feature from './Feature'
import { shallow } from 'enzyme'

it('renders without crashing', () => {
  const wrapper = shallow(
    <Feature
      iconClassName='test-ico-classname'
      title='test-title'
      description='test-description'/>
  )
  expect(wrapper.type()).toEqual('div')
  expect(wrapper.hasClass('features-box')).toBeTruthy()
  expect(wrapper.contains(<i className='fa test-ico-classname' />)).toBeTruthy()
  expect(wrapper.contains(<h4>test-title</h4>)).toBeTruthy()
  expect(wrapper.contains(<p className='text-muted'>test-description</p>)).toBeTruthy()
});
