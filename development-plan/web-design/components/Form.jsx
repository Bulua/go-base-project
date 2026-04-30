// components/Form.jsx
const FormContext = React.createContext({});

// ── Form ───────────────────────────────────────
const Form = ({ children, onSubmit, layout = 'vertical', style: xStyle }) => {
  const handleSubmit = e => { e.preventDefault(); onSubmit && onSubmit(e); };
  return (
    <FormContext.Provider value={{ layout }}>
      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: layout === 'vertical' ? '16px' : '12px', ...xStyle }}>
        {children}
      </form>
    </FormContext.Provider>
  );
};

// ── FormItem ───────────────────────────────────
const FormItem = ({ label, children, required = false, error = '', help = '', colon = true, style: xStyle }) => {
  const { layout } = React.useContext(FormContext);
  const isHoriz = layout === 'horizontal';
  return (
    <div style={{
      display: isHoriz ? 'grid' : 'flex',
      gridTemplateColumns: isHoriz ? '120px 1fr' : undefined,
      flexDirection: isHoriz ? undefined : 'column',
      gap: isHoriz ? '0 12px' : '5px',
      alignItems: isHoriz ? 'flex-start' : undefined,
      ...xStyle,
    }}>
      {label && (
        <label style={{
          fontSize: '13px', fontWeight: '500', color: 'var(--text-primary)',
          paddingTop: isHoriz ? '7px' : undefined,
          textAlign: isHoriz ? 'right' : 'left',
          display: 'flex', alignItems: isHoriz ? 'flex-start' : 'center',
          justifyContent: isHoriz ? 'flex-end' : undefined,
          gap: '3px', flexShrink: 0,
        }}>
          {required && <span style={{ color: 'var(--danger)', fontSize: '14px', lineHeight: 1 }}>*</span>}
          {label}{colon ? '' : ''}
        </label>
      )}
      <div style={{ flex: 1, minWidth: 0 }}>
        {children}
        {error && (
          <div style={{ marginTop: '4px', fontSize: '12px', color: 'var(--danger)', display: 'flex', alignItems: 'center', gap: '4px' }}>
            <i className="bi bi-exclamation-circle" />
            {error}
          </div>
        )}
        {!error && help && (
          <div style={{ marginTop: '4px', fontSize: '12px', color: 'var(--text-tertiary)' }}>
            {help}
          </div>
        )}
      </div>
    </div>
  );
};

// ── FormRow — horizontal group of FormItems ────
const FormRow = ({ children, cols = 2, style: xStyle }) => (
  <div style={{
    display: 'grid',
    gridTemplateColumns: `repeat(${cols}, 1fr)`,
    gap: '16px',
    ...xStyle,
  }}>
    {children}
  </div>
);

// ── FormDivider ────────────────────────────────
const FormDivider = ({ label }) => (
  <div style={{ display: 'flex', alignItems: 'center', gap: '10px', margin: '4px 0' }}>
    {label && <span style={{ fontSize: '12px', color: 'var(--text-tertiary)', whiteSpace: 'nowrap' }}>{label}</span>}
    <div style={{ flex: 1, height: '1px', background: 'var(--border)' }} />
  </div>
);

// ── FormActions — bottom buttons ───────────────
const FormActions = ({ children, align = 'right', style: xStyle }) => (
  <div style={{
    display: 'flex', gap: '8px',
    justifyContent: align === 'right' ? 'flex-end' : align === 'center' ? 'center' : 'flex-start',
    paddingTop: '8px',
    borderTop: '1px solid var(--border)',
    marginTop: '4px',
    ...xStyle,
  }}>
    {children}
  </div>
);

Object.assign(window, { Form, FormItem, FormRow, FormDivider, FormActions });
