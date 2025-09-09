import React from 'react';
import './AlertDialog.css'; // Create AlertDialog.css later

function AlertDialog({ message, isOpen, onClose }) {
  if (!isOpen) return null;

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <p>{message}</p>
        <div className="modal-actions">
          <button onClick={onClose} className="alert-button">OK</button>
        </div>
      </div>
    </div>
  );
}

export default AlertDialog;
