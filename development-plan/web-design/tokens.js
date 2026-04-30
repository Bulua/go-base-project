// tokens.js — BaseProject theming engine
(function () {
  var PRESET_HUES = {
    '蓝色': 250,
    '天蓝': 205,
    '绿色': 145,
    '紫色': 285,
    '玫红': 340,
    '橙色': 52,
  };

  function applyHue(hue) {
    document.documentElement.style.setProperty('--hue', hue);
    try { localStorage.setItem('bp-hue', hue); } catch(e){}
  }

  function applyDark(isDark) {
    document.documentElement.setAttribute('data-theme', isDark ? 'dark' : 'light');
    try { localStorage.setItem('bp-dark', isDark ? '1' : '0'); } catch(e){}
  }

  // Restore persisted settings
  var hue = 250, dark = false;
  try {
    var sh = localStorage.getItem('bp-hue');
    if (sh) hue = parseFloat(sh);
    dark = localStorage.getItem('bp-dark') === '1';
  } catch(e){}
  applyHue(hue);
  applyDark(dark);

  window.BaseTokens = {
    PRESET_HUES: PRESET_HUES,
    applyHue: applyHue,
    applyDark: applyDark,
    getHue: function() {
      return parseFloat(document.documentElement.style.getPropertyValue('--hue') || '250');
    },
    isDark: function() {
      return document.documentElement.getAttribute('data-theme') === 'dark';
    },
  };
})();
