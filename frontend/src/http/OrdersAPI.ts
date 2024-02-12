import { $authHost, $host } from "."
import { AxiosError } from "axios"
import { Order } from "../store/OrderStore"

export const getOrders = async (): Promise<Order[]> => {
  try {
    const {data} = await $authHost.get('/orders')
    return data.orders
  }
  catch (e) {
    console.log(e)
  }
  
  return []
}

export const deleteOrder = async (orderId: number) => {
  const {data} = await $authHost.post(`/orders/${orderId}/delete`)
  return data
}

export const getOrder = async (orderId: number): Promise<Order>  => {
  const {data} = await $authHost.get(`/orders/${orderId}`)
  return data
}
