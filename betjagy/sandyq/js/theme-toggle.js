// Theme Toggle JavaScript
// Переключение светлой/темной темы

class ThemeToggle {
  constructor() {
    this.theme = localStorage.getItem('theme') || 'light';
    this.init();
  }

  init() {
    this.setTheme(this.theme);
    this.createToggleButton();
    this.bindOffcanvasEvents();
    this.bindEvents();
  }

  setTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
    this.theme = theme;
    localStorage.setItem('theme', theme);
    
    // Update logo based on theme
    this.updateLogo(theme);
  }

  updateLogo(theme) {
    const logoElements = document.querySelectorAll('.logo');
    logoElements.forEach(logo => {
      if (theme === 'dark') {
        logo.src = logo.src.replace('logo_red.svg', 'logo_white.svg');
      } else {
        logo.src = logo.src.replace('logo_white.svg', 'logo_red.svg');
      }
    });
    
    // Update all theme toggle buttons
    const buttons = document.querySelectorAll('.theme-toggle');
    buttons.forEach(button => {
      const icon = button.querySelector('.theme-icon');
      if (icon) {
        icon.className = 'theme-icon bi ' + (theme === 'dark' ? 'bi-sun' : 'bi-moon');
      }
    });
  }

  createToggleButton() {
    const navbar = document.querySelector('.navbar .container');
    if (!navbar) return;

    // Check if theme toggle already exists
    const existingToggle = navbar.querySelector('.theme-toggle');
    if (existingToggle) return;

    const toggleButton = document.createElement('button');
    toggleButton.className = 'theme-toggle';
    toggleButton.innerHTML = `
      <span class="theme-icon bi ${this.theme === 'dark' ? 'bi-sun' : 'bi-moon'}"></span>
    `;
    
    // Add to navbar
    const nav = navbar.querySelector('.navbar-nav');
    if (nav) {
      nav.appendChild(toggleButton);
    } else {
      navbar.appendChild(toggleButton);
    }
  }

  toggleTheme() {
    const newTheme = this.theme === 'light' ? 'dark' : 'light';
    this.setTheme(newTheme);
  }

  bindOffcanvasEvents() {
    // Bind events for all theme toggle buttons
    const allToggleButtons = document.querySelectorAll('.theme-toggle');
    allToggleButtons.forEach(button => {
      button.addEventListener('click', () => this.toggleTheme());
    });
  }

  bindEvents() {
    // Listen for system theme changes
    if (window.matchMedia) {
      const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
      mediaQuery.addEventListener('change', (e) => {
        if (!localStorage.getItem('theme')) {
          this.setTheme(e.matches ? 'dark' : 'light');
        }
      });
    }
  }
}

// Initialize theme toggle when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  new ThemeToggle();
});

// Export for module usage
if (typeof module !== 'undefined' && module.exports) {
  module.exports = ThemeToggle;
}
