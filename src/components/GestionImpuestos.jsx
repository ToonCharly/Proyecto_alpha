import React, { useState, useEffect } from 'react';

const GestionImpuestos = ({ empresaId, onClose, onImpuestoActualizado }) => {
  const [impuestos, setImpuestos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [editando, setEditando] = useState(null);
  const [nuevoImpuesto, setNuevoImpuesto] = useState({
    descripcion: '',
    iva: 0,
    tipo_iva: 'Tasa',
    ieps1: 0,
    tipo_ieps1: 'Tasa',
    ieps2: 0,
    tipo_ieps2: 'Tasa',
    ieps3: 0,
    tipo_ieps3: 'Tasa'
  });

  // Cargar impuestos de la empresa
  useEffect(() => {
    if (empresaId) {
      cargarImpuestos();
    }
  }, [empresaId]);

  const cargarImpuestos = async () => {
    try {
      setLoading(true);
      const response = await fetch(`http://localhost:8080/api/impuestos?idempresa=${empresaId}`);
      
      if (!response.ok) {
        throw new Error('Error al cargar impuestos');
      }
      
      const data = await response.json();
      setImpuestos(data.impuestos || []);
    } catch (err) {
      setError(err.message);
      console.error('Error al cargar impuestos:', err);
    } finally {
      setLoading(false);
    }
  };

  const crearImpuesto = async () => {
    try {
      const impuestoData = {
        ...nuevoImpuesto,
        idempresa: empresaId
      };

      const response = await fetch('http://localhost:8080/api/impuestos', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(impuestoData),
      });

      if (!response.ok) {
        throw new Error('Error al crear impuesto');
      }

      await cargarImpuestos();
      setNuevoImpuesto({
        descripcion: '',
        iva: 0,
        tipo_iva: 'Tasa',
        ieps1: 0,
        tipo_ieps1: 'Tasa',
        ieps2: 0,
        tipo_ieps2: 'Tasa',
        ieps3: 0,
        tipo_ieps3: 'Tasa'
      });
      
      if (onImpuestoActualizado) {
        onImpuestoActualizado();
      }
    } catch (err) {
      setError(err.message);
    }
  };

  const actualizarImpuesto = async (impuesto) => {
    try {
      const response = await fetch('http://localhost:8080/api/impuestos', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(impuesto),
      });

      if (!response.ok) {
        throw new Error('Error al actualizar impuesto');
      }

      await cargarImpuestos();
      setEditando(null);
      
      if (onImpuestoActualizado) {
        onImpuestoActualizado();
      }
    } catch (err) {
      setError(err.message);
    }
  };

  const eliminarImpuesto = async (idiva) => {
    if (!window.confirm('¿Está seguro de eliminar este impuesto?')) {
      return;
    }

    try {
      const response = await fetch(`http://localhost:8080/api/impuestos?idiva=${idiva}&idempresa=${empresaId}`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        throw new Error('Error al eliminar impuesto');
      }

      await cargarImpuestos();
      
      if (onImpuestoActualizado) {
        onImpuestoActualizado();
      }
    } catch (err) {
      setError(err.message);
    }
  };

  if (loading) {
    return (
      <div className="gestion-impuestos-overlay">
        <div className="gestion-impuestos-modal">
          <div className="loading">Cargando impuestos...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="gestion-impuestos-overlay">
      <div className="gestion-impuestos-modal">
        <div className="modal-header">
          <h2>Gestión de Impuestos</h2>
          <button onClick={onClose} className="btn-cerrar">×</button>
        </div>

        {error && (
          <div className="error-message">
            {error}
          </div>
        )}

        {/* Formulario para nuevo impuesto */}
        <div className="nuevo-impuesto-form">
          <h3>Agregar Nuevo Impuesto</h3>
          <div className="form-row">
            <input
              type="text"
              placeholder="Descripción del impuesto"
              value={nuevoImpuesto.descripcion}
              onChange={(e) => setNuevoImpuesto({...nuevoImpuesto, descripcion: e.target.value})}
            />
          </div>
          
          <div className="impuestos-grid">
            <div className="impuesto-group">
              <label>IVA</label>
              <div className="impuesto-inputs">
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  max="100"
                  placeholder="% IVA"
                  value={nuevoImpuesto.iva}
                  onChange={(e) => setNuevoImpuesto({...nuevoImpuesto, iva: parseFloat(e.target.value) || 0})}
                />
                <select
                  value={nuevoImpuesto.tipo_iva}
                  onChange={(e) => setNuevoImpuesto({...nuevoImpuesto, tipo_iva: e.target.value})}
                >
                  <option value="Tasa">Tasa</option>
                  <option value="Cuota">Cuota</option>
                  <option value="Exento">Exento</option>
                </select>
              </div>
            </div>

            <div className="impuesto-group">
              <label>IEPS 1</label>
              <div className="impuesto-inputs">
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  max="100"
                  placeholder="% IEPS 1"
                  value={nuevoImpuesto.ieps1}
                  onChange={(e) => setNuevoImpuesto({...nuevoImpuesto, ieps1: parseFloat(e.target.value) || 0})}
                />
                <select
                  value={nuevoImpuesto.tipo_ieps1}
                  onChange={(e) => setNuevoImpuesto({...nuevoImpuesto, tipo_ieps1: e.target.value})}
                >
                  <option value="Tasa">Tasa</option>
                  <option value="Cuota">Cuota</option>
                  <option value="Exento">Exento</option>
                </select>
              </div>
            </div>

            <div className="impuesto-group">
              <label>IEPS 2</label>
              <div className="impuesto-inputs">
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  max="100"
                  placeholder="% IEPS 2"
                  value={nuevoImpuesto.ieps2}
                  onChange={(e) => setNuevoImpuesto({...nuevoImpuesto, ieps2: parseFloat(e.target.value) || 0})}
                />
                <select
                  value={nuevoImpuesto.tipo_ieps2}
                  onChange={(e) => setNuevoImpuesto({...nuevoImpuesto, tipo_ieps2: e.target.value})}
                >
                  <option value="Tasa">Tasa</option>
                  <option value="Cuota">Cuota</option>
                  <option value="Exento">Exento</option>
                </select>
              </div>
            </div>

            <div className="impuesto-group">
              <label>IEPS 3</label>
              <div className="impuesto-inputs">
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  max="100"
                  placeholder="% IEPS 3"
                  value={nuevoImpuesto.ieps3}
                  onChange={(e) => setNuevoImpuesto({...nuevoImpuesto, ieps3: parseFloat(e.target.value) || 0})}
                />
                <select
                  value={nuevoImpuesto.tipo_ieps3}
                  onChange={(e) => setNuevoImpuesto({...nuevoImpuesto, tipo_ieps3: e.target.value})}
                >
                  <option value="Tasa">Tasa</option>
                  <option value="Cuota">Cuota</option>
                  <option value="Exento">Exento</option>
                </select>
              </div>
            </div>
          </div>

          <button onClick={crearImpuesto} className="btn-crear">
            Crear Impuesto
          </button>
        </div>

        {/* Lista de impuestos existentes */}
        <div className="impuestos-existentes">
          <h3>Impuestos Configurados</h3>
          {impuestos.length === 0 ? (
            <p>No hay impuestos configurados para esta empresa.</p>
          ) : (
            <div className="impuestos-tabla">
              {impuestos.map((impuesto) => (
                <div key={impuesto.idiva} className="impuesto-item">
                  {editando === impuesto.idiva ? (
                    <EditarImpuestoForm 
                      impuesto={impuesto}
                      onGuardar={actualizarImpuesto}
                      onCancelar={() => setEditando(null)}
                    />
                  ) : (
                    <div className="impuesto-info">
                      <div className="impuesto-descripcion">
                        <strong>{impuesto.descripcion}</strong>
                      </div>
                      <div className="impuesto-detalles">
                        <span>IVA: {impuesto.iva}% ({impuesto.tipo_iva})</span>
                        {impuesto.ieps1 > 0 && <span>IEPS1: {impuesto.ieps1}% ({impuesto.tipo_ieps1})</span>}
                        {impuesto.ieps2 > 0 && <span>IEPS2: {impuesto.ieps2}% ({impuesto.tipo_ieps2})</span>}
                        {impuesto.ieps3 > 0 && <span>IEPS3: {impuesto.ieps3}% ({impuesto.tipo_ieps3})</span>}
                      </div>
                      <div className="impuesto-acciones">
                        <button onClick={() => setEditando(impuesto.idiva)} className="btn-editar">
                          Editar
                        </button>
                        <button onClick={() => eliminarImpuesto(impuesto.idiva)} className="btn-eliminar">
                          Eliminar
                        </button>
                      </div>
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

// Componente para editar impuesto
const EditarImpuestoForm = ({ impuesto, onGuardar, onCancelar }) => {
  const [editData, setEditData] = useState({ ...impuesto });

  return (
    <div className="editar-impuesto-form">
      <div className="form-row">
        <input
          type="text"
          value={editData.descripcion}
          onChange={(e) => setEditData({...editData, descripcion: e.target.value})}
        />
      </div>
      
      <div className="impuestos-grid">
        <div className="impuesto-group">
          <label>IVA</label>
          <div className="impuesto-inputs">
            <input
              type="number"
              step="0.01"
              min="0"
              max="100"
              value={editData.iva}
              onChange={(e) => setEditData({...editData, iva: parseFloat(e.target.value) || 0})}
            />
            <select
              value={editData.tipo_iva}
              onChange={(e) => setEditData({...editData, tipo_iva: e.target.value})}
            >
              <option value="Tasa">Tasa</option>
              <option value="Cuota">Cuota</option>
              <option value="Exento">Exento</option>
            </select>
          </div>
        </div>

        <div className="impuesto-group">
          <label>IEPS 1</label>
          <div className="impuesto-inputs">
            <input
              type="number"
              step="0.01"
              min="0"
              max="100"
              value={editData.ieps1}
              onChange={(e) => setEditData({...editData, ieps1: parseFloat(e.target.value) || 0})}
            />
            <select
              value={editData.tipo_ieps1}
              onChange={(e) => setEditData({...editData, tipo_ieps1: e.target.value})}
            >
              <option value="Tasa">Tasa</option>
              <option value="Cuota">Cuota</option>
              <option value="Exento">Exento</option>
            </select>
          </div>
        </div>

        <div className="impuesto-group">
          <label>IEPS 2</label>
          <div className="impuesto-inputs">
            <input
              type="number"
              step="0.01"
              min="0"
              max="100"
              value={editData.ieps2}
              onChange={(e) => setEditData({...editData, ieps2: parseFloat(e.target.value) || 0})}
            />
            <select
              value={editData.tipo_ieps2}
              onChange={(e) => setEditData({...editData, tipo_ieps2: e.target.value})}
            >
              <option value="Tasa">Tasa</option>
              <option value="Cuota">Cuota</option>
              <option value="Exento">Exento</option>
            </select>
          </div>
        </div>

        <div className="impuesto-group">
          <label>IEPS 3</label>
          <div className="impuesto-inputs">
            <input
              type="number"
              step="0.01"
              min="0"
              max="100"
              value={editData.ieps3}
              onChange={(e) => setEditData({...editData, ieps3: parseFloat(e.target.value) || 0})}
            />
            <select
              value={editData.tipo_ieps3}
              onChange={(e) => setEditData({...editData, tipo_ieps3: e.target.value})}
            >
              <option value="Tasa">Tasa</option>
              <option value="Cuota">Cuota</option>
              <option value="Exento">Exento</option>
            </select>
          </div>
        </div>
      </div>

      <div className="form-actions">
        <button onClick={() => onGuardar(editData)} className="btn-guardar">
          Guardar
        </button>
        <button onClick={onCancelar} className="btn-cancelar">
          Cancelar
        </button>
      </div>
    </div>
  );
};

export default GestionImpuestos;
