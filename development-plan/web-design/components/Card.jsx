// components/Card.jsx

// ── Card ───────────────────────────────────────
const Card = ({
  children, title = null, extra = null,
  padding = '20px', style: xStyle, bodyStyle,
  bordered = true, loading = false,
}) => (
  <div className="bp-card-el" style={{
    background: 'var(--bg-surface)',
    borderRadius: 'var(--radius-xl)',
    border: bordered ? '1px solid var(--border)' : 'none',
    boxShadow: 'var(--shadow-sm)',
    overflow: 'hidden',
    ...xStyle,
  }}>
    {(title || extra) && (
      <div style={{
        display: 'flex', alignItems: 'center', justifyContent: 'space-between',
        padding: '14px 20px', borderBottom: '1px solid var(--border)',
      }}>
        <span style={{ fontWeight: '600', fontSize: '14px', color: 'var(--text-primary)' }}>{title}</span>
        {extra && <span>{extra}</span>}
      </div>
    )}
    <div style={{ padding, ...bodyStyle }}>
      {loading ? (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {[80, 60, 90, 50].map((w, i) => (
            <div key={i} className="bp-skeleton" style={{ height: '14px', width: `${w}%` }} />
          ))}
        </div>
      ) : children}
    </div>
  </div>
);

// ── StatCard — dashboard KPI block ────────────
const StatCard = ({ title, value, unit = '', trend = null, trendLabel = '', icon, iconBg, color = 'var(--primary)' }) => (
  <div className="bp-card-el" style={{
    background: 'var(--bg-surface)',
    borderRadius: 'var(--radius-xl)',
    border: '1px solid var(--border)',
    boxShadow: 'var(--shadow-sm)',
    padding: '20px',
    display: 'flex', flexDirection: 'column', gap: '12px',
    transition: 'box-shadow 0.2s, transform 0.2s',
    cursor: 'default',
  }}
    onMouseEnter={e => { e.currentTarget.style.boxShadow = 'var(--shadow-md)'; e.currentTarget.style.borderLeftColor = color; e.currentTarget.style.borderLeftWidth = '3px'; }}
    onMouseLeave={e => { e.currentTarget.style.boxShadow = 'var(--shadow-sm)'; e.currentTarget.style.borderLeftColor = ''; e.currentTarget.style.borderLeftWidth = ''; }}
  >
    <div style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between' }}>
      <span style={{ fontSize: '13px', color: 'var(--text-secondary)', fontWeight: '500' }}>{title}</span>
      {icon && (
        <div style={{
          width: '40px', height: '40px', borderRadius: 'var(--radius-lg)',
          background: iconBg || 'var(--primary-light)',
          display: 'flex', alignItems: 'center', justifyContent: 'center',
          fontSize: '18px', color: color, flexShrink: 0,
        }}>
          <i className={`bi ${icon}`} />
        </div>
      )}
    </div>
    <div style={{ display: 'flex', alignItems: 'baseline', gap: '6px' }}>
      <span style={{ fontSize: '28px', fontWeight: '700', color: 'var(--text-primary)', lineHeight: 1 }}>{value}</span>
      {unit && <span style={{ fontSize: '13px', color: 'var(--text-secondary)' }}>{unit}</span>}
    </div>
    {trend !== null && (
      <div style={{ display: 'flex', alignItems: 'center', gap: '5px', fontSize: '12px' }}>
        <span style={{ color: trend >= 0 ? 'var(--success)' : 'var(--danger)', display: 'flex', alignItems: 'center', gap: '2px', fontWeight: '600' }}>
          <i className={`bi bi-arrow-${trend >= 0 ? 'up' : 'down'}-short`} />
          {Math.abs(trend)}%
        </span>
        <span style={{ color: 'var(--text-tertiary)' }}>{trendLabel}</span>
      </div>
    )}
  </div>
);

// ── SectionCard — titled content section ───────
const SectionCard = ({ title, subtitle, children, action, style: xStyle }) => (
  <Card
    title={
      <div>
        <div style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-primary)' }}>{title}</div>
        {subtitle && <div style={{ fontSize: '12px', fontWeight: '400', color: 'var(--text-tertiary)', marginTop: '2px' }}>{subtitle}</div>}
      </div>
    }
    extra={action}
    style={xStyle}
  >
    {children}
  </Card>
);

Object.assign(window, { Card, StatCard, SectionCard });
