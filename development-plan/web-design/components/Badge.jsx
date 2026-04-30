// components/Badge.jsx

// ── Tag — colored label chip ───────────────────
const Tag = ({ children, color = 'default', dot = false, style: xStyle, onClose }) => {
  const colorMap = {
    default: { bg: 'var(--bg-page)', color: 'var(--text-secondary)', border: 'var(--border)' },
    primary: { bg: 'var(--primary-light)', color: 'var(--primary)', border: 'var(--primary-dim)' },
    success: { bg: 'var(--success-light)', color: 'var(--success-text)', border: 'var(--success)' },
    warning: { bg: 'var(--warning-light)', color: 'var(--warning-text)', border: 'var(--warning)' },
    danger:  { bg: 'var(--danger-light)',  color: 'var(--danger-text)',  border: 'var(--danger)' },
    info:    { bg: 'var(--info-light)',    color: 'var(--info-text)',    border: 'var(--info)' },
  };
  const c = colorMap[color] || colorMap.default;
  return (
    <span style={{
      display: 'inline-flex', alignItems: 'center', gap: '5px',
      padding: '2px 8px', borderRadius: 'var(--radius-full)',
      fontSize: '12px', fontWeight: '500', lineHeight: '18px',
      background: c.bg, color: c.color,
      border: `1px solid ${c.border}`,
      whiteSpace: 'nowrap',
      ...xStyle,
    }}>
      {dot && (
        <span style={{
          width: '6px', height: '6px', borderRadius: '50%',
          background: 'currentColor', flexShrink: 0,
        }} />
      )}
      {children}
      {onClose && (
        <span
          onClick={onClose}
          style={{ cursor: 'pointer', display: 'flex', alignItems: 'center', marginLeft: '2px', opacity: 0.7 }}
        >
          <i className="bi bi-x" style={{ fontSize: '11px' }} />
        </span>
      )}
    </span>
  );
};

// ── StatusBadge — for table status columns ─────
const StatusBadge = ({ status }) => {
  const map = {
    active:    { label: '启用',   color: 'success' },
    inactive:  { label: '禁用',   color: 'danger'  },
    pending:   { label: '待审核', color: 'warning' },
    success:   { label: '成功',   color: 'success' },
    failed:    { label: '失败',   color: 'danger'  },
    normal:    { label: '正常',   color: 'success' },
    error:     { label: '异常',   color: 'danger'  },
    draft:     { label: '草稿',   color: 'default' },
  };
  const s = map[status] || { label: status, color: 'default' };
  return <Tag color={s.color} dot>{s.label}</Tag>;
};

// ── Badge — numeric indicator ──────────────────
const Badge = ({ count = 0, max = 99, dot = false, children, style: xStyle }) => {
  const display = dot ? null : (count > max ? `${max}+` : count);
  const show = dot ? true : count > 0;
  return (
    <div style={{ position: 'relative', display: 'inline-flex', ...xStyle }}>
      {children}
      {show && (
        <span style={{
          position: 'absolute', top: dot ? '2px' : '-6px', right: dot ? '2px' : '-6px',
          minWidth: dot ? '8px' : '18px', height: dot ? '8px' : '18px',
          background: 'var(--danger)', color: '#fff',
          borderRadius: 'var(--radius-full)',
          fontSize: '11px', fontWeight: '600', lineHeight: '18px',
          textAlign: 'center', padding: dot ? 0 : '0 4px',
          border: '2px solid var(--bg-surface)',
          pointerEvents: 'none',
        }}>
          {display}
        </span>
      )}
    </div>
  );
};

// ── MethodTag — for HTTP method labels ─────────
const MethodTag = ({ method }) => {
  const map = {
    GET:    { bg: '#e8f5e9', color: '#2e7d32' },
    POST:   { bg: '#e3f2fd', color: '#1565c0' },
    PUT:    { bg: '#fff8e1', color: '#f57f17' },
    PATCH:  { bg: '#f3e5f5', color: '#7b1fa2' },
    DELETE: { bg: '#fce4ec', color: '#c62828' },
  };
  const s = map[(method||'').toUpperCase()] || { bg: 'var(--bg-page)', color: 'var(--text-secondary)' };
  return (
    <span style={{
      ...s, padding: '2px 7px', borderRadius: 'var(--radius-sm)',
      fontSize: '11px', fontWeight: '700', fontFamily: 'var(--font-mono)',
      display: 'inline-block',
    }}>
      {(method||'').toUpperCase()}
    </span>
  );
};

Object.assign(window, { Tag, StatusBadge, Badge, MethodTag });
