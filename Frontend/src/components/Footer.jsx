import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faInstagram, faLinkedin } from '@fortawesome/free-brands-svg-icons'; // Import Instagram and LinkedIn icons
import './Footer.css';

function Footer() {
  const githubProfileImageUrl = 'https://avatars.githubusercontent.com/u/60760405?v=4'; // GitHub profile image URL

  return (
    <footer className="footer">
      <img src={githubProfileImageUrl} alt="GitHub Profile" className="github-profile-image" />
      <p>Desenvolvido por DessimA</p>
      <div className="social-links">
        <a href="https://www.instagram.com/dessim_dt" target="_blank" rel="noopener noreferrer">
          <FontAwesomeIcon icon={faInstagram} size="2x" />
        </a>
        <a href="https://www.linkedin.com/in/dessim" target="_blank" rel="noopener noreferrer">
          <FontAwesomeIcon icon={faLinkedin} size="2x" />
        </a>
      </div>
    </footer>
  );
}

export default Footer;
