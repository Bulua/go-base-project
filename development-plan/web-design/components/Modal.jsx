// components/Modal.jsx
const Modal = ({
  open = false, onClose, title, children,
  footer = null, width = 520, size,
  closable = true, maskClosable = true,
  style: xStyle,
}) => {
  const widthMap = { sm: 400, md: 520, lg: 720, xl: 960, full: '90vw' };
  const w = size ? (widthMap[size] || widthMap.md) : width;

  React.useEffect(() => {
    if (open) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => { document.body.style.overflow = ''; };
  }, [open]);

  if (!open) return null;

  return ReactDOM.createPortal(
    <div
      style={{
        position: 'fixed', inset: 0, zIndex: 5000,
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        background: 'var(--bg-mask)',
        animation: 'bpOverlayIn 0.2s ease',
        padding: '20px',
      }}
      onClick={e => { if (maskClosable && e.target === e.currentTarget) onClose && onClose(); }}
    >
      <div className="bp-modal-inner" style={{
        background: 'var(--bg-surface)',
        borderRadius: 'var(--radius-xl)',
        boxShadow: 'var(--shadow-xl)',
        width: w, maxWidth: '100%',
        maxHeight: '90vh', display: 'flex', flexDirection: 'column',
        animation: 'bpModalIn 0.22s ease',
        border: '1px solid var(--border)',
        ...xStyle,
      }}>
        {/* Header */}
        <div style={{
          display: 'flex', alignItems: 'center', justifyContent: 'space-between',
          padding: '16px 20px', borderBottom: '1px solid var(--border)',
          flexShrink: 0,
        }}>
          <span style={{ fontWeight: '600', fontSize: '15px', color: 'var(--text-primary)' }}>{title}</span>
          {closable && (
            <button
              onClick={onClose}
              style={{
                width: '28px', height: '28px', display: 'flex', alignItems: 'center', justifyContent: 'center',
                background: 'transparent', border: 'none', cursor: 'pointer',
                borderRadius: 'var(--radius-md)', color: 'var(--text-tertiary)',
                transition: 'background 0.15s, color 0.15s',
              }}
              onMouseEnter={e => { e.currentTarget.style.background = 'var(--bg-page)'; e.currentTarget.style.color = 'var(--text-primary)'; }}
              onMouseLeave={e => { e.currentTarget.style.background = ''; e.currentTarget.style.color = ''; }}
            >
              <i className="bi bi-x-lg" style={{ fontSize: '14px' }} />
            </button>
          )}
        </div>
        {/* Body */}
        <div style={{ padding: '20px', overflowY: 'auto', flex: 1 }}>
          {children}
        </div>
        {/* Footer */}
        {footer !== null && (
          <div style={{
            display: 'flex', alignItems: 'center', justifyContent: 'flex-end', gap: '8px',
            padding: '14px 20px', borderTop: '1px solid var(--border)',
            flexShrink: 0,
          }}>
            {footer}
          </div>
        )}
      </div>
    </div>,
    document.body
  );
};

// ── ConfirmModal ───────────────────────────────
const ConfirmModal = ({
  open, onClose, onConfirm, title = '确认操作',
  message, type = 'warning', loading = false,
  confirmText = '确认', cancelText = '取消',
}) => {
  const typeMap = {
    warning: { icon: 'bi-exclamation-triangle-fill', color: 'var(--warning)', bg: 'var(--warning-light)' },
    danger:  { icon: 'bi-x-circle-fill',             color: 'var(--danger)',  bg: 'var(--danger-light)' },
    info:    { icon: 'bi-info-circle-fill',           color: 'var(--info)',   bg: 'var(--info-light)' },
  };
  const t = typeMap[type] || typeMap.warning;

  return (
    <Modal open={open} onClose={onClose} title={null} width={420} footer={
      <>
        <Button variant="secondary" onClick={onClose} disabled={loading}>{cancelText}</Button>
        <Button
          variant={type === 'danger' ? 'danger' : 'primary'}
          onClick={onConfirm} loading={loading}
        >{confirmText}</Button>
      </>
    }>
      <div style={{ display: 'flex', gap: '14px', alignItems: 'flex-start' }}>
        <div style={{
          width: '40px', height: '40px', borderRadius: 'var(--radius-full)',
          background: t.bg, display: 'flex', alignItems: 'center', justifyContent: 'center',
          flexShrink: 0,
        }}>
          <i className={`bi ${t.icon}`} style={{ fontSize: '18px', color: t.color }} />
        </div>
        <div>
          <div style={{ fontWeight: '600', fontSize: '15px', color: 'var(--text-primary)', marginBottom: '6px' }}>{title}</div>
          <div style={{ fontSize: '14px', color: 'var(--text-secondary)', lineHeight: 1.6 }}>{message}</div>
        </div>
      </div>
    </Modal>
  );
};

Object.assign(window, { Modal, ConfirmModal });
