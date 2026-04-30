// components/Alert.jsx

// ── Alert Bar ──────────────────────────────────
const Alert = ({ type = 'info', title, message, closable = false, onClose, style: xStyle, action }) => {
  const [closed, setClosed] = React.useState(false);
  if (closed) return null;

  const typeMap = {
    success: { bg: 'var(--success-light)', border: 'var(--success)', color: 'var(--success-text)', icon: 'bi-check-circle-fill' },
    warning: { bg: 'var(--warning-light)', border: 'var(--warning)', color: 'var(--warning-text)', icon: 'bi-exclamation-triangle-fill' },
    error:   { bg: 'var(--danger-light)',  border: 'var(--danger)',  color: 'var(--danger-text)',  icon: 'bi-x-circle-fill' },
    info:    { bg: 'var(--info-light)',    border: 'var(--info)',    color: 'var(--info-text)',    icon: 'bi-info-circle-fill' },
  };
  const t = typeMap[type] || typeMap.info;

  return (
    <div style={{
      display: 'flex', alignItems: 'flex-start', gap: '10px',
      padding: '12px 14px',
      background: t.bg,
      border: `1px solid ${t.border}`,
      borderRadius: 'var(--radius-lg)',
      color: t.color,
      animation: 'bpSlideUp 0.2s ease',
      ...xStyle,
    }}>
      <i className={`bi ${t.icon}`} style={{ fontSize: '15px', marginTop: '1px', flexShrink: 0 }} />
      <div style={{ flex: 1, minWidth: 0 }}>
        {title && <div style={{ fontWeight: '600', fontSize: '13px', marginBottom: message ? '3px' : 0 }}>{title}</div>}
        {message && <div style={{ fontSize: '13px', opacity: 0.9 }}>{message}</div>}
      </div>
      {action && <div>{action}</div>}
      {closable && (
        <button onClick={() => { setClosed(true); onClose && onClose(); }}
          style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'inherit', opacity: 0.7, padding: '0', display: 'flex', alignItems: 'center' }}>
          <i className="bi bi-x" style={{ fontSize: '16px' }} />
        </button>
      )}
    </div>
  );
};

// ── Toast Notification System ──────────────────
const ToastContainer = () => {
  const [toasts, setToasts] = React.useState([]);

  React.useEffect(() => {
    const handler = e => {
      const { type, message, duration = 3000 } = e.detail;
      const id = Date.now() + Math.random();
      setToasts(prev => [...prev, { id, type, message, leaving: false }]);
      setTimeout(() => {
        setToasts(prev => prev.map(t => t.id === id ? { ...t, leaving: true } : t));
        setTimeout(() => setToasts(prev => prev.filter(t => t.id !== id)), 350);
      }, duration);
    };
    window.addEventListener('bp-toast', handler);
    return () => window.removeEventListener('bp-toast', handler);
  }, []);

  const typeMap = {
    success: { icon: 'bi-check-circle-fill', color: 'var(--success)', bg: '#f0fdf4' },
    error:   { icon: 'bi-x-circle-fill',     color: 'var(--danger)',  bg: '#fff5f5' },
    warning: { icon: 'bi-exclamation-triangle-fill', color: 'var(--warning)', bg: '#fffbeb' },
    info:    { icon: 'bi-info-circle-fill',   color: 'var(--info)',   bg: '#eff6ff' },
  };

  if (toasts.length === 0) return null;

  return ReactDOM.createPortal(
    <div id="bp-toast-root">
      {toasts.map(t => {
        const s = typeMap[t.type] || typeMap.info;
        return (
          <div key={t.id} style={{
            display: 'flex', alignItems: 'center', gap: '10px',
            padding: '12px 16px',
            background: 'var(--bg-surface)',
            border: '1px solid var(--border)',
            borderLeft: `4px solid ${s.color}`,
            borderRadius: 'var(--radius-lg)',
            boxShadow: 'var(--shadow-lg)',
            fontSize: '14px', color: 'var(--text-primary)',
            pointerEvents: 'auto', minWidth: '280px', maxWidth: '360px',
            animation: t.leaving ? 'bpToastOut 0.3s ease forwards' : 'bpToastIn 0.3s ease',
          }}>
            <i className={`bi ${s.icon}`} style={{ color: s.color, fontSize: '16px', flexShrink: 0 }} />
            <span style={{ flex: 1, lineHeight: 1.5 }}>{t.message}</span>
          </div>
        );
      })}
    </div>,
    document.body
  );
};

// Global toast API
window.toast = {
  success: (msg, dur) => window.dispatchEvent(new CustomEvent('bp-toast', { detail: { type: 'success', message: msg, duration: dur } })),
  error:   (msg, dur) => window.dispatchEvent(new CustomEvent('bp-toast', { detail: { type: 'error',   message: msg, duration: dur } })),
  warning: (msg, dur) => window.dispatchEvent(new CustomEvent('bp-toast', { detail: { type: 'warning', message: msg, duration: dur } })),
  info:    (msg, dur) => window.dispatchEvent(new CustomEvent('bp-toast', { detail: { type: 'info',    message: msg, duration: dur } })),
};

// ── Empty State ────────────────────────────────
const Empty = ({ description = '暂无数据', icon = 'bi-inbox', style: xStyle }) => (
  <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', padding: '48px 20px', color: 'var(--text-tertiary)', ...xStyle }}>
    <i className={`bi ${icon}`} style={{ fontSize: '40px', marginBottom: '12px', opacity: 0.5 }} />
    <span style={{ fontSize: '14px' }}>{description}</span>
  </div>
);

Object.assign(window, { Alert, ToastContainer, Empty });
