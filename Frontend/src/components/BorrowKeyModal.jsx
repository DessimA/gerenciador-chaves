import React, { useState } from 'react';
import './BorrowKeyModal.css'; // Create BorrowKeyModal.css later
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCheck, faTimes } from '@fortawesome/free-solid-svg-icons';

function BorrowKeyModal({ isOpen, onClose, onConfirmBorrow, showAlert }) {
  const [borrowerName, setBorrowerName] = useState('');

  if (!isOpen) return null;

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!borrowerName.trim()) {
      showAlert('Por favor, digite o nome do retirante.');
      return;
    }
    onConfirmBorrow(borrowerName);
    setBorrowerName(''); // Clear input after submission
    onClose();
  };

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <h2>Emprestar Chave</h2>
        <form onSubmit={handleSubmit} className="borrow-form-modal">
          <input
            type="text"
            placeholder="Nome do Retirante"
            value={borrowerName}
            onChange={(e) => setBorrowerName(e.target.value)}
          />
          <div className="modal-actions">
            <button type="submit" className="confirm-button"><FontAwesomeIcon icon={faCheck} /> Confirmar</button>
            <button type="button" onClick={onClose} className="cancel-button"><FontAwesomeIcon icon={faTimes} /> Cancelar</button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default BorrowKeyModal;
