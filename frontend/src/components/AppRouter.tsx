import React, { useContext } from "react"
import { Navigate, Route, Routes } from "react-router-dom"
import AuthPage from "../pages/AuthPage"
import RegisterPage from "../pages/RegisterPage"
import UploadPage from "../pages/UploadPage"
import LandingPage from "../pages/LandingPage"
import ProfilePage from "../pages/ProfilePage"
import OrderStatusPage from "../pages/OrderStatusPage"
import { observer } from "mobx-react-lite"
import { useStore } from "../hooks/useStore"

export const AppRouter = observer(() => {

  const {userStore} = useStore()

  return(
    <Routes>
      {userStore.isAuth && (
        <>
          <Route path="/profile" element={<ProfilePage/>}/>
          <Route path="/order/:id" element={<OrderStatusPage/>}/>
          <Route path="/upload" element={<UploadPage/>}/>
        </>
      )}

      {!userStore.isAuth && (
        <>
        <Route path="/signin" element={<AuthPage/>}/>
        <Route path="/signup" element={<RegisterPage/>}/>
        <Route path="/landing" element={<LandingPage/>}/>
        </>
      )}

      {/* <Route
        path="*" 
        element={<Navigate to="/" replace/>}
      /> */}
    </Routes>
  )
})

export default AppRouter