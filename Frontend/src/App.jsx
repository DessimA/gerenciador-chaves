import { useState, useEffect } from 'react';
import './App.css';
// import backgroundImage from './assets/background.png'; // No longer needed here
import Navbar from './components/Navbar';
import Footer from './components/Footer';
import ConfirmModal from './components/ConfirmModal';
import AlertDialog from './components/AlertDialog';
import KeyFormModal from './components/KeyFormModal';
import BorrowKeyModal from './components/BorrowKeyModal'; // Import the new borrow modal

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faHandshake, faUndo, faEdit, faTrash, faPlusCircle, faCheck, faTimes } from '@fortawesome/free-solid-svg-icons'; // Import faCheck, faTimes

function App() {
  const [keys, setKeys] = useState([]);
  const [editingKey, setEditingKey] = useState(null);
  const [editApartment, setEditApartment] = useState('');
  const [editKeyType, setEditKeyType] = useState('');
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [keyToDelete, setKeyToDelete] = useState(null);

  const [showAlertModal, setShowAlertModal] = useState(false);
  const [alertMessage, setAlertMessage] = useState('');

  const [showKeyFormModal, setShowKeyFormModal] = useState(false);
  const [showBorrowModal, setShowBorrowModal] = useState(false);
  const [keyToBorrow, setKeyToBorrow] = useState(null);

  const [showReturnModal, setShowReturnModal] = useState(false); // New state for return modal visibility
  const [keyToReturn, setKeyToReturn] = useState(null); // New state to store key ID to return
  const [borrowerToConfirmReturn, setBorrowerToConfirmReturn] = useState(''); // New state for borrower name

  const API_URL = import.meta.env.VITE_API_URL || '/api';

  const fetchKeys = async () => {
    try {
      const response = await fetch(`${API_URL}/keys`);
      const data = await response.json();
      setKeys(data || []);
      console.log('Keys after fetch:', data);
    } catch (error) {
      console.error('Erro ao buscar chaves:', error);
    }
  };

  useEffect(() => {
    fetchKeys();
    console.log('Keys on mount:', keys);
  }, []);

  // Function to show alert modal
  const showAlert = (message) => {
    setAlertMessage(message);
    setShowAlertModal(true);
  };

  // Function to close alert modal
  const closeAlert = () => {
    setShowAlertModal(false);
    setAlertMessage('');
  };

  // handleAddKey logic moved to KeyFormModal, this function will now just open the modal
  const openKeyFormModal = () => {
    setShowKeyFormModal(true);
  };

  const closeKeyFormModal = () => {
    setShowKeyFormModal(false);
  };

  // Function to open borrow modal
  const openBorrowModal = (id) => {
    setKeyToBorrow(id);
    setShowBorrowModal(true);
  };

  const closeBorrowModal = () => {
    setShowBorrowModal(false);
    setKeyToBorrow(null);
  };

  const handleConfirmBorrow = async (borrowerName) => {
    try {
      await fetch(`${API_URL}/keys/${keyToBorrow}/borrow`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ borrower_name: borrowerName }),
      });
      fetchKeys();
    } catch (error) {
      console.error('Erro ao emprestar chave:', error);
      showAlert('Erro ao emprestar chave. Tente novamente.');
    }
  };

  // Function to open return confirmation modal
  const confirmReturn = (keyId, borrowerName) => {
    setKeyToReturn(keyId);
    setBorrowerToConfirmReturn(borrowerName);
    setShowReturnModal(true);
  };

  // Function to handle actual return after modal confirmation
  const handleReturnConfirmed = async () => {
    try {
      await fetch(`${API_URL}/keys/${keyToReturn}/return`, {
        method: 'PUT',
      });
      fetchKeys();
    } catch (error) {
      console.error('Erro ao devolver chave:', error);
      showAlert('Erro ao devolver chave. Tente novamente.');
    } finally {
      setShowReturnModal(false);
      setKeyToReturn(null);
      setBorrowerToConfirmReturn('');
    }
  };

  // Function to cancel return
  const handleReturnCancelled = () => {
    setShowReturnModal(false);
    setKeyToReturn(null);
    setBorrowerToConfirmReturn('');
  };

  // Function to open the modal
  const confirmDelete = (id) => {
    setKeyToDelete(id);
    setShowDeleteModal(true);
  };

  // Function to handle actual deletion after modal confirmation
  const handleDeleteConfirmed = async () => {
    console.log('Attempting to delete key:', keyToDelete);
    try {
      await fetch(`${API_URL}/keys/${keyToDelete}`, {
        method: 'DELETE',
      });
      fetchKeys();
      console.log('Key deleted successfully.');
    } catch (error) {
      console.error('Erro ao excluir chave:', error);
      showAlert('Erro ao excluir chave. Tente novamente.');
    } finally {
      setShowDeleteModal(false);
      setKeyToDelete(null);
      console.log('Delete modal closed.');
    }
  };

  // Function to cancel deletion
  const handleDeleteCancelled = () => {
    setShowDeleteModal(false);
    setKeyToDelete(null);
  };

  const handleEditKey = (key) => {
    setEditingKey(key.id);
    setEditApartment(key.apartment_number);
    setEditKeyType(key.key_type);
  };

  const handleUpdateKey = async (e, id) => {
    e.preventDefault();
    if (!editApartment || !editKeyType) {
      showAlert('Por favor, preencha o número do apartamento e o tipo da chave.');
      return;
    }
    try {
      await fetch(`${API_URL}/keys/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ apartment_number: editApartment, key_type: editKeyType }),
      });
      setEditingKey(null);
      fetchKeys();
    } catch (error) {
      console.error('Erro ao atualizar chave:', error);
      showAlert('Erro ao atualizar chave. Tente novamente.');
    }
  };

  return (
    <>
      <div
        className="background-blur"
        // style={{ backgroundImage: `url(${backgroundImage})` }} // Removed
      ></div>
      <div className="content-container">
        <Navbar />
        <h1>Gerenciador de Chaves</h1>
        {/* Button to open the KeyFormModal */}
        <button onClick={openKeyFormModal} className="add-key-button">
          <FontAwesomeIcon icon={faPlusCircle} />
        </button>

        <div className="key-list">
          {keys.map((key) => (
            <div key={key.id} className={`key-item ${key.status}`}>
              {editingKey === key.id ? (
                <form onSubmit={(e) => handleUpdateKey(e, key.id)} className="edit-key-form">
                  <input
                    type="text"
                    value={editApartment}
                    onChange={(e) => setEditApartment(e.target.value)}
                  />
                  <input
                    type="text"
                    value={editKeyType}
                    onChange={(e) => setEditKeyType(e.target.value)}
                  />
                  <button type="submit"><FontAwesomeIcon icon={faEdit} /></button>
                  <button type="button" onClick={() => setEditingKey(null)}><FontAwesomeIcon icon={faUndo} /></button>
                </form>
              ) : (
                <> 
                  <div className="key-info">
                    <span>Apto: {key.apartment_number}</span>
                    <span>Tipo: {key.key_type}</span>
                    <span>Status: {key.status}</span>
                    {key.status === 'emprestada' && <span>Retirante: {key.borrower_name}</span>}
                  </div>
                  <div className="key-actions">
                    {key.status === 'disponivel' ? (
                      <button onClick={() => openBorrowModal(key.id)}><FontAwesomeIcon icon={faHandshake} /></button>
                    ) : (
                      <button onClick={() => confirmReturn(key.id, key.borrower_name)}><FontAwesomeIcon icon={faUndo} /></button>
                    )}
                    <button onClick={() => handleEditKey(key)}><FontAwesomeIcon icon={faEdit} /></button>
                    <button onClick={() => confirmDelete(key.id)}><FontAwesomeIcon icon={faTrash} /></button>
                  </div>
                </>
              )}
            </div>
          ))}
        </div>
        <Footer />
      </div>

      {/* Confirmation Modal */}
      <ConfirmModal
        message="Tem certeza que deseja excluir esta chave?"
        isOpen={showDeleteModal}
        onConfirm={handleDeleteConfirmed}
        onCancel={handleDeleteCancelled}
      />

      {/* Alert Modal */}
      <AlertDialog
        message={alertMessage}
        isOpen={showAlertModal}
        onClose={closeAlert}
      />

      {/* Key Form Modal */}
      <KeyFormModal
        isOpen={showKeyFormModal}
        onClose={closeKeyFormModal}
        onSave={fetchKeys}
        showAlert={showAlert}
      />

      {/* Borrow Key Modal */}
      <BorrowKeyModal
        isOpen={showBorrowModal}
        onClose={closeBorrowModal}
        onConfirmBorrow={handleConfirmBorrow}
        showAlert={showAlert}
      />

      {/* Return Confirmation Modal */}
      <ConfirmModal
        message={`Tem certeza de que ${borrowerToConfirmReturn || 'o retirante'} fez a devolução da chave?`}
        isOpen={showReturnModal}
        onConfirm={handleReturnConfirmed}
        onCancel={handleReturnCancelled}
      />
    </>
  );
}

export default App;