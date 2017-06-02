import FeatureList from './FeatureList'
import { shallow } from 'enzyme'

it('renders without crashing', () => {
  const highlights = [
    {
      iconClassName: "fa-test",
      title: "test 1",
      description: "Lorem ipsum dolor"
    },
    {
      iconClassName: "fa-test",
      title: "Test 2",
      description: "Lorem ipsum dolor"
    }
  ]

  const wrapper = shallow(
    <FeatureList
      highlights={highlights} />
  )
  expect(wrapper.find('section.section').length).toBe(1)
  expect(wrapper.find('Feature').length).toBe(highlights.length)
});
