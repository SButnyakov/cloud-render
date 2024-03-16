import { createContext, useContext } from "react"
import allStores from "../store"

const StoreContext = createContext(allStores)

export const StoreProvider = ({ children } : any) => {
  return(
    <StoreContext.Provider value={allStores}>
      {children}
    </StoreContext.Provider>
  )
}

export const useStore = () => useContext(StoreContext)
