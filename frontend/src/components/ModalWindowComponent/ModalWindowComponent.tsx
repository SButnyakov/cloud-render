import React, { useEffect, useState } from "react";
import styles from './ModalWindowComponent.module.css'



const Modal = ({ onClose, onAction, actionType }: any) => {
  const [content, setContent] = useState('')
  const [header, setHeader] = useState('')

  useEffect(() => {
    console.log(actionType)
    initModalHeader()
  })

  const initModalHeader = () => {
    switch (actionType) {
      case 'quit':
        setHeader('Quit?')
        break;

      case 'deleteOrder':
        setHeader('Delete order?')
        break;
  
      case 'subscribe':
        setHeader('Subscription')
        break;
    
      default:
        break;
    }
  }
  return (
    <div className={styles.modalBackdrop}>
      <div className={actionType !== 'subscribe' ?styles.modalContent : styles.subModalContent}>
        

        {actionType === 'subscribe'
          ?
          <>
            <p className={styles.subHeader}>{header}</p>

            <div className={styles.subContent}>
              A subscription offers numerous benefits, including priority access, which means
              subscribers often get ahead in queues for services or new products. This priority
              access not only reduces wait times but also enhances the overall customer experience
              by offering exclusivity and additional perks like early product access and special discounts.
              Essentially, a subscription elevates the consumer experience, making it more efficient and rewarding,
              thereby offering considerable value beyond just the subscribed content or products.
            </div>

            <div className={styles.buttonsBlock}>
              <button onClick={onAction}>Buy</button>
              <button onClick={onClose}>Cancel</button>
            </div>
          </>
          :
          <>
            <p className={styles.header}>{header}</p>

            <div className={styles.buttonsBlock}>
              <button onClick={onAction}>Yes</button>
              <button onClick={onClose}>No</button>
            </div>
          </>
        }
        
      </div>
    </div>
  );
};

export default Modal

