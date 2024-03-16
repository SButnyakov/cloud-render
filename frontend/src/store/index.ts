import OrderStore from "./OrderStore"
import UserStore from "./UserStore"

// eslint-disable-next-line import/no-anonymous-default-export
export default {
  userStore: new UserStore(),
  orderStore: new  OrderStore()
}
