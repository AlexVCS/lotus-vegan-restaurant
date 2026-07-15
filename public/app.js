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
    themeColor?.setAttribute("content", isDark ? "#0a1711" : "#f5f0e5");
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

  const slugify = (value) => value.toLowerCase()
    .normalize("NFD").replace(/[\u0300-\u036f]/g, "")
    .replace(/&/g, "and").replace(/[^a-z0-9]+/g, "-").replace(/(^-|-$)/g, "");

  const pricePattern = /(?:S\s*)?\$\s*\d+(?:\.\d+)?(?:\s*\/\s*L\s*\$\s*\d+(?:\.\d+)?)?|\b\d+\.\d{2}\b/gi;
  const introPattern = /^(choice of side|comes with|add topping)/i;
  const dietaryLabels = {
    GF: "Gluten-free",
    SF: "Soy-free",
  };
  const dietaryTagPattern = /\((GF|SF)\)/gi;
  let dietaryTooltipIndex = 0;

  function extractDietaryTags(name) {
    const dietaryTags = [];
    const cleanName = name.replace(dietaryTagPattern, (_, tag) => {
      const normalizedTag = tag.toUpperCase();
      if (!dietaryTags.includes(normalizedTag)) dietaryTags.push(normalizedTag);
      return " ";
    }).replace(/\s+/g, " ").trim();

    return { name: cleanName, dietaryTags };
  }

  function renderDietaryTags(tags = []) {
    if (!tags.length) return "";
    return `
      <span class="menu-item-dietary">
        ${tags.map((tag) => {
          const label = dietaryLabels[tag];
          const tooltipId = `dietary-tooltip-${++dietaryTooltipIndex}`;
          return `<span class="dietary-badge dietary-badge-tooltip dietary-badge-${tag.toLowerCase()}" tabindex="0" aria-label="${tag}" aria-describedby="${tooltipId}"><span aria-hidden="true">${tag}</span><span class="dietary-tooltip" id="${tooltipId}" role="tooltip">${label}</span></span>`;
        }).join("")}
      </span>`;
  }

  function parseItems(lines) {
    const items = [];
    let current = null;

    lines.forEach((raw) => {
      const line = raw.replace(/\s+/g, " ").trim();
      if (!line) return;
      const looksLikeItem = pricePattern.test(line) && !introPattern.test(line);
      pricePattern.lastIndex = 0;

      if (looksLikeItem) {
        if (current) items.push(current);
        const prices = line.match(pricePattern) || [];
        let name = line.replace(pricePattern, "").replace(/\s*[-/]\s*$/, "").trim();
        name = name.replace(/\s+-\s*$/, "").replace(/^[-–]\s*/, "");
        const dietary = extractDietaryTags(name);
        current = {
          name: dietary.name,
          dietaryTags: dietary.dietaryTags,
          price: prices.join(" / ").replace(/\s+/g, " "),
          description: "",
        };
      } else if (current) {
        current.description += (current.description ? " " : "") + line;
      } else {
        items.push({ name: "Good to know", price: "", description: line, note: true });
      }
    });

    if (current) items.push(current);
    return items;
  }

  function renderMenu() {
    const groups = document.querySelector("#menu-groups");
    const jump = document.querySelector("#menu-jump");
    const data = window.LOTUS_MENU || [];
    if (!data.length) {
      groups.innerHTML = '<p class="menu-empty">The menu could not be loaded. Please call us at <a href="tel:+17275491688">727-549-1688</a>.</p>';
      return;
    }

    groups.innerHTML = "";
    dietaryTooltipIndex = 0;
    data.forEach((category, index) => {
      const id = `menu-${slugify(category.title)}`;
      const shortTitle = category.title
        .replace("Beverages: ", "")
        .replace("Fresh Juices / Fruit Teas", "Juices & teas")
        .replace("Additional Beverages", "More drinks");

      const link = document.createElement("a");
      link.href = `#${id}`;
      link.textContent = shortTitle;
      jump.append(link);

      const details = document.createElement("details");
      details.className = "menu-group";
      details.id = id;
      details.open = index === 0;
      details.innerHTML = `
        <summary>
          <span class="menu-group-index">${String(index + 1).padStart(2, "0")}</span>
          <h3>${category.title}</h3>
        </summary>
        <div class="menu-items"></div>`;

      const list = details.querySelector(".menu-items");
      parseItems(category.lines).forEach((item) => {
        const article = document.createElement("article");
        article.className = `menu-item${item.note ? " menu-item-note" : ""}`;
        article.innerHTML = `
          <h4 class="menu-item-name">
            <span class="menu-item-heading">
              <span class="menu-item-title">${item.name}</span>
              ${renderDietaryTags(item.dietaryTags)}
            </span>
            ${item.price ? '<i class="dots" aria-hidden="true"></i>' : ""}
            ${item.price ? `<span class="menu-price">${item.price}</span>` : ""}
          </h4>
          ${item.description ? `<p>${item.description}</p>` : ""}`;
        list.append(article);
      });
      groups.append(details);
    });

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

  document.querySelector("#year").textContent = new Date().getFullYear();
  renderMenu();
})();
