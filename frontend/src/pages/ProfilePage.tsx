import React, { useEffect, useMemo, useRef, useState } from "react";
import { editUser, getUser } from "../http/UserApi";
import { observer } from "mobx-react-lite";
import { User } from "../store/UserStore";

import styles from './styles/ProfilePage.module.css'
import { Order } from "../store/OrderStore";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import Modal from "../components/ModalWindowComponent/ModalWindowComponent";
import { deleteOrder, getOrders } from "../http/OrdersAPI";
import { subscribe } from "../http/SubscribeAPI";

type OrderProps = {
  order: Order;
  route: (data: string) => void;
  handleDownloadImage: (data: string) => void;
  handleDelete: (data: number) => void;
}

type FormData = {
  login: string;
  email: string;
  password: string;
  passwordRepeat: string;
};

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  const day = String(date.getUTCDate()).padStart(2, '0')
  const month = String(date.getUTCMonth() + 1).padStart(2, '0')
  const year = date.getUTCFullYear()
  return `${day}.${month}.${year}`
}


const OrderRow = (props: OrderProps) => {
  const {order, route, handleDelete, handleDownloadImage} = props
  const initStatusIcon = () => {
    switch (order.status) {
      case 'in queue':
        return '🚦'
      case 'in progress':
        return '⏳'
      case 'success':
        return '✅'
      default:
        return '[ERROR CHARACTER]'
    }
      
  }
  return (
    <div className={styles.orderRow}>
      <div className={styles.rowCell}>{order.filename}</div>
      <div className={styles.rowCell} onClick={() => route(`/order/${order.id}`)}>{order.date}</div>
      <div className={styles.orderActions}>
        <div className={styles.orderAction}>{initStatusIcon()}</div>
        {order.status === 'success' ? <div className={styles.orderAction} onClick={() => handleDownloadImage(order.downloadLink.String)}>💾</div> : null}
        <div className={styles.orderAction} onClick={() => handleDelete(order.id)}>🗑️</div> 
      </div>
      
    </div>
  )
};

const OrderTable = ({ orders, handleDelete, route, handleDownloadImage }: any) => {
  return (
    <div className={styles.ordersContainer}>
      <h2>Orders</h2>
      <div className={styles.orderTable}>
        {orders.map((order: any) => (
          <OrderRow
            order={order}
            handleDelete={handleDelete}
            route={route}
            handleDownloadImage={handleDownloadImage}
            key={order.id}
          />
        ))}
      </div>
    </div>
  );
};


const ProfilePage = observer(() => {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<FormData>();
  
  const route = useNavigate()

  const [currentActionModalType, setCurrentActionModalType] = useState('')
  const [user, setUser] = useState({} as User)
  const [orders, setOrders] = useState([] as Order[])
  const [isEditSatus, setIsEditStatus] = useState(false)

  const [extraErrors, setExtraErrors] = useState('')

  const [modalVisibleSub, setSubModalVisible] = useState(false);
  const [modalVisibleOrder, setOrderModalVisible] = useState(false);

  const [deletedOrder, setDeletedOrder] = useState<number>()

  const intervalIdRef = useRef(null as unknown as NodeJS.Timer)

  const checkOrdersStatus = () => {
    getOrders()
      .then(res => {
        setOrders(res.map((order) => ({...order, date: formatDate(order.date)})))
      })
    
    const checkOrderStatusIntervalId = setInterval(() => {
      getOrders()
        .then(res => {
          setOrders(res.map((order) => ({...order, date: formatDate(order.date)})))
        })
    }, 3000)

    intervalIdRef.current = checkOrderStatusIntervalId
  }

  useEffect(() => {
    getUser()
      .then(res => {
        setUser({email: res.email, login: res.login, exparationDate: res.exparationDate})
      })
    
    checkOrdersStatus()

    return () => {
      clearInterval(intervalIdRef.current)
    }
  }, [])

  const openModalSub = () => {
    setCurrentActionModalType('subscribe')
    setSubModalVisible(true);
  };

  const closeModalSub = () => {
    setSubModalVisible(false);
  };

  const openModalOrder = (id: number) => {
    setCurrentActionModalType('deleteOrder')
    setOrderModalVisible(true);

    setDeletedOrder(id)
  };

  const closeModalOrder = () => {
    setOrderModalVisible(false);
  };

  const buySubscriptionAction = async () => {
    await subscribe()
      .then(() => closeModalSub())
    
    await getUser()
      .then(res => {
        setUser({email: res.email, login: res.login, exparationDate: res.exparationDate})
      })
  }

  const handleDeleteOrder = async () => {
    setCurrentActionModalType('deleteOrder')
    await deleteOrder(deletedOrder as number)

    await getOrders()
      .then(res => {
        setOrders(res.map((order) => ({...order, date: formatDate(order.date)})))
      })
  
    closeModalOrder();
  }
  
  let errorList = useMemo(() => {
    const errorMessages = []
    const errorsObject = {...errors}

    for (let [_, value] of Object.entries(errorsObject)) {
      errorMessages.push(value)
    }

    return errorMessages
  }, [errors])

  useEffect(() => {

  })

  /* TODO: Доделать нормальный вывод ошибок */
  const handleSetIsEditStatus = (isEdit: boolean) => {
    setIsEditStatus(isEdit)
    if (!isEdit) {
      setExtraErrors('')
    }
  }

  const onSubmit = async (data: FormData) => {
    if (data.password.localeCompare(data.passwordRepeat) !== 0) {
      console.log('password', data.password)
      console.log('repeat_password', data.passwordRepeat)

      console.log(data.password.localeCompare(data.passwordRepeat))
      setExtraErrors('Passwords must be equal!')
    }
    else {
      await editUser(data.email, data.login, data.password)
      
      const user = await getUser()

      setUser({email: user.email, login: user.login, exparationDate: user.exparationDate})
      
      setIsEditStatus(false)
    } 
  }

  /* TODO: Тут дождаться когда будет все ок с ссылкой для скачив */

  const handleDownloadImage = (downloadLink: string) => {
    const link = document.createElement('a');
    link.href = downloadLink;
    link.download = ''
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }

  return(
    <div className={styles.profilePage}>

      {!isEditSatus 
        ? 
        (<div className={styles.formContainer}>
          <form className={styles.profileForm} onSubmit={handleSubmit(onSubmit)}>
            <h2>Your Profile</h2>
            <div className={styles.formGroup}>
              Email: <div>{user.email}</div>
            </div>
            <div className={styles.formGroup}>
              Login: <div>{user.login}</div>
            </div>
            <div className={styles.formGroup} style={{visibility: user.exparationDate ? 'visible' : 'hidden'}}>
              Sub expire date: <div>{user.exparationDate}</div>
            </div>
            <div className={styles.formActions} style={{marginTop: '125px'}}>
              <button type="button" className={styles.cancelButton} onClick={() => openModalSub()}>
                Buy subscription
              </button>
              <button type="submit" className={styles.saveButton} onClick={() => handleSetIsEditStatus(true)}>
                Edit
              </button>
            </div>
          </form>
        </div>) 
        : 
        (<div className={styles.formContainer}>
          <form className={styles.profileForm} onSubmit={handleSubmit(onSubmit)}>
            <h2>Your Profile</h2>
            <div className={styles.formGroup}>
              Email: <input type="text" 
                    id="email"
                    {...register('email', {
                      required: 'Input email',
                      value: user.email,
                      pattern: {
                        value: /^[\w-]+(\.[\w-]+)*@([\w-]+\.)+[a-zA-Z]{2,7}$/,
                        message:
                          'Wrong email format!',
                      },
                    })} />
            </div>
            <div className={styles.formGroup}>
              Login: <input type="text" 
                    id="login"
                    {...register('login', {
                      required: 'Input Login',
                      value: user.login,
                      pattern: {
                        value: /^[a-zA-Z0-9]{4,15}$/,
                        message:
                          'The login must contain Latin letters and have a maximum length of 15 characters',
                      },
                    })} />
            </div>
            <div className={styles.formGroup}>
              Password:
              <input
                type="text" 
                id="password"
                {...register('password', {
                  required: 'Input password'
                })}
              />
            </div>
            <div className={styles.formGroup}>
              Confirm Password:
              <input
                type="text" 
                id="password_repeat"
                {...register('passwordRepeat', {
                  required: 'Repeat password'
                })}
              />
            </div>
            <div className={styles.errorBlockMessage}>
                  {errorList.map(el => {
                  return <p>{el && el.message}</p> 
                  })}
    
                  {extraErrors ?? null}
                </div>
            <div className={styles.formActions}>
              <button type="button" className={styles.cancelButton} onClick={() => handleSetIsEditStatus(false)}>
                Cancel
              </button>
              <button type="submit" className={styles.saveButton}>
                Save
              </button>
            </div>
          </form>
        </div>)
      }


      <OrderTable 
        orders={orders}
        handleDelete={openModalOrder}
        handleDownloadImage={handleDownloadImage}
        route={route}
      />

      {modalVisibleSub && <Modal onClose={closeModalSub} onAction={buySubscriptionAction} actionType={currentActionModalType}/>}
      {modalVisibleOrder && <Modal onClose={closeModalOrder} onAction={handleDeleteOrder} actionType={currentActionModalType}/>}

    </div>
  )
})

export default ProfilePage