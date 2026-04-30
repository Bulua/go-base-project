// pages/Login.jsx
const Login = ({ onLogin }) => {
  const [username, setUsername] = React.useState('');
  const [password, setPassword] = React.useState('');
  const [loading, setLoading] = React.useState(false);
  const [errors, setErrors] = React.useState({});
  const [captcha, setCaptcha] = React.useState('');
  const [captchaVal] = React.useState('BP24');

  const validate = () => {
    const e = {};
    if (!username.trim()) e.username = '请输入用户名';
    if (!password) e.password = '请输入密码';
    else if (password.length < 4) e.password = '密码至少 4 位';
    if (!captcha.trim()) e.captcha = '请输入验证码';
    else if (captcha.toUpperCase() !== captchaVal) e.captcha = '验证码不正确';
    return e;
  };

  const handleLogin = async e => {
    e && e.preventDefault();
    const errs = validate();
    if (Object.keys(errs).length) { setErrors(errs); return; }
    setLoading(true);
    await new Promise(r => setTimeout(r, 1000));
    setLoading(false);
    onLogin && onLogin({ username });
  };

  const presets = Object.entries(BaseTokens.PRESET_HUES);

  return (
    <div style={{
      minHeight: '100vh', display: 'flex', background: 'var(--bg-page)',
      fontFamily: 'var(--font-sans)',
    }}>
      {/* Left decorative panel */}
      <div style={{
        flex: 1, background: 'var(--primary)',
        backgroundImage: `
          radial-gradient(circle at 20% 30%, oklch(0.70 0.14 var(--hue)) 0%, transparent 50%),
          radial-gradient(circle at 80% 70%, oklch(0.38 0.20 var(--hue)) 0%, transparent 50%)
        `,
        display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center',
        padding: '60px', color: '#fff', position: 'relative', overflow: 'hidden',
        minHeight: '100vh',
      }}
        className="login-panel-left"
      >
        {/* Decorative circles */}
        <div style={{ position: 'absolute', top: '-80px', left: '-80px', width: '300px', height: '300px', borderRadius: '50%', background: 'rgba(255,255,255,0.06)' }} />
        <div style={{ position: 'absolute', bottom: '-60px', right: '-60px', width: '250px', height: '250px', borderRadius: '50%', background: 'rgba(255,255,255,0.06)' }} />
        <div style={{ position: 'absolute', top: '40%', right: '-40px', width: '160px', height: '160px', borderRadius: '50%', background: 'rgba(255,255,255,0.04)' }} />

        <div style={{ position: 'relative', zIndex: 1, maxWidth: '360px', textAlign: 'center', animation: 'bpFadeIn 0.6s ease' }}>
          <div style={{
            width: '72px', height: '72px', borderRadius: '20px',
            background: 'rgba(255,255,255,0.2)',
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            margin: '0 auto 28px', fontSize: '32px', fontWeight: '800',
            backdropFilter: 'blur(10px)',
            border: '1px solid rgba(255,255,255,0.3)',
          }}>B</div>
          <h1 style={{ fontSize: '32px', fontWeight: '800', marginBottom: '14px', lineHeight: 1.2 }}>
            BaseProject
          </h1>
          <p style={{ fontSize: '15px', opacity: 0.8, lineHeight: 1.8 }}>
            统一基座管理平台<br />
            安全、高效、可扩展的企业级解决方案
          </p>

          <div style={{ display: 'flex', gap: '24px', justifyContent: 'center', marginTop: '40px', flexWrap: 'wrap' }}>
            {[
              { icon: 'bi-shield-check', label: '安全认证' },
              { icon: 'bi-layers', label: '模块化架构' },
              { icon: 'bi-bar-chart', label: '实时监控' },
            ].map(f => (
              <div key={f.label} style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '8px', opacity: 0.85 }}>
                <div style={{ width: '44px', height: '44px', borderRadius: '12px', background: 'rgba(255,255,255,0.15)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '18px' }}>
                  <i className={`bi ${f.icon}`} />
                </div>
                <span style={{ fontSize: '12px', fontWeight: '500' }}>{f.label}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Right login form */}
      <div style={{
        width: '440px', flexShrink: 0, display: 'flex', flexDirection: 'column',
        alignItems: 'center', justifyContent: 'center',
        background: 'var(--bg-surface)', padding: '48px 48px',
        boxShadow: '-4px 0 24px rgba(0,0,0,0.06)',
      }}>
        <div style={{ width: '100%', maxWidth: '340px', animation: 'bpSlideUp 0.4s ease' }}>
          <h2 style={{ fontSize: '24px', fontWeight: '800', color: 'var(--text-primary)', marginBottom: '6px' }}>欢迎回来</h2>
          <p style={{ fontSize: '14px', color: 'var(--text-secondary)', marginBottom: '32px' }}>登录您的账号以继续</p>

          <form onSubmit={handleLogin} style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
            {/* Username */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
              <label style={{ fontSize: '13px', fontWeight: '500', color: 'var(--text-primary)' }}>用户名</label>
              <Input
                value={username} onChange={e => { setUsername(e.target.value); setErrors(p => ({...p, username: ''})); }}
                placeholder="请输入用户名"
                prefix={<i className="bi bi-person" style={{ color: 'var(--text-tertiary)', fontSize: '14px' }} />}
                error={!!errors.username}
                size="lg"
                onEnter={handleLogin}
              />
              {errors.username && <span style={{ fontSize: '12px', color: 'var(--danger)' }}><i className="bi bi-exclamation-circle" /> {errors.username}</span>}
            </div>

            {/* Password */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <label style={{ fontSize: '13px', fontWeight: '500', color: 'var(--text-primary)' }}>密码</label>
                <a href="#" onClick={e => { e.preventDefault(); window.toast && window.toast.info('请联系管理员重置密码'); }} style={{ fontSize: '12px', color: 'var(--primary)' }}>忘记密码？</a>
              </div>
              <Input
                type="password" value={password}
                onChange={e => { setPassword(e.target.value); setErrors(p => ({...p, password: ''})); }}
                placeholder="请输入密码"
                prefix={<i className="bi bi-lock" style={{ color: 'var(--text-tertiary)', fontSize: '14px' }} />}
                error={!!errors.password} size="lg" onEnter={handleLogin}
              />
              {errors.password && <span style={{ fontSize: '12px', color: 'var(--danger)' }}><i className="bi bi-exclamation-circle" /> {errors.password}</span>}
            </div>

            {/* Captcha */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
              <label style={{ fontSize: '13px', fontWeight: '500', color: 'var(--text-primary)' }}>验证码</label>
              <div style={{ display: 'flex', gap: '10px' }}>
                <Input
                  value={captcha} onChange={e => { setCaptcha(e.target.value); setErrors(p => ({...p, captcha: ''})); }}
                  placeholder="请输入验证码" error={!!errors.captcha} size="lg"
                  style={{ flex: 1 }} onEnter={handleLogin}
                />
                <div style={{
                  width: '110px', height: '40px', borderRadius: 'var(--radius-md)',
                  background: 'linear-gradient(135deg, var(--primary-light), var(--bg-page))',
                  border: '1px solid var(--border)', display: 'flex', alignItems: 'center',
                  justifyContent: 'center', cursor: 'pointer', flexShrink: 0,
                  fontFamily: 'var(--font-mono)', fontWeight: '800', letterSpacing: '6px',
                  fontSize: '16px', color: 'var(--primary)', userSelect: 'none',
                }}>
                  {captchaVal}
                </div>
              </div>
              {errors.captcha && <span style={{ fontSize: '12px', color: 'var(--danger)' }}><i className="bi bi-exclamation-circle" /> {errors.captcha}</span>}
            </div>

            <Button type="submit" block size="lg" loading={loading} style={{ marginTop: '4px' }}>
              {loading ? '登录中...' : '登 录'}
            </Button>
          </form>

          {/* Quick demo hint */}
          <div style={{ marginTop: '20px', padding: '12px', background: 'var(--bg-page)', borderRadius: 'var(--radius-lg)', border: '1px solid var(--border)' }}>
            <p style={{ fontSize: '12px', color: 'var(--text-secondary)', marginBottom: '4px' }}>
              <i className="bi bi-info-circle" style={{ marginRight: '5px', color: 'var(--info)' }} />
              演示账号（验证码：<code style={{ fontFamily: 'var(--font-mono)', fontWeight: '700', color: 'var(--primary)' }}>{captchaVal}</code>）
            </p>
            <div style={{ display: 'flex', gap: '8px', marginTop: '6px', flexWrap: 'wrap' }}>
              {[{ u: 'admin', p: 'admin123' }, { u: 'operator', p: 'pass1234' }].map(acc => (
                <button key={acc.u} onClick={() => { setUsername(acc.u); setPassword(acc.p); setErrors({}); }}
                  style={{ fontSize: '12px', padding: '3px 10px', borderRadius: 'var(--radius-full)', border: '1px solid var(--border)', background: 'var(--bg-surface)', cursor: 'pointer', color: 'var(--text-secondary)', fontFamily: 'var(--font-sans)' }}>
                  {acc.u}
                </button>
              ))}
            </div>
          </div>

          {/* Theme color picker */}
          <div style={{ marginTop: '24px', display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '10px' }}>
            <span style={{ fontSize: '12px', color: 'var(--text-tertiary)' }}>主题色</span>
            <div style={{ display: 'flex', gap: '6px' }}>
              {Object.entries(BaseTokens.PRESET_HUES).map(([name, hue]) => (
                <div key={name} title={name}
                  onClick={() => BaseTokens.applyHue(hue)}
                  style={{
                    width: '18px', height: '18px', borderRadius: '50%', cursor: 'pointer',
                    background: `oklch(0.55 0.18 ${hue})`,
                    border: BaseTokens.getHue() === hue ? '2px solid var(--text-primary)' : '2px solid transparent',
                    outline: BaseTokens.getHue() === hue ? '2px solid transparent' : 'none',
                    transition: 'transform 0.15s',
                    boxShadow: '0 0 0 1px rgba(0,0,0,0.1)',
                  }}
                  onMouseEnter={e => e.currentTarget.style.transform = 'scale(1.2)'}
                  onMouseLeave={e => e.currentTarget.style.transform = ''}
                />
              ))}
            </div>
          </div>
        </div>

        <p style={{ marginTop: '40px', fontSize: '12px', color: 'var(--text-tertiary)', textAlign: 'center' }}>
          © 2024 BaseProject. All rights reserved.
        </p>
      </div>

      {/* Responsive hide left panel */}
      <style>{`
        @media (max-width: 768px) {
          .login-panel-left { display: none !important; }
        }
      `}</style>
    </div>
  );
};

Object.assign(window, { Login });
