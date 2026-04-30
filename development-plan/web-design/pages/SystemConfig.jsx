// pages/SystemConfig.jsx
const SystemConfig = () => {
  const [activeTab, setActiveTab] = React.useState('basic');
  const [loading, setLoading] = React.useState(false);
  const [hue, setHue] = React.useState(BaseTokens.getHue());

  const handleSave = (section) => {
    setLoading(true);
    setTimeout(() => { setLoading(false); window.toast.success(`${section}配置保存成功`); }, 800);
  };

  const tabs = [
    { key:'basic',    label:'基本信息', icon:'bi-info-circle' },
    { key:'security', label:'安全策略', icon:'bi-shield-lock' },
    { key:'theme',    label:'界面主题', icon:'bi-palette' },
    { key:'notify',   label:'通知设置', icon:'bi-bell' },
    { key:'storage',  label:'存储配置', icon:'bi-hdd' },
  ];

  return (
    <div style={{ display:'flex', flexDirection:'column', gap:'16px', animation:'bpFadeIn 0.25s ease' }}>
      <div className="bp-page-header">
        <div>
          <div className="bp-page-title">系统配置</div>
          <div className="bp-page-subtitle">全局系统参数与偏好设置</div>
        </div>
      </div>

      <div style={{ display:'flex', gap:'16px', alignItems:'flex-start' }}>
        {/* Tab list */}
        <Card padding="8px" style={{ width:'200px', flexShrink:0 }}>
          <div style={{ display:'flex', flexDirection:'column', gap:'2px' }}>
            {tabs.map(t => (
              <ConfigTabBtn key={t.key} tab={t} active={activeTab===t.key} onClick={() => setActiveTab(t.key)} />
            ))}
          </div>
        </Card>

        {/* Content */}
        <div style={{ flex:1, minWidth:0 }}>
          {activeTab === 'basic' && (
            <Card title="基本信息" extra={<Button size="sm" onClick={() => handleSave('基本信息')} loading={loading}>保存</Button>}>
              <Form style={{ gap:'18px', maxWidth:'560px' }}>
                <FormItem label="系统名称" required><Input defaultValue="BaseProject" /></FormItem>
                <FormItem label="系统副标题"><Input defaultValue="统一基座管理平台" /></FormItem>
                <FormItem label="系统Logo" help="建议上传 64×64 PNG/SVG 图标">
                  <div style={{ display:'flex', gap:'12px', alignItems:'center' }}>
                    <div style={{ width:'56px', height:'56px', borderRadius:'var(--radius-lg)', background:'var(--primary)', display:'flex', alignItems:'center', justifyContent:'center', color:'#fff', fontWeight:'800', fontSize:'22px' }}>B</div>
                    <Button variant="secondary" size="sm" icon={<i className="bi bi-upload" />}>上传图标</Button>
                  </div>
                </FormItem>
                <FormItem label="版本号"><Input defaultValue="v1.2.0" readOnly /></FormItem>
                <FormItem label="版权信息"><Input defaultValue="© 2024 BaseProject. All rights reserved." /></FormItem>
                <FormItem label="ICP 备案"><Input placeholder="选填，如：粤ICP备XXXXXXXX号" /></FormItem>
                <FormItem label="系统公告"><Textarea defaultValue="欢迎使用 BaseProject 管理平台，如有问题请联系系统管理员。" rows={3} /></FormItem>
              </Form>
            </Card>
          )}

          {activeTab === 'security' && (
            <Card title="安全策略" extra={<Button size="sm" onClick={() => handleSave('安全策略')} loading={loading}>保存</Button>}>
              <Form style={{ gap:'20px', maxWidth:'560px' }}>
                <FormDivider label="密码策略" />
                <FormItem label="密码最小长度" help="建议不低于 8 位">
                  <Input type="number" defaultValue="8" style={{ width:'100px' }} />
                </FormItem>
                <FormItem label="密码复杂度要求">
                  <div style={{ display:'flex', flexDirection:'column', gap:'8px' }}>
                    <Checkbox defaultChecked>必须包含大写字母</Checkbox>
                    <Checkbox defaultChecked>必须包含小写字母</Checkbox>
                    <Checkbox defaultChecked>必须包含数字</Checkbox>
                    <Checkbox>必须包含特殊字符</Checkbox>
                  </div>
                </FormItem>
                <FormItem label="密码有效期" help="0 表示永不过期">
                  <div style={{ display:'flex', alignItems:'center', gap:'8px' }}>
                    <Input type="number" defaultValue="90" style={{ width:'100px' }} />
                    <span style={{ fontSize:'13px', color:'var(--text-secondary)' }}>天</span>
                  </div>
                </FormItem>
                <FormDivider label="登录策略" />
                <FormItem label="最大登录失败次数">
                  <div style={{ display:'flex', alignItems:'center', gap:'8px' }}>
                    <Input type="number" defaultValue="5" style={{ width:'100px' }} />
                    <span style={{ fontSize:'13px', color:'var(--text-secondary)' }}>次后锁定账号</span>
                  </div>
                </FormItem>
                <FormItem label="账号锁定时长">
                  <div style={{ display:'flex', alignItems:'center', gap:'8px' }}>
                    <Input type="number" defaultValue="30" style={{ width:'100px' }} />
                    <span style={{ fontSize:'13px', color:'var(--text-secondary)' }}>分钟</span>
                  </div>
                </FormItem>
                <FormItem label="Session 超时">
                  <div style={{ display:'flex', alignItems:'center', gap:'8px' }}>
                    <Input type="number" defaultValue="120" style={{ width:'100px' }} />
                    <span style={{ fontSize:'13px', color:'var(--text-secondary)' }}>分钟</span>
                  </div>
                </FormItem>
                <FormItem label="启用验证码"><RadioGroup value="true" options={[{value:'true',label:'启用'},{value:'false',label:'关闭'}]} /></FormItem>
                <FormItem label="允许多端同时登录"><RadioGroup value="false" options={[{value:'true',label:'允许'},{value:'false',label:'踢出旧登录'}]} /></FormItem>
              </Form>
            </Card>
          )}

          {activeTab === 'theme' && (
            <Card title="界面主题" extra={<Button size="sm" onClick={() => handleSave('界面主题')} loading={loading}>保存</Button>}>
              <Form style={{ gap:'20px', maxWidth:'560px' }}>
                <FormItem label="主题色" help="选择系统全局主色调">
                  <div style={{ display:'flex', flexDirection:'column', gap:'12px' }}>
                    <div style={{ display:'flex', gap:'10px', flexWrap:'wrap' }}>
                      {Object.entries(BaseTokens.PRESET_HUES).map(([name, h]) => (
                        <div key={name} onClick={() => { BaseTokens.applyHue(h); setHue(h); }}
                          style={{ display:'flex', flexDirection:'column', alignItems:'center', gap:'6px', cursor:'pointer' }}>
                          <div style={{
                            width:'36px', height:'36px', borderRadius:'50%',
                            background:`oklch(0.55 0.18 ${h})`,
                            border: hue === h ? '3px solid var(--text-primary)' : '3px solid transparent',
                            boxShadow:'0 0 0 1px rgba(0,0,0,0.1)',
                            transition:'transform 0.15s',
                          }}
                            onMouseEnter={e => e.currentTarget.style.transform='scale(1.1)'}
                            onMouseLeave={e => e.currentTarget.style.transform=''}
                          />
                          <span style={{ fontSize:'11px', color: hue===h ? 'var(--primary)' : 'var(--text-tertiary)', fontWeight: hue===h ? '600' : '400' }}>{name}</span>
                        </div>
                      ))}
                    </div>
                    <div style={{ display:'flex', alignItems:'center', gap:'10px' }}>
                      <span style={{ fontSize:'13px', color:'var(--text-secondary)' }}>自定义色相</span>
                      <input type="range" min="0" max="360" value={hue}
                        onChange={e => { const v=Number(e.target.value); setHue(v); BaseTokens.applyHue(v); }}
                        style={{ flex:1, accentColor:'var(--primary)' }} />
                      <div style={{ width:'28px', height:'28px', borderRadius:'50%', background:`oklch(0.55 0.18 ${hue})`, flexShrink:0 }} />
                    </div>
                  </div>
                </FormItem>
                <FormItem label="默认主题模式">
                  <RadioGroup value={BaseTokens.isDark() ? 'dark' : 'light'}
                    onChange={v => BaseTokens.applyDark(v==='dark')}
                    options={[{value:'light',label:'浅色模式'},{value:'dark',label:'深色模式'},{value:'auto',label:'跟随系统'}]} />
                </FormItem>
                <FormItem label="侧边栏风格">
                  <RadioGroup value="dark" options={[{value:'dark',label:'深色侧边栏'},{value:'light',label:'浅色侧边栏'}]} />
                </FormItem>
                <FormItem label="界面圆角">
                  <RadioGroup value="default" options={[{value:'none',label:'无圆角'},{value:'default',label:'默认'},{value:'large',label:'大圆角'}]} />
                </FormItem>
              </Form>
            </Card>
          )}

          {activeTab === 'notify' && (
            <Card title="通知设置" extra={<Button size="sm" onClick={() => handleSave('通知设置')} loading={loading}>保存</Button>}>
              <Form style={{ gap:'20px', maxWidth:'560px' }}>
                <FormDivider label="邮件通知" />
                <FormItem label="SMTP 服务器"><Input defaultValue="smtp.example.com" /></FormItem>
                <FormRow cols={2}>
                  <FormItem label="端口"><Input type="number" defaultValue="465" /></FormItem>
                  <FormItem label="加密方式"><Select defaultValue="SSL" options={['SSL','TLS','无']} /></FormItem>
                </FormRow>
                <FormRow cols={2}>
                  <FormItem label="发件账号"><Input defaultValue="noreply@example.com" /></FormItem>
                  <FormItem label="授权码"><Input type="password" placeholder="邮箱授权码" /></FormItem>
                </FormRow>
                <FormItem label="发件人名称"><Input defaultValue="BaseProject 系统" /></FormItem>
                <div style={{ display:'flex', gap:'8px' }}>
                  <Button variant="secondary" size="sm" icon={<i className="bi bi-send" />} onClick={() => window.toast.success('测试邮件已发送')}>发送测试邮件</Button>
                </div>
                <FormDivider label="通知事件" />
                {['用户注册通知','密码重置通知','异常登录告警','系统错误告警','任务完成通知'].map(item => (
                  <div key={item} style={{ display:'flex', justifyContent:'space-between', alignItems:'center', padding:'8px 0', borderBottom:'1px solid var(--border)' }}>
                    <span style={{ fontSize:'13px', color:'var(--text-primary)' }}>{item}</span>
                    <ToggleSwitch defaultOn={item.includes('异常') || item.includes('错误')} />
                  </div>
                ))}
              </Form>
            </Card>
          )}

          {activeTab === 'storage' && (
            <Card title="存储配置" extra={<Button size="sm" onClick={() => handleSave('存储配置')} loading={loading}>保存</Button>}>
              <Form style={{ gap:'18px', maxWidth:'560px' }}>
                <FormItem label="存储类型">
                  <RadioGroup value="local" options={[{value:'local',label:'本地存储'},{value:'oss',label:'阿里云 OSS'},{value:'cos',label:'腾讯云 COS'},{value:'minio',label:'MinIO'}]} />
                </FormItem>
                <FormDivider label="本地存储" />
                <FormItem label="上传路径"><Input defaultValue="/var/uploads/" style={{ fontFamily:'var(--font-mono)' }} /></FormItem>
                <FormItem label="访问域名"><Input defaultValue="https://cdn.example.com" /></FormItem>
                <FormRow cols={2}>
                  <FormItem label="单文件大小限制" help="单位：MB">
                    <Input type="number" defaultValue="10" />
                  </FormItem>
                  <FormItem label="允许的文件类型">
                    <Input defaultValue=".jpg,.png,.pdf,.xlsx" style={{ fontFamily:'var(--font-mono)', fontSize:'12px' }} />
                  </FormItem>
                </FormRow>
                <FormDivider label="存储用量" />
                {[
                  { label:'已用空间', value:68, used:'6.8 GB', total:'10 GB' },
                  { label:'文件总数', value:45, used:'4,521', total:'10,000' },
                ].map(s => (
                  <div key={s.label} style={{ display:'flex', flexDirection:'column', gap:'6px' }}>
                    <div style={{ display:'flex', justifyContent:'space-between', fontSize:'13px' }}>
                      <span style={{ color:'var(--text-secondary)' }}>{s.label}</span>
                      <span style={{ color:'var(--text-primary)', fontWeight:'500' }}>{s.used} / {s.total}</span>
                    </div>
                    <div style={{ height:'6px', background:'var(--bg-page)', borderRadius:'var(--radius-full)', overflow:'hidden' }}>
                      <div style={{ height:'100%', width:`${s.value}%`, background: s.value > 80 ? 'var(--danger)' : s.value > 60 ? 'var(--warning)' : 'var(--primary)', borderRadius:'var(--radius-full)' }} />
                    </div>
                  </div>
                ))}
              </Form>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
};

const ConfigTabBtn = ({ tab, active, onClick }) => {
  const [hov, setHov] = React.useState(false);
  return (
    <div onClick={onClick}
      onMouseEnter={() => setHov(true)} onMouseLeave={() => setHov(false)}
      style={{
        display:'flex', alignItems:'center', gap:'8px', padding:'9px 12px',
        borderRadius:'var(--radius-md)', cursor:'pointer', transition:'all 0.15s',
        background: active ? 'var(--primary-light)' : hov ? 'var(--bg-page)' : 'transparent',
        color: active ? 'var(--primary)' : 'var(--text-secondary)',
        fontWeight: active ? '600' : '400', fontSize:'13px',
      }}>
      <i className={`bi ${tab.icon}`} style={{ fontSize:'14px' }} />
      {tab.label}
    </div>
  );
};

const ToggleSwitch = ({ defaultOn = false }) => {
  const [on, setOn] = React.useState(defaultOn);
  return (
    <div onClick={() => setOn(p=>!p)}
      style={{ width:'36px', height:'20px', background: on ? 'var(--primary)' : 'var(--border-strong)', borderRadius:'10px', position:'relative', cursor:'pointer', transition:'background 0.2s', flexShrink:0 }}>
      <div style={{ position:'absolute', top:'2px', left: on ? '18px' : '2px', width:'16px', height:'16px', background:'#fff', borderRadius:'50%', transition:'left 0.2s', boxShadow:'0 1px 3px rgba(0,0,0,0.2)' }} />
    </div>
  );
};

Object.assign(window, { SystemConfig });
