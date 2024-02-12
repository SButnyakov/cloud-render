import React from 'react'
import './App.css';
import { BrowserRouter } from 'react-router-dom'
import AppRouter from './components/AppRouter'
import NavBar from './components/NavBarComponent/NavBarComponent'
import { StoreProvider } from './hooks/useStore'

const styles = {
  backgroundImage: `url(${process.env.PUBLIC_URL + '/authBackground.png'})`,
  backgroundSize: "cover",
  height: "100vh",
  backgroundRepeat: "no-repeat",
}

function App() {
  return (
    <div style={styles}>
      <StoreProvider>
        <BrowserRouter>
          <NavBar/>
          <AppRouter/>
        </BrowserRouter>
      </StoreProvider>
    </div>
    
  )
}

export default App
