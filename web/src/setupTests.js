import React from 'react'
global.React = React
// window.localStorage = {}
// console.groupCollapsed = jest.fn()
// console.log = jest.fn()
// console.groupEnd = jest.fn()

const localStorageMock = {
  getItem: jest.fn(),
  setItem: jest.fn(),
  clear: jest.fn()
};
global.localStorage = localStorageMock
