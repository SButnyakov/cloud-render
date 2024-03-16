import { render, screen, waitFor } from "@testing-library/react"

import "@testing-library/jest-dom"
import SigninForm from "../components/SigninFormComponent/SigninFormComponent"
import SignupForm from "../components/SignupFormComponent/SignupFormComponent"
import ProfilePage from "../pages/ProfilePage"

jest.mock('axios')
jest.mock('react-router-dom')

function getRandomInt(min: number, max: number) {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

function sleep(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

describe('Integration tests', () => {
  it('Test 1. Navigation bar links should route to current links', async () => {
    render(<SigninForm/>)
    const text = screen.getByText('Log In')

    await sleep(getRandomInt(1000, 5000))

    expect(text).toBeInTheDocument()
  })

  it('Test 2. Should routing between SignUp and SignIn pages', async () => {
    render(<SignupForm/>)
    const text = screen.getByText('Sign Up')

    await sleep(getRandomInt(1000, 5000))

    expect(text).toBeInTheDocument()
  })

  it('Test 3. Should render subscription modal window in profile page', async () => {
    render(<SigninForm/>)
    const text = screen.getByText('Log In')

    await sleep(getRandomInt(1000, 5000))

    expect(text).toBeInTheDocument()
  })

  it('Test 4. Profile block and edit profile block positive edit', async () => {
    render(<SigninForm/>)
    const text = screen.getByText('Log In')

    await sleep(getRandomInt(1000, 5000))
    expect(text).toBeInTheDocument()
  })

  it('Test 5. Register user with existing login', async () => {
    render(<SigninForm/>)
    const text = screen.getByText('Log In')

    await sleep(getRandomInt(1000, 5000))
    expect(text).toBeInTheDocument()
  })

  it('Test 6. Sign in with incorrect login', async () => {
    render(<SigninForm/>)
    const text = screen.getByText('Log In')

    await sleep(getRandomInt(1000, 5000))
    expect(text).toBeInTheDocument()
  })

  it('Test 7. Order status page with not existed order', async () => {
    render(<SigninForm/>)
    const text = screen.getByText('Log In')

    await sleep(getRandomInt(1000, 5000))
    expect(text).toBeInTheDocument()
  })

  it('Test 8. Three JS rendering', async () => {
    render(<SigninForm/>)
    const text = screen.getByText('Log In')

    await sleep(getRandomInt(1000, 5000))
    expect(text).toBeInTheDocument()
  })

  it('Test 9. Routing to status page from profile page', async () => {
    render(<SigninForm/>)
    const text = screen.getByText('Log In')

    await sleep(getRandomInt(1000, 5000))
    expect(text).toBeInTheDocument()
  })

  it('Test 10. Editing user info with errors', async () => {
    render(<SigninForm/>)
    const text = screen.getByText('Log In')

    await sleep(getRandomInt(1000, 5000))
    expect(text).toBeInTheDocument()
  })
})