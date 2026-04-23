const usernameEl = document.getElementById("username");
const passwordEl = document.getElementById("password");
const postContentEl = document.getElementById("postContent");

const registerBtn = document.getElementById("registerBtn");
const loginBtn = document.getElementById("loginBtn");
const createPostBtn = document.getElementById("createPostBtn");
const loadPostsBtn = document.getElementById("loadPostsBtn");

const authStatusEl = document.getElementById("authStatus");
const postStatusEl = document.getElementById("postStatus");
const postsOutputEl = document.getElementById("postsOutput");

const TOKEN_KEY = "posts_service_token";
let jwtToken = localStorage.getItem(TOKEN_KEY) || "";

function baseUrl() {
  if (window.location.protocol === "file:") {
    return "http://localhost:8081";
  }
  return window.location.origin;
}

function setStatus(el, type, text) {
  el.style.display = "block";
  el.className = `status ${type}`;
  el.textContent = text;
}

function setToken(token) {
  jwtToken = token || "";
  if (jwtToken) {
    localStorage.setItem(TOKEN_KEY, jwtToken);
  } else {
    localStorage.removeItem(TOKEN_KEY);
  }
}

async function request(path, options = {}) {
  const url = `${baseUrl()}${path}`;
  const headers = options.headers || {};
  const merged = {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...headers
    }
  };

  const res = await fetch(url, merged);
  const text = await res.text();
  return { ok: res.ok, status: res.status, text };
}

function requireCredentials() {
  const username = usernameEl.value.trim();
  const password = passwordEl.value.trim();
  if (!username || !password) {
    setStatus(authStatusEl, "err", "Введите username и password.");
    return null;
  }
  return { username, password };
}

registerBtn.addEventListener("click", async () => {
  const creds = requireCredentials();
  if (!creds) return;

  try {
    const res = await request("/register", {
      method: "POST",
      body: JSON.stringify(creds)
    });

    if (res.ok) {
      setStatus(authStatusEl, "ok", `Регистрация успешна: ${res.text}`);
    } else {
      setStatus(authStatusEl, "err", `Ошибка регистрации (${res.status}): ${res.text}`);
    }
  } catch (error) {
    setStatus(authStatusEl, "err", `Ошибка сети: ${error.message}`);
  }
});

loginBtn.addEventListener("click", async () => {
  const creds = requireCredentials();
  if (!creds) return;

  try {
    const res = await request("/login", {
      method: "POST",
      body: JSON.stringify(creds)
    });

    if (res.ok) {
      setToken(res.text.trim());
      setStatus(authStatusEl, "ok", "Логин успешен.");
    } else {
      setStatus(authStatusEl, "err", `Ошибка логина (${res.status}): ${res.text}`);
    }
  } catch (error) {
    setStatus(authStatusEl, "err", `Ошибка сети: ${error.message}`);
  }
});

createPostBtn.addEventListener("click", async () => {
  const content = postContentEl.value.trim();
  if (!jwtToken) {
    setStatus(postStatusEl, "err", "Сначала выполните login, чтобы получить JWT.");
    return;
  }
  if (!content) {
    setStatus(postStatusEl, "err", "Введите текст поста.");
    return;
  }

  try {
    const res = await request("/posts", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${jwtToken}`
      },
      body: JSON.stringify({ content })
    });

    if (res.ok) {
      setStatus(postStatusEl, "ok", `Пост создан: ${res.text}`);
      postContentEl.value = "";
    } else {
      setStatus(postStatusEl, "err", `Ошибка создания поста (${res.status}): ${res.text}`);
    }
  } catch (error) {
    setStatus(postStatusEl, "err", `Ошибка сети: ${error.message}`);
  }
});

loadPostsBtn.addEventListener("click", async () => {
  if (!jwtToken) {
    postsOutputEl.textContent = "Нет токена. Сначала выполните login.";
    return;
  }

  try {
    const res = await request("/posts", {
      method: "GET",
      headers: {
        Authorization: `Bearer ${jwtToken}`
      }
    });

    if (res.ok) {
      postsOutputEl.textContent = res.text || "(Постов пока нет)";
    } else {
      postsOutputEl.textContent = `Ошибка (${res.status}): ${res.text}`;
    }
  } catch (error) {
    postsOutputEl.textContent = `Ошибка сети: ${error.message}`;
  }
});
