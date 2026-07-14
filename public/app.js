(() => {
  const slugify = (value) => value.toLowerCase()
    .normalize("NFD").replace(/[\u0300-\u036f]/g, "")
    .replace(/&/g, "and").replace(/[^a-z0-9]+/g, "-").replace(/(^-|-$)/g, "");

  const pricePattern = /(?:S\s*)?\$\s*\d+(?:\.\d+)?(?:\s*\/\s*L\s*\$\s*\d+(?:\.\d+)?)?|\b\d+\.\d{2}\b/gi;
  const introPattern = /^(choice of side|comes with|add topping)/i;

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
        current = { name, price: prices.join(" / ").replace(/\s+/g, " "), description: "" };
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
            <span>${item.name}</span>
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
  toggle?.addEventListener("click", () => {
    const open = toggle.getAttribute("aria-expanded") === "true";
    toggle.setAttribute("aria-expanded", String(!open));
    nav.classList.toggle("open", !open);
  });
  nav?.querySelectorAll("a").forEach((link) => link.addEventListener("click", () => {
    nav.classList.remove("open");
    toggle?.setAttribute("aria-expanded", "false");
  }));

  document.querySelector("#year").textContent = new Date().getFullYear();
  renderMenu();
})();
