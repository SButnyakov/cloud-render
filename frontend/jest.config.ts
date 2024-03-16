import type { Config } from "jest"

const config: Config = {
  verbose: true,
  transform: {
    "\\.tsx?$": "ts-jest",
    "\\.ts?$": "ts-jest",
    ".+\\.(png|jpg|svg|ttf|woff|woff2)$": "jest-transform-stub",
    "^.+\\.(js|jsx)$": "babel-jest"
  },
  moduleNameMapper: {
    "src/(.*)": "<rootDir>/src/$1",
    ".+\\.(css|styl|less|sass|scss)$": "identity-obj-proxy",
    "react-leaflet": "<rootDir>/src/mocks/reactLeafletMock.tsx"
  },
  globals: {
    "IS_REACT_ACT_ENVIRONMENT": true
  },
  collectCoverageFrom: [
    "src/**/*.{js,jsx,ts,tsx}",
  ],
  testPathIgnorePatterns: ["./node_modules/"],
  coveragePathIgnorePatterns: ["index.tsx", "store.ts", "test-utils.tsx", "declararion.d.ts", "RoutingMachine.tsx", "hooks.ts", "reactLeafletMock.ts"],
  testRegex: "(/__tests__/.*|(\\.|/)(test|spec))\\.tsx?$",
  moduleFileExtensions: ["ts", "tsx", "js", "json", "node"],
  collectCoverage: false,
  clearMocks: true,
  coverageDirectory: "coverage",
  testEnvironment: "allure-jest/jsdom",
}

export default config
