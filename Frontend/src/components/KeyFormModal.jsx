import React, { useState } from 'react';
import './KeyFormModal.css'; // Create KeyFormModal.css later
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSave, faTimes } from '@fortawesome/free-solid-svg-icons';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

function KeyFormModal({ isOpen, onClose, onSave, showAlert }) {
  const [apartment, setApartment] = useState('');
  const [keyType, setKeyType] = useState('');

  if (!isOpen) return null;

  const handleAddKey = async (e) => {
    e.preventDefault();
    if (!apartment || !keyType) {
      showAlert('Por favor, preencha o número do apartamento e o tipo da chave.');
      return;
    }
    try {
      await fetch(`${API_URL}/keys`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ apartment_number: apartment, key_type: keyType }),
      });
      onSave(); // Call the callback to refresh keys in App.jsx
      setApartment('');
      setKeyType('');
      onClose(); // Close modal after saving
    } catch (error) {
      console.error('Erro ao adicionar chave:', error);
      showAlert('Erro ao adicionar chave. Tente novamente.');
    }
  };

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <h2>Adicionar Nova Chave</h2>
        <form onSubmit={handleAddKey} className="key-form-modal">
          <input
            type="text"
            placeholder="Nº do Apartamento"
            value={apartment}
            onChange={(e) => setApartment(e.target.value)}
          />
          <input
            type="text"
            placeholder="Tipo (ex: apartamento, garagem)"
            value={keyType}
            onChange={(e) => setKeyType(e.target.value)}
          />
          <div className="modal-actions">
            <button type="submit" className="confirm-button"><FontAwesomeIcon icon={faSave} /> Salvar</button>
            <button type="button" onClick={onClose} className="cancel-button"><FontAwesomeIcon icon={faTimes} /> Cancelar</button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default KeyFormModal;
