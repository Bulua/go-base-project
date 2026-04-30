// components/Dropdown.jsx
const Dropdown = ({ trigger, items = [], placement = 'bottom-right', disabled = false }) => {
  const [open, setOpen] = React.useState(false);
  const ref = React.useRef(null);

  React.useEffect(() => {
    if (!open) return;
    const handler = e => {
      if (ref.current && !ref.current.contains(e.target)) setOpen(false);
    };
    document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, [open]);

  const placementStyle = {
    'bottom-right': { top: 'calc(100% + 6px)', right: 0 },
    'bottom-left':  { top: 'calc(100% + 6px)', left: 0 },
    'top-right':    { bottom: 'calc(100% + 6px)', right: 0 },
    'top-left':     { bottom: 'calc(100% + 6px)', left: 0 },
  }[placement] || { top: 'calc(100% + 6px)', right: 0 };

  return (
    <div ref={ref} style={{ position: 'relative', display: 'inline-flex' }}>
      <div onClick={() => !disabled && setOpen(p => !p)} style={{ cursor: disabled ? 'not-allowed' : 'pointer', opacity: disabled ? 0.5 : 1 }}>
        {trigger}
      </div>
      {open && (
        <div className="bp-dropdown-menu" style={{
          position: 'absolute', zIndex: 2000,
          ...placementStyle,
          background: 'var(--bg-surface)',
          border: '1px solid var(--border)',
          borderRadius: 'var(--radius-lg)',
          boxShadow: 'var(--shadow-lg)',
          minWidth: '160px',
          padding: '4px',
          animation: 'bpSlideUp 0.15s ease',
          overflow: 'hidden',
        }}>
          {items.map((item, i) => {
            if (item.type === 'divider') {
              return <div key={i} style={{ height: '1px', background: 'var(--border)', margin: '4px 0' }} />;
            }
            if (item.type === 'header') {
              return (
                <div key={i} style={{ padding: '6px 12px 4px', fontSize: '11px', fontWeight: '600', color: 'var(--text-tertiary)', textTransform: 'uppercase', letterSpacing: '0.05em' }}>
                  {item.label}
                </div>
              );
            }
            return (
              <DropdownItem key={i} item={item} onClose={() => setOpen(false)} />
            );
          })}
        </div>
      )}
    </div>
  );
};

const DropdownItem = ({ item, onClose }) => {
  const [hov, setHov] = React.useState(false);
  return (
    <div
      onClick={() => { if (!item.disabled) { item.onClick && item.onClick(); onClose(); } }}
      onMouseEnter={() => setHov(true)}
      onMouseLeave={() => setHov(false)}
      style={{
        display: 'flex', alignItems: 'center', gap: '8px',
        padding: '7px 12px', borderRadius: 'var(--radius-md)',
        fontSize: '13px', cursor: item.disabled ? 'not-allowed' : 'pointer',
        color: item.danger ? 'var(--danger)' : item.disabled ? 'var(--text-disabled)' : 'var(--text-primary)',
        background: hov && !item.disabled ? (item.danger ? 'var(--danger-light)' : 'var(--bg-page)') : 'transparent',
        transition: 'background 0.12s',
        userSelect: 'none',
        opacity: item.disabled ? 0.5 : 1,
      }}
    >
      {item.icon && <i className={`bi ${item.icon}`} style={{ fontSize: '14px', flexShrink: 0 }} />}
      <span style={{ flex: 1 }}>{item.label}</span>
      {item.badge && (
        <span style={{ background: 'var(--danger)', color: '#fff', borderRadius: 'var(--radius-full)', padding: '1px 6px', fontSize: '11px', fontWeight: '600' }}>
          {item.badge}
        </span>
      )}
    </div>
  );
};

// ── ContextMenu (right-click style) ───────────
const MoreActions = ({ items }) => (
  <Dropdown
    placement="bottom-right"
    trigger={
      <div style={{
        width: '30px', height: '30px', display: 'flex', alignItems: 'center', justifyContent: 'center',
        borderRadius: 'var(--radius-md)', color: 'var(--text-tertiary)',
        transition: 'background 0.15s, color 0.15s',
      }}
        onMouseEnter={e => { e.currentTarget.style.background = 'var(--bg-page)'; e.currentTarget.style.color = 'var(--text-primary)'; }}
        onMouseLeave={e => { e.currentTarget.style.background = ''; e.currentTarget.style.color = ''; }}
      >
        <i className="bi bi-three-dots-vertical" />
      </div>
    }
    items={items}
  />
);

Object.assign(window, { Dropdown, MoreActions });
