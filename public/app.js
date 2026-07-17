(() => {
  const THEME_KEY = "lotus-theme";

  function getStoredTheme() {
    try {
      const value = localStorage.getItem(THEME_KEY);
      return value === "light" || value === "dark" ? value : null;
    } catch (_) {
      return null;
    }
  }

  function getSystemTheme() {
    return window.matchMedia?.("(prefers-color-scheme: dark)").matches ? "dark" : "light";
  }

  const themeToggle = document.querySelector("#theme-toggle");
  const themeLabel = themeToggle?.querySelector(".theme-toggle-label");
  const themeColor = document.querySelector('meta[name="theme-color"]');

  function applyTheme(theme) {
    const isDark = theme === "dark";
    document.documentElement.dataset.theme = theme;
    document.documentElement.style.colorScheme = theme;
    themeToggle?.setAttribute("aria-pressed", String(isDark));
    themeToggle?.setAttribute("aria-label", `Switch to ${isDark ? "light" : "dark"} mode`);
    if (themeLabel) themeLabel.textContent = isDark ? "Light" : "Dark";
    const isMenuPage = document.documentElement.dataset.page === "menu";
    const color = isMenuPage
      ? (isDark ? "#0c2a1f" : "#173c2d")
      : (isDark ? "#0a1711" : "#f5f0e5");
    themeColor?.setAttribute("content", color);
  }

  applyTheme(getStoredTheme() || document.documentElement.dataset.theme || getSystemTheme());

  themeToggle?.addEventListener("click", () => {
    const nextTheme = document.documentElement.dataset.theme === "dark" ? "light" : "dark";
    applyTheme(nextTheme);
    try {
      localStorage.setItem(THEME_KEY, nextTheme);
    } catch (_) {}
  });

  const systemTheme = window.matchMedia?.("(prefers-color-scheme: dark)");
  systemTheme?.addEventListener?.("change", (event) => {
    if (!getStoredTheme()) applyTheme(event.matches ? "dark" : "light");
  });

  function enhanceMenu() {
    const groups = document.querySelector("#menu-groups");
    const jump = document.querySelector("#menu-jump");
    if (!groups || !jump) return;

    const links = [...jump.querySelectorAll("a")];
    const openHash = () => {
      const target = location.hash && document.querySelector(location.hash);
      if (target?.matches?.("details.menu-group")) {
        groups.querySelectorAll("details.menu-group").forEach((group) => { group.open = group === target; });
        window.setTimeout(() => target.scrollIntoView({ block: "start" }), 0);
      }
    };
    links.forEach((link) => link.addEventListener("click", (event) => {
      event.preventDefault();
      const target = document.querySelector(link.hash);
      if (!target) return;
      history.pushState(null, "", link.hash);
      groups.querySelectorAll("details.menu-group").forEach((group) => { group.open = group === target; });
      target.scrollIntoView({ behavior: "smooth", block: "start" });
    }));
    window.addEventListener("hashchange", openHash);
    openHash();

    const observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (!entry.isIntersecting) return;
        links.forEach((link) => link.classList.toggle("active", link.hash === `#${entry.target.id}`));
      });
    }, { rootMargin: "-25% 0px -68% 0px" });
    groups.querySelectorAll(".menu-group").forEach((group) => observer.observe(group));
  }

  const toggle = document.querySelector(".nav-toggle");
  const nav = document.querySelector("#primary-nav");
  const navToggleLabel = toggle?.querySelector(".sr-only");
  toggle?.addEventListener("click", () => {
    const open = toggle.getAttribute("aria-expanded") === "true";
    toggle.setAttribute("aria-expanded", String(!open));
    nav.classList.toggle("open", !open);
    if (navToggleLabel) navToggleLabel.textContent = open ? "Open navigation" : "Close navigation";
  });
  nav?.querySelectorAll("a").forEach((link) => link.addEventListener("click", () => {
    nav.classList.remove("open");
    toggle?.setAttribute("aria-expanded", "false");
    if (navToggleLabel) navToggleLabel.textContent = "Open navigation";
  }));

  const year = document.querySelector("#year");
  if (year) year.textContent = new Date().getFullYear();

  const contactForm = document.querySelector("#contact-form");
  const contactStatus = document.querySelector("#contact-status");
  contactForm?.addEventListener("submit", (event) => {
    event.preventDefault();
    if (contactStatus) contactStatus.textContent = "Thanks — this form is ready to connect to the Lotus inbox.";
    contactForm.reset();
  });

  enhanceMenu();
})();
