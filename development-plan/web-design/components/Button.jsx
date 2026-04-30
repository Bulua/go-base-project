// components/Button.jsx
const Button = ({
  children, variant = 'primary', size = 'md',
  loading = false, disabled = false,
  icon = null, iconRight = null,
  block = false, onClick, type = 'button',
  style: xStyle, className = '',
}) => {
  const [hov, setHov] = React.useState(false);

  const sizes = {
    xs: { h: '26px', px: '10px', fs: '12px', gap: '4px', iconSz: 12 },
    sm: { h: '30px', px: '12px', fs: '13px', gap: '5px', iconSz: 13 },
    md: { h: '34px', px: '16px', fs: '14px', gap: '6px', iconSz: 14 },
    lg: { h: '40px', px: '20px', fs: '15px', gap: '7px', iconSz: 15 },
  };
  const sz = sizes[size] || sizes.md;

  const varMap = {
    primary: {
      bg: 'var(--primary)', color: '#fff',
      border: '1px solid var(--primary)',
      hoverFilter: 'brightness(0.90)',
    },
    secondary: {
      bg: 'var(--bg-surface)', color: 'var(--text-primary)',
      border: '1px solid var(--border)',
      hoverBg: 'var(--bg-page)', hoverBorder: '1px solid var(--border-strong)',
    },
    danger: {
      bg: 'var(--danger)', color: '#fff',
      border: '1px solid var(--danger)',
      hoverFilter: 'brightness(0.90)',
    },
    ghost: {
      bg: 'transparent', color: 'var(--primary)',
      border: '1px solid var(--primary)',
      hoverBg: 'var(--primary-light)',
    },
    text: {
      bg: 'transparent', color: 'var(--text-secondary)',
      border: '1px solid transparent',
      hoverBg: 'var(--bg-page)', hoverColor: 'var(--text-primary)',
    },
    link: {
      bg: 'transparent', color: 'var(--primary)',
      border: 'none', px: '0', h: 'auto',
      hoverDecoration: 'underline',
    },
  };

  const v = varMap[variant] || varMap.primary;
  const isDisabled = disabled || loading;

  const s = {
    display: 'inline-flex', alignItems: 'center', justifyContent: 'center',
    gap: sz.gap,
    height: variant === 'link' ? 'auto' : sz.h,
    padding: variant === 'link' ? '0' : `0 ${sz.px}`,
    fontSize: sz.fs, fontWeight: '500', fontFamily: 'var(--font-sans)',
    borderRadius: 'var(--radius-md)',
    cursor: isDisabled ? 'not-allowed' : 'pointer',
    opacity: isDisabled ? 0.5 : 1,
    transition: 'all 0.15s ease',
    whiteSpace: 'nowrap', outline: 'none',
    width: block ? '100%' : undefined,
    userSelect: 'none',
    background: (hov && v.hoverBg && !isDisabled) ? v.hoverBg : v.bg,
    color: (hov && v.hoverColor && !isDisabled) ? v.hoverColor : v.color,
    border: (hov && v.hoverBorder && !isDisabled) ? v.hoverBorder : v.border,
    filter: (hov && v.hoverFilter && !isDisabled) ? v.hoverFilter : undefined,
    textDecoration: (hov && v.hoverDecoration && !isDisabled) ? v.hoverDecoration : undefined,
    ...xStyle,
  };

  return (
    <button
      type={type} onClick={!isDisabled ? onClick : undefined}
      disabled={isDisabled} style={s} className={className}
      onMouseEnter={() => setHov(true)} onMouseLeave={() => setHov(false)}
    >
      {loading && (
        <span style={{
          width: sz.iconSz, height: sz.iconSz, flexShrink: 0,
          border: '2px solid currentColor', borderTopColor: 'transparent',
          borderRadius: '50%', display: 'inline-block',
          animation: 'bpSpin 0.6s linear infinite',
        }} />
      )}
      {!loading && icon && <span style={{ display: 'flex', fontSize: sz.iconSz, lineHeight: 1 }}>{icon}</span>}
      {children}
      {!loading && iconRight && <span style={{ display: 'flex', fontSize: sz.iconSz, lineHeight: 1 }}>{iconRight}</span>}
    </button>
  );
};

// ButtonGroup — wraps multiple buttons
const ButtonGroup = ({ children, style }) => (
  <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap', ...style }}>
    {children}
  </div>
);

Object.assign(window, { Button, ButtonGroup });
