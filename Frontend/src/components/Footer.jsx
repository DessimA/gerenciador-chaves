import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faInstagram, faLinkedin } from '@fortawesome/free-brands-svg-icons'; // Import Instagram and LinkedIn icons
import './Footer.css';

function Footer() {
  return (
    <footer className="footer">
      <p>Desenvolvido por DessimA</p>
      <div className="social-links">
        <a href="https://www.instagram.com/dessim_dt" target="_blank" rel="noopener noreferrer">
          <FontAwesomeIcon icon={faInstagram} size="2x" /> {/* Instagram Icon */}
        </a>
        <a href="https://www.linkedin.com/in/dessim" target="_blank" rel="noopener noreferrer">
          <FontAwesomeIcon icon={faLinkedin} size="2x" /> {/* LinkedIn Icon */}
        </a>
      </div>
    </footer>
  );
}

export default Footer;