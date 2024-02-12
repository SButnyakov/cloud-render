import React from "react";
import { useNavigate } from "react-router-dom";

import styles from './styles/LandingPage.module.css'

const gifStyle = {
  backgroundImage: `url(${process.env.PUBLIC_URL + '/landingVideo.gif'})`,
  height: "400px",
  width: "1500px",
  backgroundSize: "cover",
  boxShadow: "0 2px 8px rgba(0, 0, 0, 0.26)",
  margin: "10px",
  borderRadius: "10px"
}

const LandingPage = () => {
  const route = useNavigate()

  return(
    <div className={styles.pageBlock}>
      <div style={gifStyle}></div>
      <div className={styles.textBlock}>
        <div className={styles.mainText}>
        <h2>Introducing Render Service: Where Imagination Meets Reality</h2>
      <p>Welcome to Render Service, the pinnacle of architectural visualization and artistic rendering services. We breathe life into blueprints and dreams, transforming them into stunning, photo-realistic visuals. Our service is crafted for architects, designers, real estate professionals, and visionaries aiming to showcase their projects not just as structures, but as experiences waiting to unfold.</p>

      <h3>Why Choose Our Render Service?</h3>
      <ul>
        <li><b>Unparalleled Quality:</b> With a relentless focus on detail, texture, and light, our renders speak volumes of our dedication to quality. Each project is a masterpiece, ensuring your vision is realized with utter perfection.</li>
        <li><b>Innovative Technology:</b> Leveraging cutting-edge rendering software and techniques, we stay at the forefront of digital visualization. Our team continually explores new tools and methodologies to ensure your projects stand out in the digital age.</li>
        <li><b>Creative Expertise:</b> At the heart of our Render Service are creative professionals with a deep understanding of architecture and design. Our team collaborates closely with clients, breathing life into ideas through a blend of technical prowess and artistic insight.</li>
        <li><b>Customized Solutions:</b> Recognizing that every project has its unique essence, we offer tailored rendering services. Whether it's residential, commercial, interior, or landscape visualization, our solutions cater to your specific needs and aspirations.</li>
        <li><b>Timely Delivery:</b> We understand the value of your time. Our efficient workflow and dedicated team ensure that your projects are delivered with precision, on schedule, and beyond expectations.</li>
      </ul>

        </div>
      </div>
    </div>
  )
}

export default LandingPage