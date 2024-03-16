import React, { useEffect, useState } from 'react'
import cl from './NavBarComponent.module.css'
import { useNavigate } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { useStore } from '../../hooks/useStore'
import Modal from '../ModalWindowComponent/ModalWindowComponent'
import { getUser } from '../../http/UserApi'


const NavBarComponent = observer(() => {
  const currentActionModalType = 'quit'

  const route = useNavigate()
  const {userStore} = useStore()
  const [user, setUser] = useState('')

  const [isSelectVisible, setIsSelectVisible] = useState(false)

  const [isQuitModal, setIsQuitModal] = useState(false)

  useEffect(() => {
    
    if (userStore.isAuth) {
      getUser()
      .then(res => {
        setUser(res.login)
      })
    }
  }, [userStore.isAuth])

  const closeModalSub = () => {
    setIsQuitModal(false);
  };

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')
    localStorage.removeItem('orders')

    userStore.setIsAuth(false)
    setIsSelectVisible(false)

    setIsSelectVisible(false)

    closeModalSub()
    
    route('/landing')
  }


  return(
    <>
      <div className={cl.navBar}>
        {userStore.isAuth ? (<div
          className={cl.logoButton}
          onClick={() => route('/upload')}
        >
        </div>) : (<div
          className={cl.logoButton}
          onClick={() => route('/landing')}
        >
        </div>)}
        

        {userStore.isAuth ? 
          <div 
            className={cl.profileBlock}
            onClick={() => setIsSelectVisible(!isSelectVisible)}>
            <div className={cl.profileName}>
              {user}
            </div>
            <svg width="41" height="27" viewBox="0 0 41 27" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M37.5697 7.00013L19.0002 20.2355" stroke="#8C6CF3" strokeWidth="5"/>
              <line x1="3.45102" y1="6.96418" x2="22.0205" y2="20.1995" stroke="#8C6CF3" strokeWidth="5"/>
            </svg>
          </div> 
          : 
          <div className={cl.actionBarUnauth}>
            <button onClick={() => route('/signin')}>Login</button>
            <button onClick={() => route('/signup')}>Register</button>
          </div>
        }

        {isSelectVisible && 
        (
          <div className={cl.profileMenu}>
            <div onClick={() => {
              route('/profile')
              setIsSelectVisible(false)}}>
                Profile
            </div>
            <div onClick={() => setIsQuitModal(true)}>Sign Out</div>
          </div>)}
      </div>

      {isQuitModal && <Modal onClose={closeModalSub} onAction={handleLogout} actionType={currentActionModalType}/>}
    </>
  )
})

export default NavBarComponent