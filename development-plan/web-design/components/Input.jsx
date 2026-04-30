// components/Input.jsx
const inputBase = (focused, error, disabled) => ({
  width: '100%', height: '34px', padding: '0 12px',
  fontSize: '14px', fontFamily: 'var(--font-sans)',
  background: disabled ? 'var(--bg-page)' : 'var(--bg-surface)',
  color: disabled ? 'var(--text-disabled)' : 'var(--text-primary)',
  border: `1px solid ${error ? 'var(--danger)' : focused ? 'var(--primary)' : 'var(--border)'}`,
  borderRadius: 'var(--radius-md)',
  outline: 'none',
  transition: 'border-color 0.15s, box-shadow 0.15s',
  boxShadow: focused ? (error
    ? '0 0 0 3px oklch(0.97 0.04 25)'
    : '0 0 0 3px var(--primary-light)') : 'none',
  cursor: disabled ? 'not-allowed' : 'text',
});

// ── Input ──────────────────────────────────────
const Input = React.forwardRef(({
  value, onChange, placeholder = '', type = 'text',
  prefix = null, suffix = null, error = false,
  disabled = false, readOnly = false, size = 'md',
  style: xStyle, onEnter, onBlur, onFocus,
  clearable = false, ...rest
}, ref) => {
  const [focused, setFocused] = React.useState(false);
  const [showPwd, setShowPwd] = React.useState(false);

  const heights = { sm: '28px', md: '34px', lg: '40px' };
  const fontSizes = { sm: '12px', md: '14px', lg: '15px' };
  const pads = { sm: '0 8px', md: '0 12px', lg: '0 14px' };

  const wrapStyle = {
    position: 'relative', display: 'flex', alignItems: 'center',
    background: disabled ? 'var(--bg-page)' : 'var(--bg-surface)',
    border: `1px solid ${error ? 'var(--danger)' : focused ? 'var(--primary)' : 'var(--border)'}`,
    borderRadius: 'var(--radius-md)',
    boxShadow: focused ? (error ? '0 0 0 3px oklch(0.97 0.04 25)' : '0 0 0 3px var(--primary-light)') : 'none',
    transition: 'border-color 0.15s, box-shadow 0.15s',
    height: heights[size] || heights.md,
    cursor: disabled ? 'not-allowed' : undefined,
    ...xStyle,
  };

  const inputStyle = {
    flex: 1, height: '100%',
    padding: pads[size],
    fontSize: fontSizes[size],
    fontFamily: 'var(--font-sans)',
    background: 'transparent',
    color: disabled ? 'var(--text-disabled)' : 'var(--text-primary)',
    border: 'none', outline: 'none',
    cursor: disabled ? 'not-allowed' : 'text',
    minWidth: 0,
  };

  const affixStyle = {
    display: 'flex', alignItems: 'center', padding: '0 10px',
    color: 'var(--text-tertiary)', fontSize: '13px', flexShrink: 0,
  };

  const inputType = type === 'password' ? (showPwd ? 'text' : 'password') : type;

  return (
    <div style={wrapStyle}>
      {prefix && <span style={{ ...affixStyle, paddingRight: 0 }}>{prefix}</span>}
      <input
        ref={ref} type={inputType} value={value} placeholder={placeholder}
        disabled={disabled} readOnly={readOnly}
        style={inputStyle}
        onChange={onChange}
        onFocus={e => { setFocused(true); onFocus && onFocus(e); }}
        onBlur={e => { setFocused(false); onBlur && onBlur(e); }}
        onKeyDown={e => e.key === 'Enter' && onEnter && onEnter(e)}
        {...rest}
      />
      {clearable && value && !disabled && (
        <span
          style={{ ...affixStyle, cursor: 'pointer', paddingLeft: 0 }}
          onClick={() => onChange && onChange({ target: { value: '' } })}
        >
          <i className="bi bi-x-circle-fill" style={{ fontSize: '13px', color: 'var(--text-tertiary)' }} />
        </span>
      )}
      {type === 'password' && (
        <span
          style={{ ...affixStyle, cursor: 'pointer', paddingLeft: 0 }}
          onClick={() => setShowPwd(p => !p)}
        >
          <i className={`bi bi-${showPwd ? 'eye-slash' : 'eye'}`} />
        </span>
      )}
      {suffix && <span style={affixStyle}>{suffix}</span>}
    </div>
  );
});

// ── Textarea ───────────────────────────────────
const Textarea = ({ value, onChange, placeholder = '', rows = 4, error = false, disabled = false, style: xStyle }) => {
  const [focused, setFocused] = React.useState(false);
  return (
    <textarea
      value={value} onChange={onChange} placeholder={placeholder}
      rows={rows} disabled={disabled}
      style={{
        width: '100%', padding: '8px 12px', fontSize: '14px',
        fontFamily: 'var(--font-sans)', lineHeight: 1.6,
        background: disabled ? 'var(--bg-page)' : 'var(--bg-surface)',
        color: disabled ? 'var(--text-disabled)' : 'var(--text-primary)',
        border: `1px solid ${error ? 'var(--danger)' : focused ? 'var(--primary)' : 'var(--border)'}`,
        borderRadius: 'var(--radius-md)', outline: 'none', resize: 'vertical',
        boxShadow: focused ? '0 0 0 3px var(--primary-light)' : 'none',
        transition: 'border-color 0.15s, box-shadow 0.15s',
        ...xStyle,
      }}
      onFocus={() => setFocused(true)}
      onBlur={() => setFocused(false)}
    />
  );
};

// ── Select ─────────────────────────────────────
const Select = ({ value, onChange, options = [], placeholder = '请选择', disabled = false, error = false, style: xStyle, size = 'md' }) => {
  const [focused, setFocused] = React.useState(false);
  const heights = { sm: '28px', md: '34px', lg: '40px' };
  const fontSizes = { sm: '12px', md: '14px', lg: '15px' };
  return (
    <div style={{ position: 'relative', ...xStyle }}>
      <select
        value={value} onChange={onChange} disabled={disabled}
        style={{
          width: '100%', height: heights[size], padding: '0 32px 0 12px',
          fontSize: fontSizes[size], fontFamily: 'var(--font-sans)',
          background: disabled ? 'var(--bg-page)' : 'var(--bg-surface)',
          color: value ? 'var(--text-primary)' : 'var(--text-tertiary)',
          border: `1px solid ${error ? 'var(--danger)' : focused ? 'var(--primary)' : 'var(--border)'}`,
          borderRadius: 'var(--radius-md)', outline: 'none', cursor: disabled ? 'not-allowed' : 'pointer',
          appearance: 'none', transition: 'border-color 0.15s, box-shadow 0.15s',
          boxShadow: focused ? '0 0 0 3px var(--primary-light)' : 'none',
        }}
        onFocus={() => setFocused(true)}
        onBlur={() => setFocused(false)}
      >
        {placeholder && <option value="">{placeholder}</option>}
        {options.map(o => (
          <option key={typeof o === 'object' ? o.value : o} value={typeof o === 'object' ? o.value : o}>
            {typeof o === 'object' ? o.label : o}
          </option>
        ))}
      </select>
      <i className="bi bi-chevron-down" style={{
        position: 'absolute', right: '10px', top: '50%', transform: 'translateY(-50%)',
        pointerEvents: 'none', fontSize: '12px', color: 'var(--text-secondary)',
      }} />
    </div>
  );
};

// ── Checkbox ───────────────────────────────────
const Checkbox = ({ checked, onChange, children, disabled = false, indeterminate = false }) => {
  const ref = React.useRef(null);
  React.useEffect(() => {
    if (ref.current) ref.current.indeterminate = indeterminate;
  }, [indeterminate]);
  return (
    <label style={{ display: 'inline-flex', alignItems: 'center', gap: '7px', cursor: disabled ? 'not-allowed' : 'pointer', opacity: disabled ? 0.5 : 1, userSelect: 'none', fontSize: '14px', color: 'var(--text-primary)' }}>
      <input
        ref={ref} type="checkbox" checked={checked} onChange={onChange}
        disabled={disabled}
        style={{ width: '15px', height: '15px', accentColor: 'var(--primary)', cursor: disabled ? 'not-allowed' : 'pointer' }}
      />
      {children}
    </label>
  );
};

// ── Radio ──────────────────────────────────────
const Radio = ({ value, checked, onChange, children, disabled = false }) => (
  <label style={{ display: 'inline-flex', alignItems: 'center', gap: '7px', cursor: disabled ? 'not-allowed' : 'pointer', opacity: disabled ? 0.5 : 1, userSelect: 'none', fontSize: '14px', color: 'var(--text-primary)' }}>
    <input type="radio" value={value} checked={checked} onChange={onChange} disabled={disabled}
      style={{ width: '15px', height: '15px', accentColor: 'var(--primary)', cursor: 'pointer' }} />
    {children}
  </label>
);

// ── RadioGroup ─────────────────────────────────
const RadioGroup = ({ value, onChange, options = [], disabled = false }) => (
  <div style={{ display: 'flex', gap: '20px', flexWrap: 'wrap' }}>
    {options.map(o => (
      <Radio key={o.value} value={o.value} checked={value === o.value}
        onChange={() => onChange(o.value)} disabled={disabled || o.disabled}>
        {o.label}
      </Radio>
    ))}
  </div>
);

// ── SearchInput ────────────────────────────────
const SearchInput = ({ value, onChange, placeholder = '搜索...', onSearch, style: xStyle }) => (
  <Input
    value={value} onChange={onChange} placeholder={placeholder}
    prefix={<i className="bi bi-search" style={{ color: 'var(--text-tertiary)', fontSize: '13px' }} />}
    clearable onEnter={onSearch} style={xStyle}
  />
);

Object.assign(window, { Input, Textarea, Select, Checkbox, Radio, RadioGroup, SearchInput });
