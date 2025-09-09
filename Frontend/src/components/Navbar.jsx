import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faGithub } from '@fortawesome/free-brands-svg-icons'; // Import GitHub icon
import './Navbar.css';

function Navbar() {
  return (
    <nav className="navbar">
      <a href="https://github.com/DessimA" target="_blank" rel="noopener noreferrer" className="github-link">
        <FontAwesomeIcon icon={faGithub} size="2x" /> {/* GitHub Icon */}
      </a>
    </nav>
  );
}

export default Navbar;